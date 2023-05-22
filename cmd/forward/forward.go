package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/udhos/forward/cmd/forward/zlog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gopkg.in/yaml.v3"
)

func forward(c *gin.Context, app *application) {

	const me = "forward"

	ctxNew, span := newSpanGin(c, me, app)
	if span != nil {
		defer span.End()
	}

	method := c.Request.Method

	//
	// parse body
	//

	type errorResponse struct {
		Error string `json:"error"`
	}

	type requestBody struct {
		URL        string            `yaml:"url"`
		Method     string            `yaml:"method"`
		Body       string            `yaml:"body"`
		SetHeaders map[string]string `yaml:"set_headers"`
	}

	dec := yaml.NewDecoder(c.Request.Body)
	var in requestBody
	errYaml := dec.Decode(&in)
	if errYaml != nil {
		zlog.CtxErrorf(ctxNew, "%s: %v", me, errYaml)
		traceError(span, errYaml.Error())
		c.JSON(http.StatusBadRequest, errorResponse{Error: errYaml.Error()})
		return
	}

	log.Printf("%s: %s %#v", me, method, in)

	if in.Method != "" {
		method = in.Method
	}

	//
	// send request
	//

	req, errReq := http.NewRequestWithContext(ctxNew, method, in.URL, bytes.NewBufferString(in.Body))
	if errReq != nil {
		zlog.CtxErrorf(ctxNew, "%s: %v", me, errReq)
		traceError(span, errReq.Error())
		c.JSON(http.StatusBadRequest, errorResponse{Error: errReq.Error()})
		return
	}

	// set request headers
	for k, v := range in.SetHeaders {
		req.Header.Set(k, v)
	}

	client := httpClient()

	resp, errDo := client.Do(req)
	if errDo != nil {
		zlog.CtxErrorf(ctxNew, "%s: %v", me, errDo)
		traceError(span, errDo.Error())
		c.JSON(http.StatusBadRequest, errorResponse{Error: errDo.Error()})
		return
	}

	defer resp.Body.Close()

	//
	// read response body
	//

	full, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		zlog.CtxErrorf(ctxNew, "%s: %v", me, errBody)
		traceError(span, errBody.Error())
		c.JSON(http.StatusBadGateway, errorResponse{Error: errBody.Error()})
		return
	}

	//
	// copy response headers
	//

	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}

	//
	// send response
	//

	contentType := resp.Header.Get("content-type")

	c.Data(resp.StatusCode, contentType, full)
}

func httpClient() http.Client {
	return http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
}
