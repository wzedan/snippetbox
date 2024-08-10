package main

import (
	"database/sql"
	"flag"
	"github.com/wzedan/snippetbox/internal/models"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

type Application struct {
	infoLog  *log.Logger
	errLog   *log.Logger
	snippets *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dbName := flag.String("db", "snippetbox.db", "Database name")
	flag.Parse()

	db, err := openDB(*dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := Application{
		errLog:   errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	log.Printf("Listening on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.buildMux(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot serve on %s", *addr)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
