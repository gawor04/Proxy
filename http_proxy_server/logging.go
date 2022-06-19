package http_proxy_server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

type responseWriterWrapper struct {
	w          http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (w *responseWriterWrapper) Write(buf []byte) (int, error) {
	w.body.Write(buf)
	return w.w.Write(buf)
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriterWrapper) String() string {
	var buf bytes.Buffer

	for k, v := range w.w.Header() {
		buf.WriteString(fmt.Sprintf("\n%s: %v", k, v))
	}

	buf.WriteString(fmt.Sprintf("\nStatus Code: %d", w.statusCode))

	bodyStr := w.body.String()
	if len(bodyStr) > 0 {
		buf.WriteString("\nBody: ")
		buf.WriteString(bodyStr)
	}

	buf.WriteString("\n\n")

	return buf.String()
}

func logMiddleware(proxy *Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("New HTTP request")
		req_dump, err := httputil.DumpRequest(r, true)
		if err == nil {
			log.Info(string(req_dump))
		} else {
			log.Errorf("DumpRequest: %v", err)
		}

		wrapper := responseWriterWrapper{w: w}
		proxy.ServeHTTP(&wrapper, r)

		log.Info("New HTTP response")
		log.Info(wrapper.String())
	}
}
