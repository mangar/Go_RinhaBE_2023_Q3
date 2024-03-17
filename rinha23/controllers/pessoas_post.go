package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rinha23/helpers"
	"strings"

	"github.com/google/uuid"
)

// type PessoaSalvarResult struct {
// 	StatusCode	int
// 	Content		string
// 	Headers map[string]string
// }

type PessoasPostInput struct {
	Id			string `json:"id"`
	Apelido 	string `json:"apelido"`
	Nome 		interface{} `json:"nome"`
	Nascimento 	string `json:"nascimento"`
	Stack		[]interface{} `json:"stack"`
	SearchContent		string `json:"searchContent"`
}

func (i *PessoasPostInput) generateSearhContent() string {
	
	var stacks strings.Builder

	for _, stack := range i.Stack {
		ss, _ := stack.(string)
		stacks.WriteString(ss + "|")
	}

	nome, _ := i.Nome.(string)

	i.SearchContent = fmt.Sprintf("%v|%v|%v", strings.ToLower(i.Apelido), strings.ToLower(nome), strings.ToLower(stacks.String()))
	return i.SearchContent
}

func (i *PessoasPostInput) GetStack() string {	
	var stacks strings.Builder
	for _, stack := range i.Stack {
		ss, _ := stack.(string)
		stacks.WriteString(ss + ",")
	}
	return stacks.String()
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

func (r *PessoasPost) Run() Result {

	input, err := r.validateContent()
	if err := helpers.LogOnError(err, "[handler.PessoasPost.ValidateContent]"); err != nil {
		// http.Error(r.w, err.Error(), http.StatusUnprocessableEntity)
		return Result{w:r.w, StatusCode: http.StatusUnprocessableEntity}
	}

	r.input = input

	err = r.validateSintax()
	if err := helpers.LogOnError(err, "[handler.PessoasPost.ValidateSintax]"); err != nil {
		http.Error(r.w, err.Error(), http.StatusBadRequest)
		return Result{w:r.w, StatusCode: http.StatusBadRequest}
	}

	r.input.generateSearhContent()
	jsonData, _ := json.Marshal(r.input)

	if err = helpers.SetPessoa(input.Apelido, input.Id, input.SearchContent, string(jsonData)); err != nil {
		http.Error(r.w, err.Error(), http.StatusUnprocessableEntity)
		return Result{w:r.w, StatusCode: http.StatusUnprocessableEntity}
	}
	// 
	// DB. Salvar no banco de dados.....
	r.insertPessoa()
	
	r.w.Header().Set("Location", "/pessoas/" + r.input.Id)
	return Result{ w:r.w,
		StatusCode: http.StatusCreated,
		Content: string(jsonData),
		Headers: map[string]string{ "Location":"/pessoas/" + r.input.Id },
	}

}


func (r *PessoasPost) insertPessoa() error {
	ctx := context.Background()

	_, err := helpers.GetDBConnection().Exec(ctx, `
	INSERT INTO rinha23_clientes (id, apelido, nome, nascimento, stack, search_content) VALUES ($1, $2, $3, $4, $5, $6);`, 
	r.input.Id, r.input.Apelido, r.input.Nome, r.input.Nascimento, r.input.GetStack(), r.input.SearchContent)

	if err = helpers.LogOnError(err, "[NewPessoasPost.Run.insertPessoa]"); err != nil {
		return err
	}

	return nil
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
