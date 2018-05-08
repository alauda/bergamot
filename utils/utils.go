package utils

import (
	"net/http"
	"net/url"
)

// GetURL combine endpoint and path
func GetURL(endpint, path string, values url.Values) (string, error) {
	endpointURL, err := url.Parse(endpint)
	if err != nil {
		return "", err
	}
	endpointURL.Path = path
	if values != nil {
		endpointURL.RawQuery = values.Encode()
	}

	return endpointURL.String(), nil
}

// CloseResponse safe close response
func CloseResponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
