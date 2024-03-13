package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rinha23/helpers"

	"github.com/google/uuid"
)


type PessoasPostInput struct {
	Id			string `json:"id"`
	Apelido 	string `json:"apelido"`
	Nome 		interface{} `json:"nome"`
	Nascimento 	string `json:"nascimento"`
	Stack		[]interface{} `json:"stack"`
}

type PessoasPost struct {
	w http.ResponseWriter
	r *http.Request
	input *PessoasPostInput
}

func NewPessoasPost(w http.ResponseWriter, r *http.Request) PessoasPost {
	w.Header().Set("Content-Type", "application/json")
	return PessoasPost{w:w, r:r}
}

func (r *PessoasPost) Run() {

	input, err := r.validateContent()
	if err := helpers.LogOnError(err, "[handler.PessoasPost.ValidateContent]"); err != nil {
		http.Error(r.w, err.Error(), http.StatusUnprocessableEntity)
		return 
	}

	r.input = input

	err = r.validateSintax()
	if err := helpers.LogOnError(err, "[handler.PessoasPost.ValidateSintax]"); err != nil {
		http.Error(r.w, err.Error(), http.StatusBadRequest)
		return 
	}

	jsonData, _ := json.Marshal(r.input)

	if err = helpers.SetPessoa(input.Apelido, input.Id, string(jsonData)); err != nil {
		http.Error(r.w, err.Error(), http.StatusUnprocessableEntity)
		return 		
	}

	
	fmt.Fprintf(r.w, string(jsonData))
}


func (r *PessoasPost) validateContent() (*PessoasPostInput, error) {

	var input PessoasPostInput
	err := json.NewDecoder(r.r.Body).Decode(&input)
	if err = helpers.LogOnError(err, "[handler.PessoaPost.validate]"); err != nil {
		return nil, err
	} 
	
	// apelido nulo
	if input.Apelido == "" {
		return nil, errors.New("campo apelido null")
	}

	// nome nulo
	if input.Nome == "" {
		return nil, errors.New("campo nome null")
	}

	input.Id = uuid.New().String()

	return &input, nil
}


func (r *PessoasPost) validateSintax() error {

	// nome inteiro
	switch r.input.Nome.(type) {
	case float64:
		return errors.New("nome nao pode ser inteiro .. ")
	}

	// stack inteiro
	for i := range r.input.Stack {
		switch r.input.Stack[i].(type) {
		case float64:
			return errors.New("stack nao pode ser inteiro .. ")
		}
	}

	return nil
}
