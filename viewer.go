package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

type Page struct {
	Title     string
	Location  string
	PageCount int
}

var db *sql.DB

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "postgres"
	password = "AAAaaa111"
	dbname   = "archive"
)

var validPath = regexp.MustCompile("^/(reader|)/([a-zA-Z0-9/-]+)/([a-zA-Z0-9/-]+)$")

func getDocumentInfo(title string) (string, int, error) {
	var path string
	var pageCount int

	err := db.QueryRow("SELECT path, page_count FROM documents WHERE title=$1", title).Scan(&path, &pageCount)

	switch {
	case err == sql.ErrNoRows:
		return "", 0, errors.New("no document with title " + title)
	case err != nil:
		log.Fatalf("query error: %v\n", err)
		return "", 0, errors.New("query error")
	default:
		return path, pageCount, nil
	}
}

func loadPage(title string, pageNumber string) (*Page, error) {
	location, pageCount, err := getDocumentInfo(title)
	fmt.Println(location, pageCount)
	location = "/library/" + location + "/" + pageNumber + ".jpg"
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Location: location, PageCount: pageCount}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", "", errors.New("Invalid Page Title")
	}
	return m[2], m[3], nil
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	title, page, err := getTitle(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	p, err := loadPage(title, page)
	if err != nil {
		p = &Page{Title: title}
	}
	t, err := template.ParseFiles("reader_page.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, p)
}

func getConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func main() {
	db = getConnection()
	defer db.Close()
	http.Handle("/library/", http.StripPrefix("/library/", http.FileServer(http.Dir("library"))))
	http.HandleFunc("/reader/", pageHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
