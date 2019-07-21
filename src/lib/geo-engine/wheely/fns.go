package wheely

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"

	bufPool "github.com/vany-egorov/ha-eta/lib/buf-pool"
	randStr "github.com/vany-egorov/ha-eta/lib/rand-str"
)

func copyURL(dst *url.URL, src *url.URL) {
	dst.Scheme = src.Scheme
	dst.Host = src.Host
	dst.Path = src.Path
}

func patchURLCars(dst *url.URL, src *url.URL) {
	copyURL(dst, src)
	dst.Path = filepath.Join(dst.Path, pathCars)
}

func patchURLPredict(dst *url.URL, src *url.URL) {
	copyURL(dst, src)
	dst.Path = filepath.Join(dst.Path, pathPredict)
}

type logReqer interface {
	LogReq(string, *http.Request, io.Reader)
}

func tryLogReqFn(any interface{}) func(string, *http.Request, io.Reader) {
	if it, ok := any.(logReqer); ok {
		return it.LogReq
	}

	return nil
}

type logReqReser interface {
	LogReqRes(string, func(*bytes.Buffer))
}

func tryLogReqResFn(any interface{}) func(string, func(*bytes.Buffer)) {
	if it, ok := any.(logReqReser); ok {
		return it.LogReqRes
	}

	return nil
}

type logResper interface {
	LogResp(string, *http.Response, io.Writer)
}

func tryLogRespFn(any interface{}) func(string, *http.Response, io.Writer) {
	if it, ok := any.(logResper); ok {
		return it.LogResp
	}

	return nil
}

// "[>] curl -XGET \"%s\", sz: %s [<] %d in %-8s, sz: %s"
func reqResLog(prefix string, method string, u *url.URL, req *http.Request, start time.Time, resp *http.Response,
	fn func(string, func(*bytes.Buffer)),
) {
	if fn == nil {
		return
	}

	fn(prefix, func(buf *bytes.Buffer) {
		buf.WriteString("[>]")
		buf.WriteByte(' ')
		buf.WriteString(`curl -X`)
		buf.WriteString(method)
		buf.WriteByte(' ')
		buf.WriteByte('"')
		buf.WriteString(u.String())
		buf.WriteByte('"')
		buf.WriteByte(',')
		buf.WriteByte(' ')
		buf.WriteString("sz:")
		buf.WriteByte(' ')
		buf.WriteString(humanize.Bytes(uint64(req.ContentLength)))
		buf.WriteByte(' ')
		buf.WriteString("[<]")
		buf.WriteByte(' ')
		buf.WriteString(strconv.FormatInt(int64(resp.StatusCode), 10))
		buf.WriteByte(' ')
		buf.WriteString("in")
		buf.WriteByte(' ')
		buf.WriteString(fmt.Sprintf("%-8s", time.Since(start).String()))
		buf.WriteByte(',')
		buf.WriteByte(' ')
		buf.WriteString("sz:")
		buf.WriteByte(' ')
		buf.WriteString(humanize.Bytes(uint64(resp.ContentLength)))
	})
}

func doReqRes(ctx context.Context,
	c *http.Client, method string, u *url.URL,

	reqBodyData interface{}, respBody io.Writer,

	patchReqFn func(*http.Request),
	checkRespStatusCodeFn func(int) bool,

	logReqFn func(string, *http.Request, io.Reader),
	logReqResFn func(string, func(*bytes.Buffer)),
	logRespFn func(string, *http.Response, io.Writer),
) error {

	var (
		reqBody        io.Reader = nil
		reqContentType string    = ""
		start                    = time.Now()
		prefix                   = randStr.Gen(15)
	)

	if reqBodyData != nil {
		if r, ok := reqBodyData.(io.Reader); ok {
			reqBody = r
		} else {
			buf := bufPool.NewBuf()
			defer buf.Release()

			// TODO: move json/yaml/protobuf/xml fns to codec package; to lib/
			if body, e := json.Marshal(reqBodyData); e != nil {
				return errors.Wrap(ErrRequestEncode, e.Error())
			} else {
				buf.Write(body)

				reqBody = buf
				reqContentType = "application/json" // TODO: move to lib/
			}
		}
	}

	req, e := http.NewRequest(method, u.String(), reqBody)
	if e != nil {
		return errors.Wrap(ErrRequestCreate, e.Error())
	}
	req = req.WithContext(ctx)
	if v := reqContentType; v != "" {
		req.Header.Set("Content-Type", v)
	}

	if fn := patchReqFn; fn != nil {
		fn(req)
	}

	if fn := patchReqFn; fn != nil {
		fn(req)
	}

	resp, e := c.Do(req)
	if e != nil {
		return errors.Wrap(ErrRequestExecute, e.Error())
	}
	defer resp.Body.Close()

	reqResLog(prefix, method, u, req, start, resp, logReqResFn)

	if respBody != nil {
		if _, e := io.Copy(respBody, resp.Body); e != nil {
			return errors.Wrap(ErrResponseRead, e.Error())
		}

		if fn := logRespFn; fn != nil {
			fn(prefix, resp, respBody)
		}
	}

	if fn := checkRespStatusCodeFn; fn != nil && !fn(resp.StatusCode) || fn == nil && resp.StatusCode >= 300 {
		return errors.Wrap(BadStatusCode, fmt.Sprintf("(:code %d)", resp.StatusCode))
	}

	return nil
}
