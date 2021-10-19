package handler

import (
	"net/http"
	"html/template"
	"database/sql"
	"encoding/json"
	"strconv"
	"log"

	"github.com/NSMichelJ/go-todo-list/app/model"
)

var templates = template.Must(template.ParseGlob("public/templates/*"))

func MainHandler(w http.ResponseWriter, r *http.Request){
	err := templates.ExecuteTemplate(w, "indexPage", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func TaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.URL.Path != "/api/v1/task/" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string {
			"message": "not found",
		})
		return
	}

	if r.Method == "GET" {
		query := r.URL.Query()

		id, err := strconv.ParseInt(query.Get("id"), 10, 64)
		if err != nil {
			log.Print("Error: ", err)
			return
		}
		
		var task model.Task
		task.ID = id
		task, err = task.Get(db)
		if err != nil{
			log.Print("Error: ", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)

	} else if r.Method == "DELETE" {
		query := r.URL.Query()

		id, err := strconv.ParseInt(query.Get("id"), 10, 64)
		if err != nil {
			log.Print("Error: ", err)
			return
		}

		task := model.Task{}
		task.ID = id
		tasks, err := task.Delete(db)
		if err != nil{
			log.Print("Error: ", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	} else if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
		var task model.Task
    	decoder.Decode(&task)
		

		tasks, err := task.Update(db)
		if err != nil {
			log.Print("Error: ",err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tasks)

	}else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Method "+ r.Method + " no allowed"))	
	}

}

func TasksHandler(db *sql.DB, w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.URL.Path != "/api/v1/tasks/" {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string {
			"message": "not found",
		})
		return
	}

	if r.Method == "GET" {
		task := model.Task{}
		tasks, err := task.FethTasks(db)
		if err != nil {
			log.Print("Error: ", err)
        	return
    	}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var task model.Task
    	decoder.Decode(&task)

		id, err := task.Add(db)
		if err != nil {
			log.Print("Error: ", err)
			return
		}
		task.ID = id
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	} else if r.Method == "DELETE" {
		_, err := db.Query("DELETE FROM task")
		if err != nil {
			log.Print("Error: ", err)
		}
		tasks := []model.Task{}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Method "+ r.Method + " no allowed"))	
	}
	
}
