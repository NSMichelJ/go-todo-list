package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/NSMichelJ/go-todo-list/config"
	"github.com/NSMichelJ/go-todo-list/app/handler"
)

type App struct{
	DB *sql.DB
}

func (app *App) InitDB(c *config.Config) {
	
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		c.DB.Username,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.DBName,
		c.DB.Charset,
	)

	db, err := sql.Open(c.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}
	
	app.DB = db
}

func (app *App) Run(port string){
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/", handler.MainHandler)
	http.HandleFunc("/api/v1/task/", app.handleRequest(handler.TaskHandler))
	http.HandleFunc("/api/v1/tasks/",app.handleRequest(handler.TasksHandler))

	log.Print("Server runnig on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type RequestHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (app *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(app.DB, w, r)
	}
}
