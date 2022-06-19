package http_proxy_server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	targetAddr *url.URL
}

func NewHttpProxyHandler(targetAddr string) (http.Handler, error) {
	targetUrl, err := url.Parse(targetAddr)
	if err != nil {
		return nil, err
	}

	return logMiddleware(&Proxy{
		targetAddr: targetUrl,
	}), err
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rp := httputil.NewSingleHostReverseProxy(p.targetAddr)
	r.Host = r.URL.Host
	rp.ServeHTTP(w, r)
}
