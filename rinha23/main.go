package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rinha23/controllers"
	"rinha23/helpers"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var contagemPessoasResultChan chan controllers.ContagemPessoasResult
var pessoaDetalheResultChan chan controllers.PessoaDetalheResult
var pessoaBuscarResultChan chan controllers.PessoaBuscarResult
var pessoaSalvarResultChan chan controllers.PessoaSalvarResult

// worker é a função que cada worker vai executar.
// Recebe tarefas de um channel, processa-as e notifica a finalização no WaitGroup.
func worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done() // Notifica a conclusão deste worker ao WaitGroup.
	for task := range tasks {
		logrus.Debug(fmt.Sprintf(">>>>>> Worker %v, Controller:%v, Id:%v", id, task.Controller, task.ID))

		switch task.Controller {
		case "hello":
			fmt.Fprintf(task.w, "[" + task.Controller + "]Hello 👋!")

		case "pessoaSalvar":
			a := controllers.NewPessoasPost(task.w, task.r)
			pessoaSalvarResultChan <- a.Run()

		case "pessoaBuscar":
			a := controllers.NewPessoaBuscar(task.w, task.r)
			pessoaBuscarResultChan <- a.Run()			

		case "pessoaDetalhe":
			vars := mux.Vars(task.r)
			c := controllers.NewPessoaDetalhe(task.w, task.r, vars["ID"])
			pessoaDetalheResultChan <- c.Run()

		case "contagem-pessoas":
			c := controllers.NewContagemPessoas(task.w, task.r)
			contagemPessoasResultChan <- c.Run()

		default:
			task.w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(task.w, "[" + task.Controller + "] Controller not found")
		}

		logrus.Debug(fmt.Sprintf("<<<<<< Worker %v, Controller:%v, Id:%v", id, task.Controller, task.ID))
	}
}


type Task struct {
	ID			int
	Controller 	string
	w 			http.ResponseWriter
	r 			*http.Request
}

func CreateTask(controller string, w http.ResponseWriter, r *http.Request) Task {
	return Task{ ID: time.Now().Nanosecond(), w:w, r:r, Controller: controller }
}


func main() {
	godotenv.Load(".env")
	helpers.SetupLog()	
	logrus.Info(">>>>>>>>>>   " + os.Getenv("SERVER_NAME")  + "   <<<<<<<<<< ")

	// DB
	helpers.TestDBConnection()

	// Redis
	helpers.TestRedisConnection()


	// logrus.Debug("Starting server at port ", os.Getenv("WEB_PORT"))
    // if err := http.ListenAndServe(":" + os.Getenv("WEB_PORT"), controllers.SetupRoutes()); err != nil {
    //     log.Fatal(err)
    // }

	logrus.Debug(">> 0")

	const numWorkers = 1
	taskChannel := make(chan Task, 50)
	contagemPessoasResultChan =	make(chan controllers.ContagemPessoasResult)
	pessoaDetalheResultChan = make(chan controllers.PessoaDetalheResult)
	pessoaBuscarResultChan = make(chan controllers.PessoaBuscarResult)
	pessoaSalvarResultChan = make(chan controllers.PessoaSalvarResult)

	logrus.Debug(">> 1")

	var wg sync.WaitGroup

	// Initialize the worker pool.
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChannel, &wg)
	}


	logrus.Debug(">> 2")

	routers := setupRoutes(taskChannel)

	logrus.Debug(">> 3")


	go func() {
		logrus.Debug("Starting server at port ", os.Getenv("WEB_PORT"))
		if err := http.ListenAndServe(":" + os.Getenv("WEB_PORT"), routers); err != nil {
			log.Fatal(err)
		}
	}()

	

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(taskChannel)
	wg.Wait() // Espera todos os workers terminarem.

}



func setupRoutes(taskChannel chan Task) *mux.Router {
	logrus.Debug("[Routes] Seting up routes.")

	router := mux.NewRouter()

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// if !isTestRequest(w, r) {
			taskChannel <- CreateTask("hello", w, r)
		// }
	})

	router.HandleFunc("/pessoas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			taskChannel <- CreateTask("pessoaSalvar", w, r)
			result := <- pessoaSalvarResultChan

			// r.w.Header().Set("Location", "/pessoas/" + r.input.Id)
			// r.w.WriteHeader(http.StatusCreated)
			// fmt.Fprintf(r.w, string(jsonData))

			// if result.StatusCode == 200 {
			// 	w.Header().Set("Location", "/pessoas/" + r.input.Id)
			// }

			for key, value  := range result.Headers {
				w.Header().Set(key, value)
			}
			w.WriteHeader(result.StatusCode)
			fmt.Fprintf(w, fmt.Sprintf("%v", result.Content))

		} else if r.Method == http.MethodGet {
			taskChannel <- CreateTask("pessoaBuscar", w, r)
			result := <- pessoaBuscarResultChan

			w.WriteHeader(result.StatusCode)
			fmt.Fprintf(w, fmt.Sprintf("%v", result.Content))
		}
	})

	router.HandleFunc("/pessoas/{ID}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskChannel <- CreateTask("pessoaDetalhe", w, r)
			result := <- pessoaDetalheResultChan

			w.WriteHeader(result.StatusCode)
			fmt.Fprintf(w, fmt.Sprintf("%v", result.Content))
		}
	})

	router.HandleFunc("/contagem-pessoas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskChannel <- CreateTask("contagem-pessoas", w, r)
			result := <- contagemPessoasResultChan

			w.WriteHeader(result.StatusCode)
			fmt.Fprintf(w, fmt.Sprintf("%v", result.Count))
		}
	})

	return router
}
