package utils

import (
	"context"
	"io"
	"net/http"
)

func Fetch(ctx context.Context, link string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return "", Errorc(err)
	}

	data, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", Errorc(err)
	}

	body, err := io.ReadAll(data.Body)
	if err != nil {
		return "", Errorc(err)
	}

	return string(body), nil
}
