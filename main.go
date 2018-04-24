package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	redirectCode = http.StatusFound
	redirects    = make(map[string]string)
	port         = "8080"
	bufferSize   = 100000
	logs         chan *http.Request
	disableCH    = false
)

var methods = [...]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func getKey(uri string, noSearch bool) string {
	path := strings.Split(uri, "*")
	params := strings.Split(path[0], "?")

	if _, ok := redirects[params[0]]; noSearch || ok {
		return params[0]
	} else {
		//Find most nearest key by largest key length
		var nearestKeys string
		var maxKeyLen int
		for key := range redirects {
			if strings.Contains(params[0], key) && len(key) > maxKeyLen {
				maxKeyLen = len(key)
				nearestKeys = key
			}
		}
		return nearestKeys
	}
}

func generateAnswer(r *http.Request) string {
	target := redirects[getKey(r.RequestURI, false)]
	target = strings.Replace(target, "{URI}", r.RequestURI, -1)
	return strings.Replace(target, "{HOST}", r.Host, -1)
}

func redirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, generateAnswer(r), redirectCode)
	if !disableCH {
		logs <- r
	}
}

func load(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "%d", len(logs))
}

func main() {
	if portEnv, present := os.LookupEnv("PORT"); present {
		if v, err := strconv.Atoi(portEnv); err == nil {
			if v > 0 && v < 65536 {
				port = portEnv
			}
		}
	}

	if bufferSizeEnv, present := os.LookupEnv("BUFFER"); present {
		if v, err := strconv.Atoi(bufferSizeEnv); err == nil {
			if v > 0 {
				bufferSize = v
			}
		}
	}

	if v, present := os.LookupEnv("DISABLE_CH"); present && v == "true" {
		disableCH = true
	}

	logs = make(chan *http.Request, bufferSize)

	if !disableCH {
		go logger(logs)
	}

	router := httprouter.New()
	router.GET("/load", load)

	if v, present := os.LookupEnv("PERMANENTLY"); present && v == "true" {
		redirectCode = http.StatusMovedPermanently
	}

	if redirectsEnv, present := os.LookupEnv("REDIRECTS"); present {
		entries := strings.Split(redirectsEnv, "|")
		for _, entry := range entries {
			kv := strings.Split(entry, " ")
			if len(kv) == 2 {
				redirects[getKey(kv[0], true)] = kv[1]
				for _, method := range methods {
					router.Handle(method, kv[0], redirect)
					log.Println("Registred " + method + " " + kv[0])
				}
			} else {
				log.Fatal("Failded to parse: " + entry)
			}
		}
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
