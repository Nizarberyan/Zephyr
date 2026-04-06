package services

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SubsonicService struct {
	BaseURL string
	Client  *http.Client
}

func NewSubsonicService(baseURL string) *SubsonicService {
	return &SubsonicService{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SubsonicResponse matches the standard Subsonic API response
type SubsonicResponse struct {
	XMLName xml.Name `xml:"subsonic-response"`
	Status  string   `xml:"status,attr"`
	Error   *struct {
		Code    int    `xml:"code,attr"`
		Message string `xml:"message,attr"`
	} `xml:"error,omitempty"`
}

// VerifySubsonic checks credentials against the Navidrome Subsonic API
func (s *SubsonicService) VerifySubsonic(username, password, token, salt string) (bool, error) {
	u, err := url.Parse(s.BaseURL + "/rest/ping.view")
	if err != nil {
		return false, err
	}

	q := u.Query()
	q.Set("u", username)
	q.Set("v", "1.16.1") // Standard Subsonic version
	q.Set("c", "Zephyr") // Client name

	if token != "" && salt != "" {
		q.Set("t", token)
		q.Set("s", salt)
	} else if password != "" {
		q.Set("p", password)
	} else {
		return false, fmt.Errorf("missing credentials (password or token/salt)")
	}

	u.RawQuery = q.Encode()

	resp, err := s.Client.Get(u.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var subResp SubsonicResponse
	if err := xml.NewDecoder(resp.Body).Decode(&subResp); err != nil {
		return false, err
	}

	if subResp.Status == "ok" {
		return true, nil
	}

	return false, nil
}
