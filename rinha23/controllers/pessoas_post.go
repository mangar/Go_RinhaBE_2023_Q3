package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rinha23/helpers"

	"github.com/sirupsen/logrus"
)


type PessoasPostInput struct {
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
	if err := LogOnError(err, "[handler.PessoasPost.ValidateContent]"); err != nil {
		http.Error(r.w, err.Error(), http.StatusUnprocessableEntity)
		return 
	}

	r.input = input

	err = r.validateSintax()
	if err := LogOnError(err, "[handler.PessoasPost.ValidateSintax]"); err != nil {
		http.Error(r.w, err.Error(), http.StatusBadRequest)
		return 
	}

	rdb := helpers.GetRedisConnection()
	rdb.Set(context.Background(), "aa", "aa")
	// rdb.Del(context.Background(), "cliente_extrato_" + strconv.Itoa(idCliente))
	// rdb.Del(context.Background(), "cliente_saldo_" + strconv.Itoa(idCliente))
	
	jsonData, _ := json.Marshal(r.input)
	fmt.Fprintf(r.w, string(jsonData))
}


func (r *PessoasPost) validateContent() (*PessoasPostInput, error) {

	var input PessoasPostInput
	err := json.NewDecoder(r.r.Body).Decode(&input)
	if err = LogOnError(err, "[handler.PessoaPost.validate]"); err != nil {
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





func LogOnError(err error, msg string) error {
	if err != nil {
		logrus.Error(msg + " .. " + err.Error())
		return errors.New(msg + " .. " + err.Error())
	} else {
		return nil
	}
}

