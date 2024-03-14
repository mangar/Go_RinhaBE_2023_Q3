package controllers

// // worker é a função que cada worker vai executar.
// // Recebe tarefas de um channel, processa-as e notifica a finalização no WaitGroup.
// func Worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
// 	defer wg.Done() // Notifica a conclusão deste worker ao WaitGroup.
// 	for task := range tasks {
// 		logrus.Debug("> > > > > > Worker %d iniciou tarefa %v\n", id, task)

// 		logrus.Debug(">> CONTROLLER: ", task.Controller)
// 		fmt.Fprintf(task.w, "[" + task.Controller + "]Hello 👋!")

// 		logrus.Debug("< < < < < < Worker %d completou tarefa %v\n", id, task)
// 	}
// }

// type Task struct {
// 	ID			int
// 	Controller 	string
// 	w 			http.ResponseWriter
// 	r 			*http.Request
// }

// func SetupRoutes(taskChannel chan Task) *mux.Router {
// 	logrus.Debug("[Routes] Seting up routes.")

// 	router := mux.NewRouter()

// 	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		if !isTestRequest(w, r) {

// 			task := Task{
// 				ID:      time.Now().Nanosecond(),
// 				Controller: "extrato",
// 				w:w,
// 				r:r,
// 			}
// 			taskChannel <- task

// 		}
// 	})

// 	router.HandleFunc("/pessoas", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodPost {
// 			a := NewPessoasPost(w,r)
// 			a.Run()

// 		} else if r.Method == http.MethodGet {
// 			a := NewPessoaBuscar(w,r)
// 			a.Run()
// 		}
// 	})

// 	router.HandleFunc("/pessoas/{ID}", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodGet {
// 			vars := mux.Vars(r)
// 			c := NewPessoaDetalhe(w,r,vars["ID"])
// 			c.Run()
// 		}
// 	})

// 	router.HandleFunc("/contagem-pessoas", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodGet {
// 			c := NewContagemPessoas(w,r)
// 			c.Run()
// 		}
// 	})

// 	return router
// }

// func isTestRequest(w http.ResponseWriter, r *http.Request) bool {
// 	isTest := false
// 	content := make(map[string]string)

// 	testValue := r.Header.Get("X-Test")
// 	if testValue == "true" {

// 		w.Header().Set("Content-Type", "application/json")

// 		// Headers
// 		for name, values := range r.Header {
// 			for _, value := range values {
// 				content["HEADER:" + name] = value
// 			}
// 		}

// 		// PATH PARAMS
// 		for k, v := range mux.Vars(r) {
// 			content["PATH:" + k] = v
// 		}

// 		// QUERY
// 		for k, v := range r.URL.Query() {
// 			content["PATH:" + k] = v[0]
// 		}

// 		content["URL:Path"] = r.URL.Path
// 		content["URL:RawQuery"] = r.URL.RawQuery
// 		content["URL:FullPath"] = r.URL.Path + r.URL.RawQuery

// 		jsonData, _ := json.Marshal(content)

// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprintf(w, string(jsonData))
// 		isTest = true
// 	}
// 	return isTest
// }


