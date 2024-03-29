package controllers

import (
	"context"
	"net/http"
	"rinha23/helpers"
	"strconv"
)

// type ContagemPessoasResult struct {
// 	Content		string
// 	StatusCode	int
// }

type ContagemPessoas struct {
	w 			http.ResponseWriter
	r 			*http.Request
}

func NewContagemPessoas(w http.ResponseWriter, r *http.Request) ContagemPessoas {
	w.Header().Set("Content-Type", "application/json")
	return ContagemPessoas{w:w, r:r}
}

func (r *ContagemPessoas) Run() Result {

	count := 0
	ctx := context.Background()

	rows, err := helpers.GetDBConnection().Query(ctx, "select count(*) from rinha23_clientes")
	defer rows.Close()

	helpers.LogOnError(err, "[NewContagemPessoas]")

	for rows.Next() {
		err = rows.Scan(&count)
		helpers.LogOnError(err, "[NewContagemPessoas.Next]")
	}

	return Result{w:r.w, StatusCode: 200, Content: strconv.Itoa(count)}
}