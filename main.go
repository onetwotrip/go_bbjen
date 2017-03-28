package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var jenkins_url, jenkins_scheme string

func send(url *url.URL, payload string) {
	q := url.Query()
	q.Set("payload", payload)
	url.RawQuery = q.Encode()
	url.Scheme = jenkins_scheme
	url.Host = jenkins_url

	log.Println(url.String())
	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		_, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parse(w http.ResponseWriter, req *http.Request) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)
	send(req.URL, bodyString)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func init() {
	jenkins_url = getenv("BBJEN_JENKINS_URL", "localhost:8080")
	jenkins_scheme = getenv("BBJEN_JENKINS_SCHEME", "http")
}

func main() {
	log.Println("Running on port :8082")

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{[a-z]+}/{[a-z]+}", parse).Methods("POST")
	http.Handle("/", rtr)

	log.Fatal(http.ListenAndServe(":8082", nil))
}
