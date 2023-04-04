package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"mango-admin/pkg/sdk/api"
)

var (
	ResponseLog      io.Writer = os.Stdout
	slowMilliseconds int64     = 300 // slow api
)

func ApiLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer printLog(ctx, copyRequestBody(ctx), time.Now())

		ctx.Next()
	}
}

type apiLogFormat struct {
	Time       string          `json:"time"`
	Path       string          `json:"path"`
	Method     string          `json:"method"`
	Cost       int64           `json:"cost" `          // Milliseconds
	Slow       bool            `json:"slow,omitempty"` // Slow request
	Uid        int             `json:"uid,omitempty"`
	Query      string          `json:"query,omitempty"`
	BodyString string          `json:"body_string,omitempty"`
	BodyJson   json.RawMessage `json:"body_json,omitempty"`
	Response   json.RawMessage `json:"response"`
}

func copyRequestBody(ctx *gin.Context) []byte {
	if ctx.Request.Method == http.MethodGet {
		return nil
	}

	bodyBytes, _ := ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

func printLog(ctx *gin.Context, bodyData []byte, start time.Time) {
	log := api.GetRequestLogger(ctx)

	apiLoger := apiLogFormat{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Cost:   time.Since(start).Milliseconds(),
		Path:   ctx.Request.URL.Path,
		Method: ctx.Request.Method,
		Uid:    ctx.GetInt("uid"),
		Query:  ctx.Request.URL.RawQuery,
	}

	if len(bodyData) > 0 && bodyData[0] == '{' {
		apiLoger.BodyJson = bodyData
	} else {
		apiLoger.BodyString, _ = url.QueryUnescape(string(bodyData))
	}

	if apiLoger.Cost > slowMilliseconds {
		apiLoger.Slow = true
	}

	rt, bl := ctx.Get("result")
	if bl {
		rb, err := json.Marshal(rt)
		if err != nil {
			log.Warnf("json Marshal result error, %s", err.Error())
		} else {
			apiLoger.Response = rb
		}
	}

	b, err := json.Marshal(apiLoger)
	if err != nil {
		log.Error(err)
		return
	}
	ResponseLog.Write(b)
	ResponseLog.Write([]byte("\n"))
}
