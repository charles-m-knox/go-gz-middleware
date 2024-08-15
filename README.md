# go-gz-middleware

Offers a simple gzip middleware http handler that can be used like this:

```go
package main

import (
    "net/http"

    gzm "git.cmcode.dev/cmcode/go-gz-middleware"
)

// payload has a length of >860
var payload = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/short"
        w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello world"))
    case "/long":
        fallthrough
	default:
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload) // writes the longer payload
	}
}

func main() {
    http.Handle("/", gzm.GzipHandler(
        http.HandlerFunc(router),
        gzm.BestCompression,
        gzm.DefaultMinContentLen,
    ))

    http.ListenAndServe(":8080", nil)
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
