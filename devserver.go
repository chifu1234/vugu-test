//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/vugu/vugu/devutil"
)

// Folder matches requests that access paths under "/api".
var Folder = devutil.RequestMatcherFunc(func(r *http.Request) bool {
	// Clean and ensure the path starts with "/".
	cleanPath := path.Clean("/" + r.URL.Path)

	// Check if the path starts with "/api/".
	// The trailing slash is important to ensure it matches the directory concept.
	// Also, check if the path exactly equals "/api" to include it as well.
	return strings.HasPrefix(cleanPath, "/api/") || cleanPath == "/api"
})

type proxy struct {
}

func (p proxy) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Parse the destination URL, where you want to proxy the requests to
	if req.Method == "OPTIONS" {
		return
	}
	targetURL, _ := url.Parse("http://localhost:8001")

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Update the headers to allow for SSL redirection
	req.URL.Host = targetURL.Host
	req.URL.Scheme = targetURL.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = targetURL.Host

	// Note: You might want to adjust the request before forwarding it,
	// For example, by modifying the req.Header here.
	// ServeHTTP uses the reverse proxy to forward the request and writes the response back to the original client

	proxy.ServeHTTP(resp, req)

}
func main() {
	l := "127.0.0.1:8844"
	log.Printf("Starting HTTP Server at %q", l)
	p := proxy{}
	wc := devutil.NewWasmCompiler().SetDir(".")
	mux := devutil.NewMux()
	mux.Match(Folder, p)
	mux.Match(devutil.NoFileExt, devutil.DefaultAutoReloadIndex.Replace(
		`<!-- styles -->`,
		`<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css">`))
	mux.Exact("/main.wasm", devutil.NewMainWasmHandler(wc))
	mux.Exact("/wasm_exec.js", devutil.NewWasmExecJSHandler(wc))
	mux.Default(devutil.NewFileServer().SetDir("."))

	log.Fatal(http.ListenAndServe(l, mux))
}
