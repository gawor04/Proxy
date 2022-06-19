# http-proxy

## Install

```
go install github.com/gawor04/proxy/cmd/http-proxy@latest
```

## Usage

```
Usage of http-proxy:
  -config string
    	configuration file path, other arguments can be specified in yaml configuration file
  -listen string
    	proxy server listen address (default: 0.0.0.0:80, if ssl-striping used: 0.0.0.0:443)
  -target string
    	proxy target address (required)
  -log string
    	log file path (optional)
  -cert-dir
    	CA certificate and private key directory (default: ./)
  -ssl-striping
    	HTTPS connection between client and proxy server (default:`s false)
```

## Example configuration - HTTP:

```
---
listen: localhost:80
target: http://cat-fact.herokuapp.com
log: ./log.txt
ssl-striping: false
cert-dir: ./
```

## Example configuration - HTTPS:

```
---
listen: localhost:443
target: https://cat-fact.herokuapp.com
log: ./log.txt
ssl-striping: true
cert-dir: ./
```

also CA certificate created by http-proxy should be imported to web browser (cert path is listed in logs)