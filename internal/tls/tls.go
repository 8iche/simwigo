package tls

//Source from https://go.dev/src/crypto/tls/generate_cert.go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func GenerateTLSCertificate(host string, dir string, isRSA bool) (certFile string, keyFile string, err error) {

	if len(host) == 0 {
		return "", "", errors.New("invalid host(s)")
	}

	var rsaBits = 4096
	var priv interface{}

	certFile = dir + "cert.pem"
	keyFile = dir + "key.pem"

	switch isRSA {
	case true:
		priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	default:
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		//log.Fatalf("Unrecognized elliptic curve: %q", *ecdsaCurve)
	}
	if err != nil {
		//log.Fatalf("Failed to generate private key: %v", err)
		return "", "", err
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	notBefore := time.Now()

	notAfter := notBefore.Add(1000)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		//log.Fatalf("Failed to generate serial number: %v", err)
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Simwigo"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		//log.Fatalf("Failed to create certificate: %v", err)
		return "", "", err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		//log.Fatalf("Failed to open cert.pem for writing: %v", err)
		return "", "", err
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		//log.Fatalf("Failed to write data to cert.pem: %v", err)
		return "", "", err
	}
	if err := certOut.Close(); err != nil {
		//log.Fatalf("Error closing cert.pem: %v", err)
		return "", "", err
	}
	//log.Print("wrote cert.pem\n")

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		//log.Fatalf("Failed to open key.pem for writing: %v", err)
		return "", "", err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		//log.Fatalf("Unable to marshal private key: %v", err)
		return "", "", err
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		//log.Fatalf("Failed to write data to key.pem: %v", err)
		return "", "", err
	}
	if err := keyOut.Close(); err != nil {
		//log.Fatalf("Error closing key.pem: %v", err)
		return "", "", err
	}
	//log.Print("wrote key.pem\n")
	return certFile, keyFile, nil
}
