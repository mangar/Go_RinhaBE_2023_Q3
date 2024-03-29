package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rinha23/helpers"
)

// type PessoaBuscarResult struct {
// 	StatusCode	int
// 	Content		string
// }

type PessoaData struct {
	Id			string `json:"id"`
	Apelido 	string `json:"apelido"`
	Nome 		string `json:"nome"`
	Nascimento 	string `json:"nascimento"`
	Stack		[]string `json:"stack"`
}

type PessoaBuscar struct {
	w 			http.ResponseWriter
	r 			*http.Request
	t			string
	output		[]PessoaData
}

func NewPessoaBuscar(w http.ResponseWriter, r *http.Request) PessoaBuscar {
	w.Header().Set("Content-Type", "application/json")
	return PessoaBuscar{w:w, r:r, output: make([]PessoaData, 0)}
}

func (r *PessoaBuscar) Run() {

	if _, err := r.validateQueryParams(); err == nil {
		pessoasRedis, err := helpers.GetPessoaByTermo(r.t)
		helpers.LogOnError(err, "[PessoaBuscar.Run.01]")

		for _, pessoaRedis := range pessoasRedis {
			var pessoaData PessoaData
			json.Unmarshal([]byte(pessoaRedis), &pessoaData)
			r.output = append(r.output, pessoaData)
		}
		
	} else {
		r.w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(r.w, err.Error())
		return 
	}

	jsonData, _ := json.Marshal(r.output)
	fmt.Fprintf(r.w, string(jsonData))
	// return Result{w:r.w, StatusCode: http.StatusOK, Content: string(jsonData)}
}


func (r *PessoaBuscar) validateQueryParams() (*string, error) {

	queryParams := r.r.URL.Query()

	for key, values := range queryParams { 
		if key == "t" {
			for _, value := range values {
				r.t = value
				return &r.t, nil
			}
		}
	}
	return nil, errors.New("parametro 't' nao encontrado")
}
