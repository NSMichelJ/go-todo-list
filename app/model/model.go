package model

import (
	"database/sql"
)

type Task struct{
	ID int64 `json:"id"`
	Content string `json:"content"`
	Created string `json:"created"`
}

func (task Task) FethTasks(db *sql.DB) ([]Task, error) {
	tasks := []Task{}
	rows, err := db.Query("SELECT * FROM task")
    if err != nil {
		return tasks, err
    }
    defer rows.Close()
	
    for rows.Next() {
        var t Task
        if err := rows.Scan(&t.ID, &t.Content, &t.Created); err != nil {
            return tasks, err
        }
        tasks = append(tasks, t)
    }

    if err := rows.Err(); err != nil {
		return tasks, err
    }
	return tasks, nil
}

func (task Task) Add(db *sql.DB) (int64, error) {
	var id int64
	query, err := db.Prepare("INSERT INTO task (content, created) VALUES(?, ?)")
	if err != nil {
		return id, err
	}
	defer query.Close()

	result, err := query.Exec(task.Content, task.Created)
	if err != nil {
		return id, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func (task Task) Get(db *sql.DB) (Task, error){
	var t Task
    row := db.QueryRow("SELECT * FROM task WHERE id = ?", task.ID)
    if err := row.Scan(&t.ID, &t.Content, &t.Created); err != nil {
        return t, err
    }
    return t, nil
}

func (task Task) Update(db *sql.DB) ([]Task, error) {
	var tasks []Task
	query, err := db.Prepare("UPDATE task SET content=? WHERE id=?")
	if err != nil {
		return tasks, err
	}
	defer query.Close()

	_, err = query.Exec(task.Content, task.ID)
	if err != nil {
		return tasks, err
	}

	tasks, err = task.FethTasks(db)
	if err != nil {
        return tasks, err
    }

	return tasks, nil	
}

func (task Task) Delete(db *sql.DB) ([]Task, error) {
	var tasks []Task

	query, err := db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		return tasks, err
	}
	defer query.Close()

	_, err = query.Exec(task.ID)
	if err != nil {
		return tasks, err
	}
		
	tasks, err = task.FethTasks(db)
	if err != nil {
        return tasks, err
    }

	return tasks, nil
}
