package main

import (
	"log"
	"net/http"
	"html/template"
	"encoding/gob"
	"fmt"
	"strings"
	"github.com/gorilla/sessions"
	"reflect"
	"strconv"
	// "github.com/go-ldap/ldap"
)

var store = sessions.NewCookieStore([]byte("samplekey"))
var loggedIn = false

//handle initial load, and unknown addresses
func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		gob.Register(map[string]string{})

		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		}
		session.Values["loadSize"] = 100
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		if _, ok := session.Values["loggedIn"]; ok {
			if r.Method == "GET" && session.Values["loggedIn"].(bool) {
				//Case for GET request
		
				lists := getTables()
				setSortDefaults(w,r)
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

func setSortDefaults(w http.ResponseWriter,r *http.Request) {
	session, _ := store.Get(r, "session-name")
	gob.Register(map[string]string{})
	session.Values["loadSize"] = 100

	lists := getTables()
	if _, ok := session.Values["sort"]; ok {
	
	} else {
		//here, set defaults for sort stuff
		sortValues := make(map[string]string)
		for _, ele := range lists {
			// fmt.Println(ele.TblNames.Column)
			for _,column := range ele.TblNames.Column {
				//table is column.cat
				// tableName := column.Cat
				backName := column.BackName
				sortValues[backName] = "ID,1"

			}
		}
		
		session.Values["sort"] = sortValues

	}
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		setSortDefaults(w,r)
		session, _ := store.Get(r, "session-name")
		
		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				// fmt.Println(session.Values["sort"])
				//sort values in session.Values["sort"]
				spl := strings.Split(r.URL.String(),"/")
				typ := spl[2]
				tName := spl[3]
				pagNum,err := strconv.Atoi(spl[4])
				//based on these, build the table from this info

				sortParam := session.Values["sort"].(map[string]string)[tName]
				tmpl := template.Must(template.ParseFiles("html/table.html"))
				limit := session.Values["loadSize"].(int)

				tbl := getTable(typ,tName,limit,pagNum,sortParam)
				
				err = tmpl.Execute(w,tbl)
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

func (s *server) handleSort() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSortDefaults(w,r)
		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				spl := strings.Split(r.URL.String(),"/")
				//get the tablename and the column name
				tType := spl[2]
				tName := spl[3]
				colName := spl[4]
			
				//check the table
				ind := session.Values["sort"].(map[string]string)[tName]
				indSplit := strings.Split(ind,",")
				curCol := indSplit[0]
				curInd := indSplit[1]
				
				newEntry := ""
				if(curCol == colName) {
					//increment curInd, then mod by 3
					curIntInd,_ := strconv.Atoi(curInd)
					curIntInd++
					curIntInd %= 3
					newEntry = colName + "," + strconv.Itoa(curIntInd)
				} else {
					newEntry = colName +","+"1"
				}
				session.Values["sort"].(map[string]string)[tName] = newEntry

				err := session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//need to get the table type, then redirect back to the /table/type/tname/1
				http.Redirect(w, r, "/table/"+tType+"/"+tName+"/1", http.StatusSeeOther)
				
			}
		}
	}
}

func (s *server) handleNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSortDefaults(w,r)
		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				spl := strings.Split(r.URL.String(),"/")
				if(spl[1] != "new") {

				} else {
					typ := spl[2]
					tName := spl[3]
					tmpl := template.Must(template.ParseFiles("html/newentry.html"))
					limit := session.Values["loadSize"].(int)
					sortParam := session.Values["sort"].(map[string]string)[tName]

					tbl := getTable(typ,tName,limit,1,sortParam)
					
				
					err := tmpl.Execute(w,tbl)
					if err != nil {
						fmt.Println("REEEEE")
						panic(err)
					}
	
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
		setSortDefaults(w,r)
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
		setSortDefaults(w,r)
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

func (s *server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSortDefaults(w,r)
		session, _ := store.Get(r, "session-name")
		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				if(r.Method == "POST") {
					r.ParseForm()
					if(len(r.Form) != 0) {
						// parse for the search term, then pass this to a modified get table function essentially
						searchTerm := r.Form["searchValue"][0]
						spl := strings.Split(r.URL.String(),"/")
						typ := spl[2]
						tName := spl[3]
						tbl := searchTable(typ,tName,1000,1,searchTerm)

						tmpl := template.Must(template.ParseFiles("html/table.html"))
				
				
						err := tmpl.Execute(w,tbl)
						if err != nil {
							fmt.Println("REEEEE")
							panic(err)
						}
	
					}
	
				}
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		
	}
}

func (s *server) handleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSortDefaults(w,r)
		session, _ := store.Get(r, "session-name")

		if _, ok := session.Values["loggedIn"]; ok {
			if(session.Values["loggedIn"].(bool)) {
				//here, parse a file that has been uploaded
			}
		}
	}
}



func (s *server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSortDefaults(w,r)
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
