package license

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAjU1UBhgEkMqXUxHvJfL5LD+IIcTwtgowXWhTS2ieb2Y=
-----END PUBLIC KEY-----
`

type Document struct {
	Product          string   `json:"product"`
	LicenseID        string   `json:"license_id"`
	Customer         string   `json:"customer"`
	Email            string   `json:"email,omitempty"`
	Company          string   `json:"company,omitempty"`
	IssuedAt         string   `json:"issued_at"`
	ValidUntil       string   `json:"valid_until"`
	MaintenanceUntil string   `json:"maintenance_until,omitempty"`
	MaxNodes         int      `json:"max_nodes,omitempty"`
	Features         []string `json:"features,omitempty"`
	Notes            string   `json:"notes,omitempty"`
	Signature        string   `json:"signature"`
}

type payload struct {
	Product          string   `json:"product"`
	LicenseID        string   `json:"license_id"`
	Customer         string   `json:"customer"`
	Email            string   `json:"email,omitempty"`
	Company          string   `json:"company,omitempty"`
	IssuedAt         string   `json:"issued_at"`
	ValidUntil       string   `json:"valid_until"`
	MaintenanceUntil string   `json:"maintenance_until,omitempty"`
	MaxNodes         int      `json:"max_nodes,omitempty"`
	Features         []string `json:"features,omitempty"`
	Notes            string   `json:"notes,omitempty"`
}

type Status struct {
	Mode             string   `json:"mode"`
	Valid            bool     `json:"valid"`
	Message          string   `json:"message"`
	PublicKeySource  string   `json:"public_key_source"`
	Product          string   `json:"product,omitempty"`
	LicenseID        string   `json:"license_id,omitempty"`
	Customer         string   `json:"customer,omitempty"`
	Email            string   `json:"email,omitempty"`
	Company          string   `json:"company,omitempty"`
	IssuedAt         string   `json:"issued_at,omitempty"`
	ValidUntil       string   `json:"valid_until,omitempty"`
	MaintenanceUntil string   `json:"maintenance_until,omitempty"`
	MaxNodes         int      `json:"max_nodes,omitempty"`
	Features         []string `json:"features,omitempty"`
	UsingEmbeddedKey bool     `json:"using_embedded_key"`
}

type Manager struct {
	dataDir string
}

func NewManager(dataDir string) *Manager {
	return &Manager{dataDir: dataDir}
}

func (m *Manager) LicenseDir() string {
	return filepath.Join(m.dataDir, "license")
}

func (m *Manager) LicensePath() string {
	return filepath.Join(m.LicenseDir(), "license.json")
}

func (m *Manager) PublicKeyPath() string {
	return filepath.Join(m.LicenseDir(), "public.pem")
}

func (m *Manager) EnsureDirs() error {
	return os.MkdirAll(m.LicenseDir(), 0755)
}

func (m *Manager) SaveLicense(raw string) error {
	if err := m.EnsureDirs(); err != nil {
		return err
	}
	var doc Document
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return err
	}
	pretty, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.LicensePath(), pretty, 0644)
}

func (m *Manager) SavePublicKey(raw string) error {
	if err := m.EnsureDirs(); err != nil {
		return err
	}
	if _, err := ParsePublicKeyPEM([]byte(raw)); err != nil {
		return err
	}
	return os.WriteFile(m.PublicKeyPath(), []byte(strings.TrimSpace(raw)+"\n"), 0644)
}

func (m *Manager) Status(now time.Time) Status {
	pub, source, embedded, err := m.loadPublicKey()
	if err != nil {
		return Status{Mode: "invalid", Message: err.Error(), PublicKeySource: source, UsingEmbeddedKey: embedded}
	}
	data, err := os.ReadFile(m.LicensePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Status{Mode: "unlicensed", Message: "Lisans yuklenmedi", PublicKeySource: source, UsingEmbeddedKey: embedded}
		}
		return Status{Mode: "invalid", Message: err.Error(), PublicKeySource: source, UsingEmbeddedKey: embedded}
	}
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return Status{Mode: "invalid", Message: "Lisans JSON okunamadi", PublicKeySource: source, UsingEmbeddedKey: embedded}
	}
	if err := Verify(doc, pub); err != nil {
		return Status{Mode: "invalid", Message: err.Error(), PublicKeySource: source, UsingEmbeddedKey: embedded}
	}
	status := Status{
		Mode:             "active",
		Valid:            true,
		Message:          "Lisans dogrulandi",
		PublicKeySource:  source,
		Product:          doc.Product,
		LicenseID:        doc.LicenseID,
		Customer:         doc.Customer,
		Email:            doc.Email,
		Company:          doc.Company,
		IssuedAt:         doc.IssuedAt,
		ValidUntil:       doc.ValidUntil,
		MaintenanceUntil: doc.MaintenanceUntil,
		MaxNodes:         doc.MaxNodes,
		Features:         append([]string(nil), doc.Features...),
		UsingEmbeddedKey: embedded,
	}
	if t, err := parseDate(doc.ValidUntil); err == nil && now.After(t) {
		status.Mode = "expired"
		status.Valid = false
		status.Message = "Lisans suresi doldu"
	}
	return status
}

func (m *Manager) loadPublicKey() (ed25519.PublicKey, string, bool, error) {
	if data, err := os.ReadFile(m.PublicKeyPath()); err == nil {
		pub, err := ParsePublicKeyPEM(data)
		return pub, m.PublicKeyPath(), false, err
	}
	pub, err := ParsePublicKeyPEM([]byte(defaultPublicKeyPEM))
	return pub, "embedded-dev-key", true, err
}

func ParsePublicKeyPEM(data []byte) (ed25519.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("gecerli public key bulunamadi")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := parsed.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("desteklenmeyen public key tipi")
	}
	return pub, nil
}

func ParsePrivateKeyPEM(data []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("gecerli private key bulunamadi")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("desteklenmeyen private key tipi")
	}
	return priv, nil
}

func Sign(doc *Document, priv ed25519.PrivateKey) error {
	payloadBytes, err := canonicalPayload(*doc)
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, payloadBytes)
	doc.Signature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func Verify(doc Document, pub ed25519.PublicKey) error {
	if strings.TrimSpace(doc.Signature) == "" {
		return fmt.Errorf("imza eksik")
	}
	sig, err := base64.StdEncoding.DecodeString(doc.Signature)
	if err != nil {
		return fmt.Errorf("imza cozulmedi")
	}
	payloadBytes, err := canonicalPayload(doc)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, payloadBytes, sig) {
		return fmt.Errorf("lisans imzasi gecersiz")
	}
	if _, err := parseDate(doc.IssuedAt); err != nil {
		return fmt.Errorf("issued_at gecersiz")
	}
	if _, err := parseDate(doc.ValidUntil); err != nil {
		return fmt.Errorf("valid_until gecersiz")
	}
	if doc.MaintenanceUntil != "" {
		if _, err := parseDate(doc.MaintenanceUntil); err != nil {
			return fmt.Errorf("maintenance_until gecersiz")
		}
	}
	return nil
}

func canonicalPayload(doc Document) ([]byte, error) {
	features := append([]string(nil), doc.Features...)
	sort.Strings(features)
	return json.Marshal(payload{
		Product:          doc.Product,
		LicenseID:        doc.LicenseID,
		Customer:         doc.Customer,
		Email:            doc.Email,
		Company:          doc.Company,
		IssuedAt:         doc.IssuedAt,
		ValidUntil:       doc.ValidUntil,
		MaintenanceUntil: doc.MaintenanceUntil,
		MaxNodes:         doc.MaxNodes,
		Features:         features,
		Notes:            doc.Notes,
	})
}

func parseDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("bos tarih")
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", value)
}

func SampleDocument() Document {
	return Document{
		Product:          "FluxStream",
		LicenseID:        "lic_demo_001",
		Customer:         "Ornek Kurum",
		Email:            "it@example.com",
		Company:          "Ornek Kurum",
		IssuedAt:         time.Now().UTC().Format("2006-01-02"),
		ValidUntil:       time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02"),
		MaintenanceUntil: time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02"),
		MaxNodes:         1,
		Features:         []string{"abr", "rtmps", "recording", "branding"},
		Notes:            "Offline signed license sample",
	}
}
