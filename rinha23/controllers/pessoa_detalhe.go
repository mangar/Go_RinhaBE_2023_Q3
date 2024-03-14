package controllers

import (
	"net/http"
	"rinha23/helpers"
)

type PessoaDetalheResult struct {
	StatusCode	int
	Content		string
}

type PessoaDetalhe struct {
	w 			http.ResponseWriter
	r 			*http.Request
	IdPessoa	string
}

func NewPessoaDetalhe(w http.ResponseWriter, r *http.Request, id string) PessoaDetalhe {
	w.Header().Set("Content-Type", "application/json")
	return PessoaDetalhe{w:w, r:r, IdPessoa: id}
}

func (r *PessoaDetalhe) Run() PessoaDetalheResult {
	jsonDataPessoa, err := helpers.GetPessoaById(r.IdPessoa)
	if err := helpers.LogOnError(err, "[handler.PessoaDetalhe.GetPessoaById]"); err != nil {
		return PessoaDetalheResult{StatusCode: http.StatusNotFound, Content:err.Error()}
	}
	return PessoaDetalheResult{StatusCode: 200, Content:jsonDataPessoa}
}

