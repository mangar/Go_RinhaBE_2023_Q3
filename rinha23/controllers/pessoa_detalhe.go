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

	jsonDataPessoa, err := helpers.GetPessoaById(r.IdPessoa)
	if err := helpers.LogOnError(err, "[handler.PessoaDetalhe.GetPessoaById]"); err != nil {
		http.Error(r.w, err.Error(), http.StatusNotFound)
		return 
	}

	fmt.Fprintf(r.w, jsonDataPessoa)
}

