package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/test-get", testGet)
	http.HandleFunc("/test-file", testReceiveFile)
	http.HandleFunc("/read-header", readHeader)
	http.HandleFunc("/", generic)

	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.DisableKeepAlives = true
		tr.MaxIdleConnsPerHost = 1
		tr.CloseIdleConnections()
	}

	//run a client-test
	server := &http.Server{Addr: ":8081", ReadTimeout: 2 * time.Second, WriteTimeout: 2 * time.Second}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed start client test", err)
		os.Exit(1)
	}
}

func testGet(w http.ResponseWriter, r *http.Request) {
	generic(w, r)
}

func testPost(w http.ResponseWriter, r *http.Request) {
	generic(w, r)
}

func testReceiveFile(w http.ResponseWriter, r *http.Request) {
	//check if post 1 file
	_, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%q", dump)
	content := "====================== DumpRequest ================================= <br>"
	content += string(dump)
	content += "==================================================================== <br>"
	content += fmt.Sprintf("Received file [%s] size [%s]", fh.Filename, fh.Size)
	fmt.Fprintf(w, content)
}

func generic(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%q", dump)
}
func readHeader(w http.ResponseWriter, r *http.Request) {

	content := "Header <br />"
	for key, value := range r.Header {
		content += fmt.Sprintf("Key [%s]: %s <br />", key, strings.Join(value, ","))
	}
	content += "<br /> Trailer"
	for key, value := range r.Trailer {
		content += fmt.Sprintf("Key [%s]: %s <br />", key, strings.Join(value, ","))
	}
	content += "<br /> Another info"
	content += fmt.Sprintf("Key [Method]: %v <br />", r.Method)
	content += fmt.Sprintf("Key [URL INFO]: (%v)(%v)(%v)(%v)(%v)(%v)(%v)(%v) <br />", r.URL.Scheme, r.URL.Opaque, r.URL.Host, r.URL.Path, r.URL.RawPath, r.URL.ForceQuery, r.URL.RawQuery, r.URL.Fragment)
	content += fmt.Sprintf("Key [Proto]: %v <br />", r.Proto)
	content += fmt.Sprintf("Key [TransferEncoding]: %v <br />", r.TransferEncoding)
	content += fmt.Sprintf("Key [Host]: %v <br />", r.Host)
	content += fmt.Sprintf("Key [RequestURI]: %v <br />", r.RequestURI)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, content)
}
