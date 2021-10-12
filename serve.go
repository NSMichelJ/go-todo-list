package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct{
	ID int64 `json:"id"`
	Content string `json:"content"`
	Created string `json:"created"`
}

var templates = template.Must(template.ParseGlob("public/templates/*"))

var db *sql.DB

func main(){
	var err error
	db, err = getDB()
	if err != nil{
		log.Print("error", err)
		return
	}

	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/api/v1/task/", taskHandler)
	http.HandleFunc("/api/v1/tasks/",tasksHandler)

	log.Print("Server runnig on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request){
	err := templates.ExecuteTemplate(w, "indexPage", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func taskHandler(w http.ResponseWriter, r *http.Request){
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

		task, err := get(id)
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

		tasks, err := delete(id)
		if err != nil{
			log.Print("Error: ", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	} else if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
		var params Task
    	decoder.Decode(&params)
		
		tasks, err := update(params.Content, params.ID)
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

func tasksHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.URL.Path != "/api/v1/tasks/" {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string {
			"message": "not found",
		})
		return
	}

	if r.Method == "GET" {
		tasks, err := fethTasks()
		if err != nil {
			log.Print("Error: ", err)
        	return
    	}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var params Task
    	decoder.Decode(&params)

		id, err := add(params.Content, params.Created)
		if err != nil {
			log.Print("Error: ", err)
			return
		}

		task, err := get(id)
		if err != nil {
			log.Print("Error: ", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	} else if r.Method == "DELETE" {
		_, err := db.Query("DELETE FROM task")
		if err != nil {
			log.Print("Error: ", err)
		}
		tasks := []Task{}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Method "+ r.Method + " no allowed"))	
	}
	
}

func getDB() (db *sql.DB, e error) {
	user := "root"
	pass := ""
	host := "tcp(127.0.0.1)"
	dbName := "task"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", user, pass, host, dbName))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func fethTasks() ([]Task, error) {
	tasks := []Task{}
	rows, err := db.Query("SELECT * FROM task")
    if err != nil {
		return tasks, err
    }
    defer rows.Close()
	
    for rows.Next() {
        var task Task
        if err := rows.Scan(&task.ID, &task.Content, &task.Created); err != nil {
            return tasks, err
			
        }
        tasks = append(tasks, task)
    }

    if err := rows.Err(); err != nil {
		return tasks, err
    }
	return tasks, nil
}


func get(id int64) (Task, error){
	var task Task
    row := db.QueryRow("SELECT * FROM task WHERE id = ?", id)
    if err := row.Scan(&task.ID, &task.Content, &task.Created); err != nil {
        return task, err
    }
    return task, nil
}

func add(content string, created string) (int64, error) {
	var id int64
	query, err := db.Prepare("INSERT INTO task (content, created) VALUES(?, ?)")
	if err != nil {
		return id, err
	}
	defer query.Close()

	result, err := query.Exec(content, created)
	if err != nil {
		return id, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func update(content string, id int64) ([]Task, error) {
	var tasks []Task
	query, err := db.Prepare("UPDATE task SET content=? WHERE id=?")
	if err != nil {
		return tasks, err
	}
	defer query.Close()

	_, err = query.Exec(content, id)
	if err != nil {
		return tasks, err
	}

	tasks, err = fethTasks()
	if err != nil {
		return tasks, err
    }

	return tasks, nil	
}

func delete(id int64) ([]Task, error){
	var tasks []Task

	query, err := db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		return tasks, err
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		return tasks, err
	}
		
	tasks, err = fethTasks()
	if err != nil {
        return tasks, err
    }

	return tasks, nil
}
