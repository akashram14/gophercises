package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/debug", sourceCodeHandler)
	log.Fatal(http.ListenAndServe(":8080", recoverMw(mux, true)))

}

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	b := bytes.NewBuffer(nil)
	io.Copy(b, file)
	_ = quick.Highlight(w, b.String(), "go", "html", "monokai")
}

func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				if !dev {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, string(stack))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

//func (rw *responseWriter) Write(b []byte) (int, error) {
//	rw.writes = append(rw.writes, b)
//	return len(b), nil
//}
//
//func (rw *responseWriter) WriteHeader(statuscode int) {
//	rw.status = statuscode
//}
//
//func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
//	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
//	if !ok {
//		return nil, nil, fmt.Errorf("the responseWriter does not implement the Hijacker interface")
//	}
//	return hijacker.Hijack()
//}
//
//func (rw *responseWriter) Flush() {
//	flusher, ok := rw.ResponseWriter.(http.Flusher)
//	if !ok {
//		return
//	}
//	flusher.Flush()
//}

//func (rw *responseWriter) flush() error {
//	if rw.status != 0 {
//		rw.ResponseWriter.WriteHeader(rw.status)
//	}
//	for _, write := range rw.writes {
//		_, err := rw.ResponseWriter.Write(write)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>helolo!</h1")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1")
}
