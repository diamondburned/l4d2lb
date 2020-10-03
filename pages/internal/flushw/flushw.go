package flushw

import "net/http"

// Wrap wraps the given ResponseWriter to always flush on write.
func Wrap(w http.ResponseWriter) http.ResponseWriter {
	flusher, _ := w.(http.Flusher)
	return flushWriter{w, flusher}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(Wrap(w), r)
	})
}

type flushWriter struct {
	http.ResponseWriter
	flusher http.Flusher
}

func (f flushWriter) Write(b []byte) (int, error) {
	n, err := f.ResponseWriter.Write(b)
	if err != nil {
		return n, err
	}

	if f.flusher != nil {
		f.flusher.Flush()
	}

	return n, nil
}
