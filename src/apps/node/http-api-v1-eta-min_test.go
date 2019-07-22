package node

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vany-egorov/ha-eta/lib/log"
	"github.com/vany-egorov/ha-eta/models"
)

func fnRandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func fnRandomLat() float64 {
	return fnRandomFloat64(-90, 90)
}

func fnRandomLng() float64 {
	return fnRandomFloat64(-180, 180)
}

func fnRandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

type geoEngineMock struct {
	minETA models.ETA

	carsLimit uint64
	carsLat   float64
	carsLng   float64

	predictLat float64
	predictLng float64

	points models.Points
	etas   models.ETAs
}

func (it *geoEngineMock) CarsLimit() uint64 { return it.carsLimit }
func (it *geoEngineMock) DoCars(ctx context.Context, lat, lng float64, limit uint64, any, events interface{}) error {
	it.carsLat = lat
	it.carsLng = lng

	if lat < -90 || lat > 90 {
		return fmt.Errorf("lat is out of range: %f", lat)
	}

	if lng < -180 || lng > 180 {
		return fmt.Errorf("lng is out of range: %f", lng)
	}

	out := any.(*models.Points)

	for i := 0; uint64(i) < limit; i++ {
		*out = append(*out, models.Point{fnRandomLat(), fnRandomLng()})
	}

	it.points = *out

	return nil
}
func (it *geoEngineMock) DoPredict(ctx context.Context, lat, lng float64, anySrc, anyDst, events interface{}) error {
	it.predictLat = lat
	it.predictLng = lng

	in := anySrc.(models.Points)
	out := anyDst.(*models.ETAs)

	for range in {
		*out = append(*out, models.ETA(fnRandomInt(1000, 10000)))
	}
	if len(in) > 0 {
		(*out)[len(in)-1] = it.minETA
	}

	it.etas = *out

	return nil
}

type geoEngineMockFail struct{}

func (it *geoEngineMockFail) CarsLimit() uint64 { return 0 }
func (it *geoEngineMockFail) DoCars(ctx context.Context, lat, lng float64, limit uint64, any, events interface{}) error {
	return fmt.Errorf("geo-engine do-cars error")
}
func (it *geoEngineMockFail) DoPredict(ctx context.Context, lat, lng float64, anySrc, anyDst, events interface{}) error {
	return fmt.Errorf("geo-predict error")
}

type cacheMock struct {
	points *models.Points

	etas *models.ETAs

	hitPoints bool
	hitETAs   bool
}

func (it *cacheMock) GetPoints(point models.Point, limit uint64, points *models.Points) bool {
	if it.points == nil {
		it.hitPoints = false
		return false
	}
	it.hitPoints = true
	points = it.points
	return true
}
func (it *cacheMock) SetPoints(point models.Point, limit uint64, points models.Points) {
	it.points = &points
}

func (it *cacheMock) GetETAs(point models.Point, all models.Points, hits, miss *models.Points, etas *models.ETAs) {
	if it.etas == nil {
		it.hitETAs = false
		*miss = all
		return
	}

	it.hitETAs = true
	*hits = all
	*etas = *it.etas
}
func (it *cacheMock) SetETAs(point models.Point, points models.Points, etas models.ETAs) {
	it.etas = &etas
}

func (it *cacheMock) Flush() {}

const (
	path          = "/api/v1/eta/min"
	method string = http.MethodGet
)

