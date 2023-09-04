package emailoctopus

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	eoapi = "https://emailoctopus.com/api/1.6/"
)

type EmailOctopus struct {
	ListID string
	APIKey string
}

type eoCreateContactParams struct {
	APIKey string            `json:"api_key"`
	Email  string            `json:"email_address"`
	Status string            `json:"status"`
	Fields map[string]string `json:"fields"`
}

// https://emailoctopus.com/api-documentation/lists/create-contact
func (eo *EmailOctopus) Add(ctx context.Context, email string) error {
	params := eoCreateContactParams{
		APIKey: eo.APIKey,
		Email:  email,
		Status: "SUBSCRIBED",
	}
	p, err := json.Marshal(params)
	if err != nil {
		return Errorc(err)
	}

	u := fmt.Sprintf(eoapi+"lists/%s/contacts", eo.ListID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(p))
	if err != nil {
		return Errorc(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Errorc(err)
	}
	resp.Body.Close()

	return nil
}

// https://emailoctopus.com/api-documentation/lists/create-contact
func (eo *EmailOctopus) Remove(ctx context.Context, email string) error {
	id := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(email))))

	u := fmt.Sprintf(eoapi+"lists/%s/contacts/%s?api_key=%s",
		eo.ListID, id, eo.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return Errorc(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Errorc(err)
	}
	resp.Body.Close()

	return nil
}
