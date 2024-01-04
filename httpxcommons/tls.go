package httpxcommons

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/coffeehc/base/errors"
)

func LoadCertificate(raw []byte) (tls.Certificate, error) {
	var cert tls.Certificate
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		var err error
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		} else {
			cert.PrivateKey, err = parsePrivateKey(block.Bytes)
			if err != nil {
				return cert, err
			}
		}
		raw = rest
	}

	if len(cert.Certificate) == 0 {
		return cert, errors.SystemError("没有公钥证书")
	} else if cert.PrivateKey == nil {
		return cert, errors.SystemError("没有私钥证书")
	}
	return cert, nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.SystemError("加载PKCS8私钥失败")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, errors.SystemError("无法解析私钥")
}
