package httpx

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
)

func loopbackCertificate(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-time.Minute),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(parsed)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: privateKey}, roots
}

func TestStartWithCertificateUsesTLS(t *testing.T) {
	certificate, roots := loopbackCertificate(t)
	service := NewService(GetDefaultConfig("127.0.0.1:0", "tls-regression"))
	if service == nil {
		t.Fatal("NewService returned nil")
	}
	app := service.GetEngine()
	app.Get("/health", func(c fiber.Ctx) error { return c.SendString("secure") })
	ready := make(chan struct{}, 1)
	app.Hooks().OnListen(func(fiber.ListenData) error {
		ready <- struct{}{}
		return nil
	})
	stopped := service.StartWithCertificate(certificate, nil)
	t.Cleanup(func() {
		if err := service.Shutdown(); err != nil {
			t.Error(err)
		}
		select {
		case err := <-stopped:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(5 * time.Second):
			t.Error("TLS listener did not shut down")
		}
	})
	select {
	case <-ready:
	case err := <-stopped:
		t.Fatalf("listener failed before readiness: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("listener did not become ready")
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}}
	t.Cleanup(transport.CloseIdleConnections)
	client := &http.Client{Transport: transport, Timeout: 3 * time.Second}
	response, err := client.Get("https://" + service.GetServerAddress() + "/health")
	if err != nil {
		t.Errorf("certificate listener did not provide HTTPS: %v", err)
	} else {
		body, readErr := io.ReadAll(response.Body)
		response.Body.Close()
		if readErr != nil || response.StatusCode != http.StatusOK || string(body) != "secure" || response.TLS == nil {
			t.Errorf("TLS response status=%d body=%q tls=%v err=%v", response.StatusCode, body, response.TLS != nil, readErr)
		}
	}
	plainTransport := &http.Transport{}
	t.Cleanup(plainTransport.CloseIdleConnections)
	plainClient := &http.Client{Transport: plainTransport, Timeout: 3 * time.Second}
	if response, err := plainClient.Get("http://" + service.GetServerAddress() + "/health"); err == nil {
		defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			t.Error("certificate listener accepted plaintext HTTP")
		}
	}
}
