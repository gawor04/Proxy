package cert

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/square/certstrap/pkix"
)

var (
	certExpires = time.Now().Add(time.Hour * time.Duration(36*30*24))
	keyBits     = 4096
)

type cert struct {
	commonName     string
	domain         string
	caCertPath     string
	caKeyPath      string
	serverCertPath string
	serverKeyPath  string
}

type CreateCertOutput struct {
	CaCertPath     string
	ServerCertPath string
	ServerKeyPath  string
}

func NewCert(commonName, domain, caDir string) *cert {
	return &cert{
		commonName:     commonName,
		domain:         domain,
		caCertPath:     path.Join(caDir, fmt.Sprintf("%sCA.crt", commonName)),
		caKeyPath:      path.Join(caDir, fmt.Sprintf("%sCA.key", commonName)),
		serverCertPath: path.Join(caDir, fmt.Sprintf("%s.crt", domain)),
		serverKeyPath:  path.Join(caDir, fmt.Sprintf("%s.key", domain)),
	}
}

func (c *cert) CreateCert() (*CreateCertOutput, error) {
	caKey, caCrt, err := c.initCa()
	if err != nil {
		return nil, err
	}

	_, serverCsr, err := c.requestServerCert()
	if err != nil {
		return nil, err
	}

	_, err = c.signByCa(caCrt, caKey, serverCsr)
	if err != nil {
		return nil, err
	}

	return &CreateCertOutput{
		CaCertPath:     c.caCertPath,
		ServerCertPath: c.serverCertPath,
		ServerKeyPath:  c.serverKeyPath,
	}, nil
}

func (c *cert) initCa() (key *pkix.Key, crt *pkix.Certificate, err error) {
	key, err = c.readOrCreateCaKey()
	if err != nil {
		return
	}
	crt, err = c.readOrCreateCaCertificate(key)
	if err != nil {
		return
	}

	return
}

func (c *cert) readOrCreateCaKey() (key *pkix.Key, err error) {
	keyBytes, err := os.ReadFile(c.caKeyPath)

	if err == nil {
		key, err = pkix.NewKeyFromPrivateKeyPEM(keyBytes)
		if err != nil {
			return
		}
	} else if os.IsNotExist(err) {
		key, err = pkix.CreateRSAKey(keyBits)
		if err != nil {
			return
		}
		keyBytes, err = key.ExportPrivate()
		if err != nil {
			return
		}
		if err = os.WriteFile(c.caKeyPath, keyBytes, 0644); err != nil {
			return
		}
	} else {
		return
	}

	if key == nil {
		err = fmt.Errorf("CA key is nil")
		return
	}

	return
}

func (c *cert) readOrCreateCaCertificate(key *pkix.Key) (crt *pkix.Certificate, err error) {
	crtBytes, err := os.ReadFile(c.caCertPath)

	if err == nil {
		crt, err = pkix.NewCertificateFromPEM(crtBytes)
		if err != nil {
			return
		}
	} else if os.IsNotExist(err) {
		crt, err = pkix.CreateCertificateAuthority(key, "", certExpires, "", "", "", "", c.commonName)
		if err != nil {
			return
		}
		crtBytes, err = crt.Export()
		if err != nil {
			return
		}
		if err = os.WriteFile(c.caCertPath, crtBytes, 0644); err != nil {
			return
		}
	} else {
		return
	}

	if crt == nil {
		err = fmt.Errorf("CA certificate is nil")
		return
	}

	return
}

func (c *cert) requestServerCert() (key *pkix.Key, csr *pkix.CertificateSigningRequest, err error) {
	key, err = pkix.CreateRSAKey(keyBits)
	if err != nil {
		return
	}
	csr, err = pkix.CreateCertificateSigningRequest(key, "", nil, []string{c.domain}, nil, "", "", "", "", c.domain)
	if err != nil {
		return
	}

	var keyBytes []byte
	keyBytes, err = key.ExportPrivate()
	if err != nil {
		return
	}
	if err = os.WriteFile(c.serverKeyPath, keyBytes, 0644); err != nil {
		return
	}

	if key == nil {
		err = fmt.Errorf("Server key is nil")
		return
	}
	if csr == nil {
		err = fmt.Errorf("Server certificate signing request is nil")
		return
	}

	return
}

func (c *cert) signByCa(caCrt *pkix.Certificate, caKey *pkix.Key, csr *pkix.CertificateSigningRequest) (crt []byte, err error) {
	serverCrt, err := pkix.CreateCertificateHost(caCrt, caKey, csr, certExpires)
	if err != nil {
		return
	}
	crt, err = serverCrt.Export()
	if err != nil {
		return
	}
	if err = os.WriteFile(c.serverCertPath, crt, 0644); err != nil {
		return
	}

	return
}
