package enroll

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// TLSConfig holds the CSR subject fields the server expects, as reported by
// its /Marti/api/tls/config endpoint.
type TLSConfig struct {
	Organization       string
	OrganizationalUnit string
}

type certificateConfigXML struct {
	NameEntries struct {
		NameEntry []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"nameEntry"`
	} `xml:"nameEntries"`
}

// FetchTLSConfig retrieves the CSR subject configuration from baseURL
// (e.g. "https://192.168.1.50:8446"). The endpoint requires no
// authentication and is queried before any trust is established with the
// server.
func FetchTLSConfig(ctx context.Context, httpClient *http.Client, baseURL string) (TLSConfig, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/Marti/api/tls/config", nil)
	if err != nil {
		return TLSConfig{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return TLSConfig{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TLSConfig{}, fmt.Errorf("enroll: fetching TLS config: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TLSConfig{}, err
	}

	var parsed certificateConfigXML
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return TLSConfig{}, fmt.Errorf("enroll: parsing TLS config response: %w", err)
	}

	cfg := TLSConfig{}
	for _, entry := range parsed.NameEntries.NameEntry {
		switch entry.Name {
		case "O":
			cfg.Organization = entry.Value
		case "OU":
			cfg.OrganizationalUnit = entry.Value
		}
	}

	return cfg, nil
}