var _ = Describe("GET ETA min", func() {
	var (
		app      *App
		router   *gin.Engine
		recorder *httptest.ResponseRecorder

		urlRaw string

		req  *http.Request
		eReq error

		cchMock *cacheMock
	)

	var ( // factorize
		lat = float64(0)
		lng = float64(0)

		minETA    = models.ETA(fnRandomInt(1, 20))
		carsLimit = uint64(fnRandomInt(1, 100))
	)

	BeforeEach(func() {
		lat = fnRandomLat()
		lng = fnRandomLng()

		app = &App{}

		app.ctx.setFnLog(log.LogStd)

		router = app.NewRouter()
		recorder = httptest.NewRecorder()

		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
		q.Set("lng", strconv.FormatFloat(lng, 'f', -1, 64))

		urlRaw = fmt.Sprintf("%s?%s", path, q.Encode())
	})

	Context("OK", func() {
		var (
			geoMock *geoEngineMock
		)

		BeforeEach(func() {
			geoMock = &geoEngineMock{
				minETA:    minETA,
				carsLimit: carsLimit,
			}
			app.ctx.setGeoEngine(geoMock)

			cchMock = &cacheMock{}
			app.ctx.setCache(cchMock)

			req, eReq = http.NewRequest(method, urlRaw, nil)
			recorder.Body.Reset()
			router.ServeHTTP(recorder, req)
		})

		It("should execute OK", func() {
			Expect(eReq).ShouldNot(HaveOccurred())
			Expect(recorder.Code).Should(BeEquivalentTo(http.StatusOK))
		})

		It("url values should be parsed and passed well", func() {
			Expect(geoMock.carsLat).Should(BeEquivalentTo(lat))
			Expect(geoMock.carsLng).Should(BeEquivalentTo(lng))

			Expect(geoMock.predictLat).Should(BeEquivalentTo(lat))
			Expect(geoMock.predictLng).Should(BeEquivalentTo(lng))
		})

		It("points and etas should be cached", func() {
			Expect(cchMock.etas).Should(Not(BeNil()))
			Expect(cchMock.points).Should(Not(BeNil()))

			Expect(geoMock.etas).Should(BeEquivalentTo(*cchMock.etas))
			Expect(geoMock.points).Should(BeEquivalentTo(*cchMock.points))
		})

		It("for first time cache should be MISS", func() {
			Expect(cchMock.hitETAs).Should(BeFalse())
			Expect(cchMock.hitPoints).Should(BeFalse())
		})

		Describe("output value", func() {
			var (
				minETAActual = models.ETA(0)
				eParse       error
			)

			BeforeEach(func() {
				v, e := strconv.ParseUint(recorder.Body.String(), 0, 64)
				minETAActual = models.ETA(v)
				eParse = e
			})

			It("min-eta value should be parsed with no errors", func() {
				Expect(eParse).Should(Not(HaveOccurred()))
			})

			It("should return valid min-ETA value", func() {
				Expect(minETA).Should(BeEquivalentTo(minETAActual))
			})
		})

		Describe("one more time visit", func() {
			var (
				minETAActual = models.ETA(0)
				eParse       error
			)

			BeforeEach(func() {
				req, eReq = http.NewRequest(method, urlRaw, nil)
				recorder.Body.Reset()
				router.ServeHTTP(recorder, req)

				v, e := strconv.ParseUint(recorder.Body.String(), 0, 64)
				minETAActual = models.ETA(v)
				eParse = e
			})

			It("should execute OK", func() {
				Expect(eReq).ShouldNot(HaveOccurred())
				Expect(recorder.Code).Should(BeEquivalentTo(http.StatusOK))
			})

			It("for second time cache should be HIT", func() {
				Expect(cchMock.hitETAs).Should(BeTrue())
				Expect(cchMock.hitPoints).Should(BeTrue())
			})

			It("min-eta value should be parsed with no errors again", func() {
				Expect(eParse).Should(Not(HaveOccurred()))
			})

			It("should return valid min-ETA value again", func() {
				Expect(minETA).Should(BeEquivalentTo(minETAActual))
			})
		})
	})

	Context("FAIL", func() {
		var (
			geoMock *geoEngineMockFail
		)

		BeforeEach(func() {
			geoMock = &geoEngineMockFail{}
			app.ctx.setGeoEngine(geoMock)

			cchMock = &cacheMock{}
			app.ctx.setCache(cchMock)

			req, eReq = http.NewRequest(method, urlRaw, nil)
			recorder.Body.Reset()
			router.ServeHTTP(recorder, req)
		})

		It("should execute with no errors", func() {
			Expect(eReq).ShouldNot(HaveOccurred())
		})
		It("response should contain some data with error description", func() {
			Expect(recorder.Body.Len()).ShouldNot(BeZero())
		})
		It("content type should be JSON", func() {
			Expect(recorder.Header().Get("Content-Type")).Should(BeEquivalentTo("application/json"))
		})
		It("geo-engine fail should handle well", func() {
			Expect(recorder.Code).Should(BeEquivalentTo(http.StatusBadGateway))
		})
	})
})
