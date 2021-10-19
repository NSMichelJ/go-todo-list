package handler

import (
	"net/http"
	"html/template"
	"database/sql"
	"encoding/json"
	"strconv"
	"log"
)

type Task struct{
	ID int64 `json:"id"`
	Content string `json:"content"`
	Created string `json:"created"`
}

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

		task, err := get(db, id)
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

		tasks, err := delete(db, id)
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
		
		tasks, err := update(db, params.Content, params.ID)
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
		tasks, err := fethTasks(db)
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

		id, err := add(db, params.Content, params.Created)
		if err != nil {
			log.Print("Error: ", err)
			return
		}

		task, err := get(db, id)
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

func fethTasks(db *sql.DB) ([]Task, error) {
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


func get(db *sql.DB, id int64) (Task, error){
	var task Task
    row := db.QueryRow("SELECT * FROM task WHERE id = ?", id)
    if err := row.Scan(&task.ID, &task.Content, &task.Created); err != nil {
        return task, err
    }
    return task, nil
}

func add(db *sql.DB, content string, created string) (int64, error) {
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

func update(db *sql.DB, content string, id int64) ([]Task, error) {
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

	tasks, err = fethTasks(db)
	if err != nil {
		return tasks, err
    }

	return tasks, nil	
}

func delete(db *sql.DB, id int64) ([]Task, error){
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
		
	tasks, err = fethTasks(db)
	if err != nil {
        return tasks, err
    }

	return tasks, nil
}
