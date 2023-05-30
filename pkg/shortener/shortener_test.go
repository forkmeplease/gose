// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package shortener_test

import (
	"net/url"
	"testing"

	"github.com/stv0g/gose/pkg/config"
	"github.com/stv0g/gose/pkg/shortener"
)

func TestShortener(t *testing.T) {
	s, err := shortener.NewShortener(&config.ShortenerConfig{
		Endpoint: "https://l.0l.de/rest/v2/short-urls/shorten?apiKey=952eee41-ad41-4743-8e7c-a1571168fd22&format=txt&longUrl={{.Url}}",
		Method:   "GET",
		Response: "raw",
	})
	if err != nil {
		t.FailNow()
	}

	long, err := url.Parse("http://a-very-long-url.com")
	if err != nil {
		t.FailNow()
	}

	short, err := s.Shorten(long)
	if err != nil {
		t.Fatalf("Failed to shorten: %s", err)
	}

	t.Logf("Shorted: %s", short)
}
