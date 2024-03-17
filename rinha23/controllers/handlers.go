package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func SetupRoutes() *mux.Router {
	logrus.Debug("[Routes] Seting up routes.")

	router := mux.NewRouter()

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if !isTestRequest(w, r) {
			fmt.Fprintf(w, "Hello 👋!")
		}
	})

	router.HandleFunc("/pessoas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			a := NewPessoasPost(w,r)
			result := a.Run()
			writeOut(result, w)

		} else if r.Method == http.MethodGet {
			a := NewPessoaBuscar(w,r)
			a.Run()
			// result := a.Run()
			// writeOut(result, w)
		}
	})

	router.HandleFunc("/pessoas/{ID}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			vars := mux.Vars(r)
			c := NewPessoaDetalhe(w,r,vars["ID"])
			c.Run()
			// result := c.Run()
			// writeOut(result, w)
		}
	})

	router.HandleFunc("/contagem-pessoas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			c := NewContagemPessoas(w,r)
			result := c.Run()
			writeOut(result, w)
		}
	})

	logrus.Debug("[Routes] DONE Seting up routes.")
	return router
}


type Result struct {
	w 			http.ResponseWriter
	StatusCode	int
	Content		string
	Headers map[string]string
}

func writeOut(result Result, w http.ResponseWriter) {

	if result.StatusCode != http.StatusOK {
		result.w.WriteHeader(result.StatusCode)
	}

	// logrus.Debug("HEADER:", result.Headers)

	// for key, value := range result.Headers {
	// 	logrus.Debug(" - ", key, ":", value)
	// 	result.w.Header().Set(key, value)

	// 	result.w.Header().Set("w1", "w1")

	// 	w.Header().Set(key, value)

	// 	w.Header().Set("w2", "w2")
	// }

	fmt.Fprintf(result.w, result.Content)
}


func isTestRequest(w http.ResponseWriter, r *http.Request) bool {
	isTest := false
	content := make(map[string]string)

	testValue := r.Header.Get("X-Test")
	if testValue == "true" {

		w.Header().Set("Content-Type", "application/json")
		
		// Headers
		for name, values := range r.Header {
			for _, value := range values {
				content["HEADER:" + name] = value
			}
		}

		// PATH PARAMS
		for k, v := range mux.Vars(r) {
			content["PATH:" + k] = v
		}

		// QUERY
		for k, v := range r.URL.Query() {
			content["PATH:" + k] = v[0]
		}


		content["URL:Path"] = r.URL.Path
		content["URL:RawQuery"] = r.URL.RawQuery
		content["URL:FullPath"] = r.URL.Path + r.URL.RawQuery

		jsonData, _ := json.Marshal(content)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(jsonData))
		isTest = true
	}
	return isTest
}


