# go-gz-middleware

Offers a simple, minimal gzip middleware http handler that can be used like this:

```go
package main

import (
	"log"
	"net/http"

	gzm "git.cmcode.dev/cmcode/go-gz-middleware"
)

var (
	// longPayload has a length of >860
	longPayload = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	// shortPayload has a length of <860
	shortPayload = []byte("Hello world")
)

func router(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	switch r.URL.Path {
	case "/short":
		_, _ = w.Write(shortPayload)
	case "/long":
		fallthrough
	default:
		_, _ = w.Write(longPayload) // writes the longer payload
	}
}

func main() {
	http.Handle("/", gzm.GzipHandler(
		http.HandlerFunc(router),
		gzm.BestCompression,
		gzm.DefaultMinContentLen,
	))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("failed to listen and serve: %v", err.Error())
	}
}
```

Available options for the handler are

- **Compression level**
  - `gzm.NoCompression`
  - `gzm.BestSpeed`
  - `gzm.BestCompression`
  - `gzm.DefaultCompression`
  - `gzm.HuffmanOnly`
- **Minimum content length** - gzip compression is counterproductive if the payload is too small.
  - `gzm.DefaultMinContentLen`

Under the hood, this project leverages [go-http-response-recorder](https://git.cmcode.dev/cmcode/go-http-response-recorder) for determining content length before writing responses to the next middleware handler.
