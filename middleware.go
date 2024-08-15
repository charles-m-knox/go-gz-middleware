package gzm

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	recorder "git.cmcode.dev/cmcode/go-http-response-recorder"
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

const (
	// Default minimum content length for enabling gzip compression. Not
	// currently implemented.
	DefaultMinContentLen = 860
)

// gzipResponseWriter enables compatibility with the http.ResponseWriter struct
// so that it can be used as a middleware processor.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write is necessary in order to properly implement the io.Writer interface.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	return n, err
}

// GzipHandler returns a handler that gzip-compresses all suitable requests. It
// can be chained with other HTTP handlers.
//
// It is recommended to use the exported compression levels included with this
// module, such as NoCompression, BestSpeed, BestCompression, etc.
//
// It chooses to use gzip or not, based on the minimum content length threshold.
// Values less than the minLen threshold will not be compressed.
func GzipHandler(next http.Handler, compressionLevel int, minLen int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Header.Get("Accept-Encoding")
		if !strings.Contains(c, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		b := &bytes.Buffer{}

		rec := recorder.NewResponseRecorder(w, b, func(status int, header http.Header) bool { return true })
		next.ServeHTTP(rec, r)

		s := rec.Size()

		if s < minLen {
			// No need to gzip, just return the original response
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			log.Printf("failed to instantiate gzip middleware: %v", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		defer gz.Close()

		gzrw := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzrw, r)
	})
}
