//go:build !desktop

package main

import "net/http"

func runDesktopApp(handler http.Handler) {
	// Web server mode - no desktop GUI required
}
