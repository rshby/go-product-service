package middleware

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
)

var bufferPool = sync.Pool{New: func() any {
	logrus.Info("create new buffer from pool")
	return &bytes.Buffer{}
}}

type ResponseBodyWriter struct {
	ginResponseWriter gin.ResponseWriter
	body              *bytes.Buffer
}

func NewResponseBodyWriter(ginResponseWriter gin.ResponseWriter) *ResponseBodyWriter {
	return &ResponseBodyWriter{
		ginResponseWriter: ginResponseWriter,
	}
}

func (r *ResponseBodyWriter) Header() http.Header {
	return r.ginResponseWriter.Header()
}

func (r *ResponseBodyWriter) WriteHeader(statusCode int) {
	r.ginResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseBodyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ginResponseWriter.Hijack()
}

func (r *ResponseBodyWriter) Flush() {
	r.ginResponseWriter.Flush()
}

func (r *ResponseBodyWriter) CloseNotify() <-chan bool {
	return r.ginResponseWriter.CloseNotify()
}

func (r *ResponseBodyWriter) Status() int {
	return r.ginResponseWriter.Status()
}

func (r *ResponseBodyWriter) Size() int {
	return r.ginResponseWriter.Size()
}

func (r *ResponseBodyWriter) WriteString(s string) (int, error) {
	return r.ginResponseWriter.WriteString(s)
}

func (r *ResponseBodyWriter) Written() bool {
	return r.ginResponseWriter.Written()
}

func (r *ResponseBodyWriter) WriteHeaderNow() {
	r.ginResponseWriter.WriteHeaderNow()
}

func (r *ResponseBodyWriter) Pusher() http.Pusher {
	return r.ginResponseWriter.Pusher()
}

func (r *ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body = bufferPool.Get().(*bytes.Buffer)

	// write to body
	r.body.Write(b)

	// also write to gin respons body
	return r.ginResponseWriter.Write(b)
}

func (r *ResponseBodyWriter) PutBack() {
	if r.body == nil {
		return
	}

	r.body.Reset()
	bufferPool.Put(r.body)
}
