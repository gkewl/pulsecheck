package sse

import (
	"fmt"
	"net/url"
)

func makeOrigin(url *url.URL) string {
	return fmt.Sprintf("%s://%s", url.Scheme, url.Host)
}

func chomp(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\n' {
		return data[:len(data)-2]
	} else {
		return data
	}
}
