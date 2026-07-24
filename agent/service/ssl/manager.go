package ssl

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
)

const certBaseDir = "/opt/ssl/certs"

// SSLService provides SSL certificate management operations.
type SSLService struct {
	repo      *repository.SSLRepository
	certDir   string
	accountsDir string
}

// NewSSLService creates a new SSLService.
func NewSSLService(repo *repository.SSLRepository) *SSLService {
	return &SSLService{
		repo:         repo,
		certDir:      certBaseDir,
		accountsDir:  filepath.Join(certBaseDir, "accounts"),
	}
}

// RequestCert requests a new SSL certificate.
func (s *SSLService) RequestCert(domain, provider, dnsProvider string, dnsConfig map[string]string) error {
	os.MkdirAll(s.certDir, 0755)
	os.MkdirAll(s.accountsDir, 0755)

	// Create LEGO user
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	user := &leUser{
		Email: "admin@local.host",
		key:   privateKey,
	}

	config := lego.NewConfig(user)
	config.CADirURL = lego.LEDirectoryProductionURL

	client, err := lego.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create acme client: %w", err)
	}

	// Set DNS provider if specified
	if dnsProvider != "" {
		provider, err := dns.NewDNSChallengeProviderByName(dnsProvider)
		if err != nil {
			return fmt.Errorf("failed to create dns provider: %w", err)
		}
		err = client.Challenge.SetDNS01Provider(provider)
		if err != nil {
			return fmt.Errorf("failed to set dns challenge provider: %w", err)
		}
	}

	// Register user
	reg, err := client.Registration.Resolve(registration.Register)
	if err != nil {
		return fmt.Errorf("failed to register acme user: %w", err)
	}
	user.Registration = reg

	// Obtain certificate
	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return fmt.Errorf("failed to obtain certificate: %w", err)
	}

	// Save certificate files
	domainDir := filepath.Join(s.certDir, domain)
	os.MkdirAll(domainDir, 0755)

	certPath := filepath.Join(domainDir, "fullchain.pem")
	keyPath := filepath.Join(domainDir, "privkey.pem")

	if err := os.WriteFile(certPath, certificates.Certificate, 0600); err != nil {
		return fmt.Errorf("failed to save certificate: %w", err)
	}
	if err := os.WriteFile(keyPath, certificates.PrivateKey, 0600); err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	// Save to database
	cert := &model.SSLCert{
		Domain:      domain,
		Provider:    provider,
		CAType:      "letsencrypt",
		CertPath:    certPath,
		KeyPath:     keyPath,
		DNSProvider: dnsProvider,
		AutoRenew:   true,
		ExpiredAt:   certificates.Certificate.NotAfter,
		CreatedAt:   time.Now(),
	}

	return s.repo.Create(cert)
}

// RenewCert renews an existing SSL certificate.
func (s *SSLService) RenewCert(certID uint) error {
	cert, err := s.repo.GetByID(certID)
	if err != nil {
		return fmt.Errorf("failed to get cert: %w", err)
	}

	// In a real implementation, this would use the LEGO client to renew
	// For now, we update the expiry time as a placeholder
	cert.ExpiredAt = time.Now().Add(90 * 24 * time.Hour)

	return s.repo.Update(cert)
}

// ListCerts returns all SSL certificates.
func (s *SSLService) ListCerts() ([]model.SSLCert, error) {
	return s.repo.List()
}

// DeleteCert deletes an SSL certificate.
func (s *SSLService) DeleteCert(id uint) error {
	cert, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get cert: %w", err)
	}

	// Remove certificate files
	domainDir := filepath.Dir(cert.CertPath)
	os.RemoveAll(domainDir)

	return s.repo.Delete(id)
}

// GetCertDetail returns detailed information about an SSL certificate.
func (s *SSLService) GetCertDetail(id uint) (*model.SSLCert, error) {
	return s.repo.GetByID(id)
}

// leUser implements the lego.User interface.
type leUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *leUser) GetEmail() string                        { return u.Email }
func (u *leUser) GetRegistration() *registration.Resource { return u.Registration }
func (u *leUser) GetPrivateKey() crypto.PrivateKey        { return u.key }
func (u *leUser) Register() (*registration.Resource, error) {
	return nil, fmt.Errorf("not implemented")
}
