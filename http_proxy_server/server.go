package http_proxy_server

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	// Wrap specifies a function for optionally wrapping upstream for
	// inspecting the decrypted HTTP request and response.
	Wrap func(upstream http.Handler) http.Handler

	// CA specifies the root CA for generating leaf certs for each incoming
	// TLS request.
	serverCrt  tls.Certificate
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
	if r.Method == http.MethodConnect {
		p.serveConnect(w, r)
		return
	}

	rp := httputil.NewSingleHostReverseProxy(p.targetAddr)
	r.Host = r.URL.Host
	rp.ServeHTTP(w, r)
}

func (p *Proxy) serveConnect(w http.ResponseWriter, r *http.Request) {
	// sConfig := new(tls.Config)
	// if p.TLSServerConfig != nil {
	// 	*sConfig = *p.TLSServerConfig
	// }
	// sConfig.Certificates = []tls.Certificate{*provisionalCert}
	// sConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	// 	cConfig := new(tls.Config)
	// 	if p.TLSClientConfig != nil {
	// 		*cConfig = *p.TLSClientConfig
	// 	}
	// 	cConfig.ServerName = hello.ServerName
	// 	sconn, err = tls.Dial("tcp", r.Host, cConfig)
	// 	if err != nil {
	// 		log.Println("dial", r.Host, err)
	// 		return nil, err
	// 	}
	// 	return p.cert(hello.ServerName)
	// }
}

// func (p *Proxy) newServerConfig(w http.ResponseWriter, r *http.Request) tls.Config {
// 	srvConfig := tls.Config{
// 		Certificates: []tls.Certificate{
// 			p.serverCert,
// 		},
// 		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {

// 		}
// 	}
// }
