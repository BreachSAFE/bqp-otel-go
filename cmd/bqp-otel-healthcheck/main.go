// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

// Command bqp-otel-healthcheck probes the Collector health_check extension.
package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	client := http.Client{Timeout: 2 * time.Second}
	response, err := client.Get("http://127.0.0.1:13133/")
	if err != nil {
		os.Exit(1)
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1024))
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
