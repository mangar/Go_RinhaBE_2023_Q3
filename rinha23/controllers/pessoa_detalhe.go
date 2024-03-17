package controllers

import (
	"fmt"
	"net/http"
	"rinha23/helpers"
)

type PessoaDetalhe struct {
	w 			http.ResponseWriter
	r 			*http.Request
	IdPessoa	string
}

func NewPessoaDetalhe(w http.ResponseWriter, r *http.Request, id string) PessoaDetalhe {
	w.Header().Set("Content-Type", "application/json")
	return PessoaDetalhe{w:w, r:r, IdPessoa: id}
}

func (r *PessoaDetalhe) Run() {
	jsonData, err := helpers.GetPessoaById(r.IdPessoa)
	if err := helpers.LogOnError(err, "[handler.PessoaDetalhe.GetPessoaById]"); err != nil {

		r.w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(r.w, err.Error())
		return 
		// return Result{w:r.w, StatusCode: http.StatusNotFound, Content:err.Error()}
	}

	fmt.Fprintf(r.w, string(jsonData))
	// return Result{w:r.w, StatusCode: 200, Content:jsonDataPessoa}
}

