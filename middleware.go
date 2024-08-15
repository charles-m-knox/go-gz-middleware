package gzm

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// These constants are copied from the gzip package, so that code that imports this
// module does not also have to import "compress/gzip".
const (
	NoCompression      = gzip.NoCompression
	BestSpeed          = gzip.BestSpeed
	BestCompression    = gzip.BestCompression
	DefaultCompression = gzip.DefaultCompression
	HuffmanOnly        = gzip.HuffmanOnly
)

// const (
// 	// Default minimum content length for enabling gzip compression. Not
// 	// currently implemented.
// 	DefaultMinContentLen = 860
// )

// gzipResponseWriter enables compatibility with the http.ResponseWriter struct
// so that it can be used as a middleware processor.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write is necessary in order to properly implement the io.Writer interface.
// It chooses to use gzip or not, based on the minimum content length threshold.
// Values less than the minLen threshold will not be compressed.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandler returns a handler that gzip-compresses all suitable requests. It
// can be chained with other HTTP handlers.
//
// It is recommended to use the exported compression levels included with this
// module, such as NoCompression, BestSpeed, BestCompression, etc.
//
// Currently, GzipHandler naively compresses anything it receives, regardless
// of the content length. This means that some payloads below a certain
// threshold will actually result in larger transfers.
func GzipHandler(next http.Handler, compressionLevel int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Header.Get("Accept-Encoding")
		if strings.Contains(c, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				log.Printf("failed to instantiate gzip middleware: %v", err.Error())
				next.ServeHTTP(w, r)
			}
			defer gz.Close()
			gzrw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(gzrw, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
