package main

import (
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

type server struct {}

func (s *server) serve() {
	log.Fatalln(http.ListenAndServe(":420", nil))
}

func (s *server) routes() {
	http.HandleFunc("/", s.handleIndex())
	http.HandleFunc("/table/",s.handleTableLoad())
	http.HandleFunc("/new/",s.handleNew())
	http.HandleFunc("/add/",s.handleAdd())
	http.HandleFunc("/rec/",s.handlePage())
	http.HandleFunc("/del/",s.handleDelete())
	http.HandleFunc("/upload/",s.handleUpload())
	http.HandleFunc("/sort/",s.handleSort())
	http.HandleFunc("/search/",s.handleSearch())
	http.HandleFunc("/logout/",s.handleLogout())

	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs)) 
	
}


func main() {
	serv := server{}
	serv.routes()
	serv.serve()
}

