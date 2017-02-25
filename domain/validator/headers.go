package validator

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// ValidHeader checks if a given response header has
// an expected value.
func ValidHeader(hostname, header, expected string) bool {
	v, err := readHeader(hostname, hostname, header)
	if err != nil {
		return false
	}

	return v == expected
}

// readHeader sends an HTTP HEAD request to the given address
// asking for the hostname information.
// It looks for the given header and returns the value of that header.
// It returns an error if the server doesn't reply with a 2xx or 3xx
// status code or the header is not present.
func readHeader(address, hostname, header string) (string, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   address,
	}
	req, err := http.NewRequest("HEAD", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Host = hostname

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 304 {
		return "", errors.Errorf("error requesting domain headers: %s - %s", hostname, resp.Status)
	}

	v := resp.Header.Get(header)
	if v == "" {
		return "", errors.Errorf("error looking up header: %s - %s", hostname, header)
	}

	return v, nil
}
