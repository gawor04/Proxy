package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"proxy/cert"
	"proxy/config"
	"proxy/http_proxy_server"

	log "github.com/sirupsen/logrus"
)

func runHttpsServer(cfg config.Config, handler http.Handler) {
	userName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	host, port, err := net.SplitHostPort(cfg.Listen)
	if err != nil {
		host = cfg.Listen
	}

	c := cert.NewCert(userName, host, cfg.CertDir)
	certOut, err := c.CreateCert()
	if err != nil {
		log.Fatalf("Create cert", err)
	}
	certAbsPath, err := filepath.Abs(certOut.CaCertPath)
	if err != nil {
		log.Fatalf("Getting CA certificate abs path", err)
	}
	log.Infof("\nImport CA certificate to your web browser:\n - %s", certAbsPath)

	if port == "" {
		port = "443"
	}
	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", port),
		TLSConfig: &tls.Config{ServerName: host},
		Handler:   handler,
	}

	err = server.ListenAndServeTLS(certOut.ServerCertPath, certOut.ServerKeyPath)
	if err != nil {
		log.Fatalf("ListenAndServeTLS: %v", err)
	}
}

func main() {
	cfg, err := config.LoadConfig(os.Args[1:])
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	/* log to stdout and log file */
	file, err := os.OpenFile(cfg.Log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Open log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))

	/* do not ignore new lines */
	formatter := &log.TextFormatter{}
	formatter.DisableQuote = true
	log.SetFormatter(formatter)

	handler, err := http_proxy_server.NewHttpProxyHandler(cfg.Target)
	if err != nil {
		log.Fatalf("NewHttpProxyHandler: %v", err)
	}

	log.Infof("Start proxy: %s -> %s", cfg.Listen, cfg.Target)

	if cfg.SslStriping {
		runHttpsServer(cfg, handler)
	} else {
		_, port, err := net.SplitHostPort(cfg.Listen)
		if err != nil {
			port = "80"
		}
		server := &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: handler,
		}
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}
}
