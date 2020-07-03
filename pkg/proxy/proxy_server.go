package proxy

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bmaynard/apimock/pkg/filesystem"
	l "github.com/bmaynard/apimock/pkg/utils/logger"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	l.Log.Warn("Capturing TLS requests currently not supported")
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request, isK8s bool) {
	originalHost := req.URL.Hostname()
	host := strings.Split(req.Host, ":")

	if originalHost == "" {
		originalHost = host[0]
	}

	if isK8s {
		port := "80"

		if len(host) >= 2 {
			port = host[1]
		}

		req.URL.Host = "localhost:" + port
		req.Host = "localhost:" + port

		if os.Getenv("SERVICE_HOST_NAME") != "" {
			originalHost = os.Getenv("SERVICE_HOST_NAME")
		}
	}

	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}

	l.Log.Infof("Processing HTTP Request for: %s", originalHost)
	resp, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		// We need to read the body to write the mock file, then re-set resp.body for the response
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		err := filesystem.GetAdapter().WriteMockFile(resp, bodyBytes, originalHost)

		if err != nil {
			l.Log.Fatal(err)
		}
	}

	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func Listen(pemPath string, keyPath string, addr string, isK8s bool) {
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r, isK8s)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	l.Log.Info("Starting proxy server")

	if len(pemPath) > 1 && len(keyPath) > 1 {
		l.Log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
	} else {
		l.Log.Fatal(server.ListenAndServe())
	}
}
