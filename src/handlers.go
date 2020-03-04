package main

import (
	"log"
	"net/http"
	"html/template"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"github.com/gorilla/sessions"
	"reflect"
	// "github.com/go-ldap/ldap"
)

var store = sessions.NewCookieStore([]byte("samplekey"))
var loggedIn = false

//handle initial load, and unknown addresses
func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		}

		fmt.Println(session.Values["loggedIn"].(bool))
		if _, ok := session.Values["loggedIn"]; ok {
			if r.Method == "GET" && session.Values["loggedIn"].(bool) {
				//Case for GET request
				db, err := sql.Open("mysql","alec:gregor@tcp(localhost:3306)/SCHOOLDB")
				if err != nil {
					fmt.Println("Opening DB failed")
					log.Fatal(err)
					
				}
				defer db.Close()		
		
				 lists := getTables()
	
				t := template.Must(template.ParseFiles("html/loggedInIndex.html"))
				if err := t.Execute(w, lists); err != nil {
					 log.Fatalln(err)
				}
			} else if r.Method == "POST" && !session.Values["loggedIn"].(bool) {
					r.ParseForm()
					if(len(r.Form) != 0) {
	
						username := (r.Form["username"][0])
						password := (r.Form["password"][0])
						loggedIn = login(username,password)
						if(loggedIn) {
							session.Values["loggedIn"] = true
							err := session.Save(r, w)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
								return
							}
							
						} else {
							http.Redirect(w,r,"/",http.StatusSeeOther)
						}		
					}
				
				http.Redirect(w, r, "/", http.StatusSeeOther)
			
			} else {
				//Case for POST request
				session.Values["loggedIn"] = false
				err := session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				t := template.Must(template.ParseFiles("html/index.html"))
				if err := t.Execute(w, nil); err != nil {
					log.Fatalln(err)
				}
	
				
			}
	
		} else {
			session.Values["loggedIn"] = false
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		
	}
}
func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("LOGOUT TRIGGERED")
		session, _ := store.Get(r, "session-name")

		session.Values["loggedIn"] = false
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
}

func (s *server) handleTableLoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		
		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				spl := strings.Split(r.URL.String(),"/")
				typ := spl[2]
				tName := spl[3]
				//based on these, build the table from this info
				tmpl := template.Must(template.ParseFiles("html/table.html"))
				tbl := getTable(typ,tName)
				
				err := tmpl.Execute(w,tbl)
				if err != nil {
					fmt.Println("REEEEE")
					panic(err)
				}

			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
				
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	
		
	}
}


func (s *server) handleNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				spl := strings.Split(r.URL.String(),"/")

				typ := spl[2]
				tName := spl[3]
				tmpl := template.Must(template.ParseFiles("html/newentry.html"))

				tbl := getTable(typ,tName)
				
			
				err := tmpl.Execute(w,tbl)
				if err != nil {
					fmt.Println("REEEEE")
					panic(err)
				}

			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
		
}
func (s *server) handleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				spl := strings.Split(r.URL.String(),"/")
				typ := spl[2]
				tName := spl[3]
				//parse form
				if(r.Method == "POST") {
					r.ParseForm()
					if(len(r.Form) != 0) {
						//figure out how to generically parse out the elements of the table
						keys := reflect.ValueOf(r.Form).MapKeys()
						strkeys := make([]string, len(keys))
						for i := 0; i < len(keys); i++ {
							strkeys[i] = keys[i].String()
						}
						cont := make([]string,len(strkeys))
						// var cont [len(strkeys)]string
						index := 0
						for key,ele := range r.Form {
							// fmt.Println("IND",index)
							// fmt.Println("ELE",ele[0])
							// index++
							index = 0
							for _,keyEle := range strkeys {
								
								if(key == keyEle) {
									cont[index] = ele[0]
									break
								} else{
									index++
								}
	
							}
							//need to make sure the keys[index] and key are the same
						}
						setRow(typ,tName,strkeys,cont)
						
					}
					//first_name,last_name,email,phoneNumber
				}

				http.Redirect(w, r, "/table/"+ typ+"/"+tName, http.StatusSeeOther)


			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		
	}
}

func (s *server) handlePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				spl := strings.Split(r.URL.String(),"/")
				tableName := spl[2]
				Id := spl[3]
				tbl := getPage(tableName,Id)
				tmpl := template.Must(template.ParseFiles("html/detailpage.html"))
				
				
				err := tmpl.Execute(w,tbl)
				if err != nil {
					fmt.Println("REEEEE")
					panic(err)
				}

			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}


	}
}

func (s *server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
		
				cond := strings.Split(r.URL.String(),"/")
				table := (cond[2])
				id := (cond[3])
				deleteRow(table,id)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
}