package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"io/ioutil"
)

const hdclubUrl = "http://hdclub.org/rss.php"

var imgRE *regexp.Regexp
var baseUrl = []byte("http://hdclub.org/")

func replaceBadHrefs(in []byte) []byte {
	out := imgRE.FindSubmatch(in)
	ok := (len(out[2]) >= 7 && string(out[2][:7]) == "http://") || (len(out[2]) >= 8 && string(out[2][:8]) == "https://")
	if ok {
		return in
	} else {
		var outBytes = make([]byte, 0, len(in)+len(baseUrl))
		outBytes = append(outBytes, out[1]...)
		outBytes = append(outBytes, baseUrl...)
		outBytes = append(outBytes, out[2]...)
		outBytes = append(outBytes, out[3]...)
		return outBytes
	}
}

func init() {
	imgRE = regexp.MustCompile(`(<img src=")(.*)(")`)
}

type HDClubHandler struct{}

func (HDClubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Only GET method allowed\n")
		return
	}

	var url string = fmt.Sprintf("%s?%s", hdclubUrl, r.URL.RawQuery)
	var resp *http.Response
	var err error
	resp, err = http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Bad response from hdclub: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Bad status code from hdclub: %d", resp.StatusCode)
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Can't read content from hdclub: %v", err)
		return
	}

	// copy headers from original server
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.Write(imgRE.ReplaceAllFunc(data, replaceBadHrefs))
}

func main() {

	http.Handle("/", HDClubHandler{})

	log.Fatal(http.ListenAndServe(":6060", nil))
}
