package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

func main() {
	var err error
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Println("Error listening:", err.Error())
	}
	log.Println("Listening on", l.Addr())
	hdlr := RegexHandler{}
	err = http.Serve(l, &hdlr)
	if err != nil {
		log.Println("Error running server:", err.Error())
	}
}

type RegexHandler struct {
}

func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		index, _ := os.ReadFile("./index.html")
		w.Write(index)
		return
	}
	if r.URL.Path == "/api" {
		var matches []string
		errstr := ""
		rbytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Failed to read body:", err.Error())
			errstr = "Failed to read request body."
		} else {
			rdata := struct {
				Regex string
				Text  string
			}{}
			err = json.Unmarshal(rbytes, &rdata)
			if err != nil {
				log.Println("Failed to parse form data:", err.Error())
				errstr = "Failed to parse form data."
			} else {
				re, err := regexp.Compile(rdata.Regex)
				if err != nil {
					errstr = err.Error()
				} else {
					matches = re.FindStringSubmatch(rdata.Text)
				}
			}
		}
		data := struct {
			Error   string
			Matches []string
		}{errstr, matches}
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data:", err.Error())
			bytes = []byte(`{"Error":"failed to marshal data","Matches":[]}`)
		}
		w.Write(bytes)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
