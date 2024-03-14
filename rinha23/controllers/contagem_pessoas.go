package controllers

import (
	"context"
	"fmt"
	"net/http"
	"rinha23/helpers"
)

type ContagemPessoas struct {
	w 			http.ResponseWriter
	r 			*http.Request
}

func NewContagemPessoas(w http.ResponseWriter, r *http.Request) ContagemPessoas {
	w.Header().Set("Content-Type", "application/json")
	return ContagemPessoas{w:w, r:r}
}

func (r *ContagemPessoas) Run() {

	count := 0
	ctx := context.Background()

	rows, err := helpers.GetDBConnection().Query(ctx, "select count(*) from rinha23_clientes")
	defer rows.Close()

	helpers.LogOnError(err, "[NewContagemPessoas]")

	for rows.Next() {
		err = rows.Scan(&count)
		helpers.LogOnError(err, "[NewContagemPessoas.Next]")
	}

	fmt.Fprintf(r.w, fmt.Sprintf("%v", count))
}