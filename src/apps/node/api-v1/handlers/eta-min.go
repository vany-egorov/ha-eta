package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	v1 "github.com/vany-egorov/ha-eta/apps/node/api-v1"
	apiErrors "github.com/vany-egorov/ha-eta/apps/node/api-v1/errors"
	bufPool "github.com/vany-egorov/ha-eta/lib/buf-pool"
	cache "github.com/vany-egorov/ha-eta/lib/cache"
	geoEngine "github.com/vany-egorov/ha-eta/lib/geo-engine"
	"github.com/vany-egorov/ha-eta/lib/log"
	"github.com/vany-egorov/ha-eta/models"
)

const (
	timeout = 1000*time.Millisecond - 30*time.Millisecond // TODO: config
)

type geoDelegate struct {
	prefix string
	fnLog  func(log.Level, string)
}

/* event-delegate */
func (it *geoDelegate) LogReq(_p string, req *http.Request, r io.Reader) {
	sz := uint64(req.ContentLength)

	if s, ok := r.(fmt.Stringer); ok && sz != 0 {
		buf := bufPool.NewBuf()
		defer buf.Release()

		fmt.Fprintf(buf, "%s (:geo) [>] (sz: %s :: %s)", it.prefix, humanize.Bytes(sz), s.String())
		it.fnLog(log.Debug, buf.String())
	}
}

/* event-delegate */
func (it *geoDelegate) LogReqRes(_p string, emit func(*bytes.Buffer)) {
	buf := bufPool.NewBuf()
	defer buf.Release()

	fmt.Fprintf(buf, "%s (:geo) [<>] ", it.prefix)
	emit(&buf.Buffer)

	it.fnLog(log.Info, buf.String())
}

/* event-delegate */
func (it *geoDelegate) LogResp(_p string, resp *http.Response, w io.Writer) {
	sz := uint64(resp.ContentLength)

	if s, ok := w.(fmt.Stringer); ok && sz != 0 {
		buf := bufPool.NewBuf()
		defer buf.Release()

		fmt.Fprintf(buf, "%s (:geo) [<] (:sz %s :: %s)", it.prefix, humanize.Bytes(sz), s.String())
		it.fnLog(log.Debug, buf.String())
	}
}

func etaMin(ctx context.Context, sctx ETAMinCtx, prefix string, point models.Point) (models.ETA, error) {
	var (
		points = models.Points{}
		etas   = models.ETAs{}

		pointsHits = models.Points{}
		etasHits   = models.ETAs{}

		pointsMiss = models.Points{}
		etasMiss   = models.ETAs{}

		geo       = sctx.GeoEngine()
		cch       = sctx.Cache()
		carsLimit = geo.CarsLimit()
	)

	sctx.FnLog()(log.Info, fmt.Sprintf("%s [<] (:lat %.7f :lng %.7f)",
		prefix, point.Lat, point.Lng))

	events := geoDelegate{prefix, sctx.FnLog()}

	if ok := cch.GetPoints(point, carsLimit, &points); !ok {
		if err := geo.DoCars(ctx, point.Lat, point.Lng, carsLimit, &points, &events); err != nil {
			return 0, errors.Wrap(apiErrors.ETAMinGeoEngineCars, err.Error())
		}

		cch.SetPoints(point, carsLimit, points)
	}

	if cch.GetETAs(point, points, &pointsHits, &pointsMiss, &etasHits); len(pointsMiss) != 0 {
		if err := geo.DoPredict(ctx, point.Lat, point.Lng, pointsMiss, &etasMiss, &events); err != nil {
			return 0, errors.Wrap(apiErrors.ETAMinGeoEnginePredict, err.Error())
		}

		etasMiss.Merge(etasHits)
		pointsMiss.Merge(pointsHits)

		etas = etasMiss
		points = pointsMiss

		cch.SetETAs(point, points, etas)
	} else {
		etas = etasHits
	}

	if len(etas) == 0 {
		return 0, apiErrors.ETAMinNoETAsFound
	}

	return etas.Min(), nil
}

type ETAMinCtx interface {
	GeoEngine() geoEngine.Engine
	Cache() cache.Cache
	FnLog() func(log.Level, string)
}

func ETAMin(c *gin.Context) {
	req := v1.ReqETAMin{}

	if err := c.Bind(&req); err != nil {
		v1.Send(c.Writer, apiErrors.ETAMinReqParse)
		return
	}

	if err := req.Validate(); err != nil {
		v1.Send(c.Writer, err)
		return
	}

	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	sctx := c.MustGet("service-ctx").(ETAMinCtx)
	prefix := c.MustGet("prefix").(string)

	minETA, err := etaMin(ctx, sctx, prefix, req.Point)
	if err != nil {
		v1.Send(c.Writer, err)
		return
	}

	c.Writer.Write(minETA.Bytes())
}
