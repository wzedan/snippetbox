package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := Application{
		errLog:  errorLog,
		infoLog: infoLog,
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/create", app.snippetCreate)
	mux.HandleFunc("/get", app.snippetView)

	log.Printf("Listening on %s", *addr)

	srv := &http.Server{Addr: *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot serve on %s", *addr)
	}
}
