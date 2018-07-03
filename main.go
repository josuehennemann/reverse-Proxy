package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/josuehennemann/logger"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	flag.StringVar(&pathConfig, "config", "", "path config")
	flag.Parse()
	if isEmpty(pathConfig) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	NewConfig()
	var err error

	Logger, err = logger.New(config.GetLogFile()+"reverse-proxy.log", logger.LEVEL_ALL, true)
	CheckErrorAndKillMe(err)

	rules, err = initRuleList()
	CheckErrorAndKillMe(err)
	//check listen https
	if config.GetHttpsListen() != "" {
		go ServerHttps()
	}
	ServerHttp()
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func ServerHttp() {
	// http handlers

	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/admin/reload-rules", reloadRules)
	http.HandleFunc("/", capture)

	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.DisableKeepAlives = true
		tr.MaxIdleConnsPerHost = 1
		tr.CloseIdleConnections()
	}
	Logger.Printf(logger.INFO, "Initiating the service  [%s] ...", config.GetHttpListen())

	//listen http service
	server := &http.Server{Addr: config.GetHttpListen(), ReadTimeout: HTTP_READ_TIMEOUT, WriteTimeout: HTTP_WRITE_TIMEOUT}
	err := server.ListenAndServe()
	if err != nil {
		Logger.Fatalf("Can't start service http: %s\n", err.Error())
	}
}

func ServerHttps() {

	configTLS := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{ // define TLS
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		},
	}

	configTLS.RootCAs = x509.NewCertPool()
	//read a CA file and append in PEM
	if ca, err := ioutil.ReadFile(config.GetHttpsCertificate() + "server.ca"); err == nil {
		configTLS.RootCAs.AppendCertsFromPEM(ca)
	}

	server := &http.Server{Addr: config.GetHttpsListen(), TLSConfig: configTLS, ReadTimeout: HTTP_READ_TIMEOUT, WriteTimeout: HTTP_WRITE_TIMEOUT}
	err := server.ListenAndServeTLS(config.GetHttpsCertificate()+"server.crt", config.GetHttpsCertificate()+"server.key")
	if err != nil {
		Logger.Fatalf("Can't start service http: %s\n", err.Error())
	}
}
func capture(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	var finalDestination string
	// write log in any return
	defer func(destiny *string) {
		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		Logger.Printf(logger.INFO, "Redirecting from [%s] to [%s]. Ip Request: [%s]", r.URL.Path, *destiny, clientIP)
	}(&finalDestination)

	referer := r.Header.Get("Referer")

	tmp := strings.Split(r.URL.Path, "/")
	uri := tmp[1]
	rule := rules.GetRule(uri)
	if rule == nil && referer == "" {
		http.NotFound(w, r)
		return
	}
	// check if referer in request
	if rule == nil && referer != "" {
		urlReferer, err := url.Parse(referer)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		tmp := strings.Split(urlReferer.Path, "/")
		uri = tmp[1]
		rule = rules.GetRule(uri)
		if rule == nil {

			http.NotFound(w, r)
			return
		}
	}

	//rediret not repass headers
	if rule.Redirect {
		params := ""
		if tmp := r.URL.Query(); len(tmp) > 0 {
			params = "?" + tmp.Encode()
		}
		finalDestination = rule.Destiny + params
		http.Redirect(w, r, finalDestination, http.StatusFound)
		return
	}

	u, err := url.Parse(rule.Destiny)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	// "build" a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(u)
	if rule.RemoveOrigin {
		//build a new path
		r.URL.Path = strings.Replace(r.URL.Path, "/"+rule.Origin+"/", "/", -1)
	}
	//catch a orginal director func
	nickFury := proxy.Director

	//"override" a director function
	proxy.Director = func(req *http.Request) {
		//exec "original" director
		nickFury(req)

		req.Header.Set("Access-Control-Allow-Origin", "*")
		req.Header.Set("Access-Control-Allow-Headers", "X-Requested-With")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:60.0) Gecko/20100101 Firefox/60.0")

		req.Host = u.Host // alter host to host of destiny

		if rule.RemoveOrigin {
			r.URL.Path = strings.Replace(r.URL.Path, "/"+rule.Origin+"/", "/", -1)
		}
	}

	//TODO: implement function, if necessary
	/*proxy.ModifyResponse = func(resp *http.Response) error {
		return nil
	}*/

	// set var to write log
	finalDestination = u.Host + r.URL.Path
	proxy.ServeHTTP(w, r)
}

func reloadRules(w http.ResponseWriter, r *http.Request) {
	if err := rules.Reload(); err != nil {
		Logger.Fatalf("Failed reload rules [%s]", err.Error())
	}
	text := ""
	for k, v := range rules.List() {
		text += fmt.Sprintf("%s = %v <br>", k, *v)
	}
	fmt.Fprintln(w, "Reload success <br> "+text)
}
