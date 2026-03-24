package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/license"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "genkey":
		if err := handleGenKey(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "sign":
		if err := handleSign(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("FluxStream license tool")
	fmt.Println("  genkey -out-dir ./license-keys")
	fmt.Println("  sign -key private.pem -customer 'Acme' -valid-until 2027-12-31 -features abr,rtmps -out license.json")
}

func handleGenKey(args []string) error {
	fs := flag.NewFlagSet("genkey", flag.ExitOnError)
	outDir := fs.String("out-dir", ".", "output directory")
	prefix := fs.String("prefix", "fluxstream-license", "file prefix")
	fs.Parse(args)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*outDir, 0755); err != nil {
		return err
	}
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}
	privPath := filepath.Join(*outDir, *prefix+".private.pem")
	pubPath := filepath.Join(*outDir, *prefix+".public.pem")
	if err := os.WriteFile(privPath, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER}), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(pubPath, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}), 0644); err != nil {
		return err
	}
	fmt.Println("private:", privPath)
	fmt.Println("public:", pubPath)
	return nil
}

func handleSign(args []string) error {
	fs := flag.NewFlagSet("sign", flag.ExitOnError)
	keyPath := fs.String("key", "", "private key PEM path")
	outPath := fs.String("out", "license.json", "output license path")
	product := fs.String("product", "FluxStream", "product")
	licenseID := fs.String("license-id", fmt.Sprintf("lic_%d", time.Now().Unix()), "license id")
	customer := fs.String("customer", "", "customer name")
	email := fs.String("email", "", "contact email")
	company := fs.String("company", "", "company")
	issuedAt := fs.String("issued-at", time.Now().UTC().Format("2006-01-02"), "YYYY-MM-DD or RFC3339")
	validUntil := fs.String("valid-until", time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02"), "YYYY-MM-DD or RFC3339")
	maintenanceUntil := fs.String("maintenance-until", time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02"), "YYYY-MM-DD or RFC3339")
	features := fs.String("features", "abr,rtmps,recording,branding", "comma separated features")
	maxNodes := fs.Int("max-nodes", 1, "max nodes")
	notes := fs.String("notes", "", "notes")
	fs.Parse(args)

	if strings.TrimSpace(*keyPath) == "" {
		return fmt.Errorf("-key gerekli")
	}
	if strings.TrimSpace(*customer) == "" {
		return fmt.Errorf("-customer gerekli")
	}
	keyData, err := os.ReadFile(*keyPath)
	if err != nil {
		return err
	}
	priv, err := license.ParsePrivateKeyPEM(keyData)
	if err != nil {
		return err
	}
	doc := license.Document{
		Product:          *product,
		LicenseID:        *licenseID,
		Customer:         *customer,
		Email:            *email,
		Company:          *company,
		IssuedAt:         *issuedAt,
		ValidUntil:       *validUntil,
		MaintenanceUntil: *maintenanceUntil,
		MaxNodes:         *maxNodes,
		Features:         splitCSV(*features),
		Notes:            *notes,
	}
	if err := license.Sign(&doc, priv); err != nil {
		return err
	}
	payload, err := jsonPretty(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(*outPath, payload, 0644)
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func jsonPretty(doc license.Document) ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}
