package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type topLevel struct {
	Tbls table
	TblNames []tableNames
}

type tableNames struct {
	DBName string
	TblNames tblRow
}

type table struct {
	Name string
	BackName string
	Type string
	Titles []ent
	ToAdd []ent
	Drop []dropEnt
	Rows []rowEnt
}

type tblRow struct {
	Column []tblTit
}

type tblTit struct {
	Cat string
	Content string
	BackName string
}

type row struct {
	Column []ent
}

type rowEnt struct {
	Column []ent
	ID string
	Type string
}

type ent struct {
	Content string
}

type dropEnt struct {
	Content string
	Options []idEnt
}

type idEnt struct {
	Content string
	ID string
}
type complexent struct {
	Content string
	Name string
	ID string
	Table string
}

type multidoub struct {
	Name ent
	Cont complexent
}

type doub struct {
	Name ent
	Cont ent
}

type trip struct {
	Name ent
	Cont []complexent
}

type top struct {
	TblNames []tableNames
	Type string
	ID string
	Metadata []doub
	SingleRel []multidoub
	DoubleRel []trip
	Table string
}


func getPage(tableName string, ID string) top{
	fmt.Println(tableName)
	fmt.Println(ID)
	//get columns for this table
	toSend := fmt.Sprintf("{\"query\":\"query MyQuery {\\n  Tables(where: {Name: {_eq: \\\"%s\\\"}}) {\\n    DisplayName\\n    Columns {\\n      Name\\n    }\\n  }\\n}\\n\",\"variables\":{}}",tableName)

	body := makeQuery(toSend)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	nameAndColumns := result["data"].(map[string]interface{})["Tables"].([]interface{})
	
	displayName := nameAndColumns[0].(map[string]interface{})["DisplayName"].(string)
	columns := nameAndColumns[0].(map[string]interface{})["Columns"].([]interface{})

	//break this up, we now have the display name "DisplayName" and the lsit of Columns "Columns" and their names "Name"

	var columnList []string
	for _,ele := range columns {
		//add the column names to list
		columnList= append(columnList,ele.(map[string]interface{})["Name"].(string))
	}
	toGet := ""
	//write query ti graphql
	for _,ele := range columnList {
		toGet += ele +"\n"
	}
	toGet += "ID\n"

	boilerPlate := fmt.Sprintf("{\"query\":\"query MyQuery {%s(where: {ID: {_eq: %s}}){\n%s}" + "}\",\"variables\":{}}",tableName,ID,toGet)
	
	body2 := makeQuery(boilerPlate)

	var result2 map[string]interface{}
	json.Unmarshal([]byte(body2), &result2)
	elements := result2["data"].(map[string]interface{})[tableName].([]interface{})[0].(map[string]interface{})

	var out top 
	out.TblNames = getTables()
	out.Type = displayName
	out.ID = ID
	var tempDoub []doub
	
	for _,ele := range columnList {
		switch ty := elements[ele].(type) {
		case float64:
			//x is a float64
			tempDoub = append(tempDoub,doub{ent{ele},ent{fmt.Sprintf("%d",int(ty))}})
			
		case nil:
			fmt.Println("NIL")
		default:
			//x is a string
			tempDoub = append(tempDoub,doub{ent{ele},ent{ty.(string)}})
			} 
	}
	out.Metadata = tempDoub
	//get to the array of fields

	// for _,ele := range 

	//then, get the appropraite entry based on ID and table ID
	//then, get the foreign relations
	return out
}

func getTable(typ string,tableName string) topLevel {
	//retrieve the column names for that tableName
	var tbl table
	var tlTab topLevel

	toSend := fmt.Sprintf("{\"query\":\"query MyQuery {\\n  Tables(where: {Name: {_eq: \\\"%s\\\"}}) {\\n    DisplayName\\n    Columns {\\n      Name\\n    }\\n  }\\n}\\n\",\"variables\":{}}",tableName)

	body := makeQuery(toSend)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	nameAndColumns := result["data"].(map[string]interface{})["Tables"].([]interface{})
	
	displayName := nameAndColumns[0].(map[string]interface{})["DisplayName"].(string)
	columns := nameAndColumns[0].(map[string]interface{})["Columns"].([]interface{})

	//break this up, we now have the display name "DisplayName" and the lsit of Columns "Columns" and their names "Name"

	var columnList []string
	for _,ele := range columns {
		//add the column names to list
		columnList= append(columnList,ele.(map[string]interface{})["Name"].(string))
	}

	//iterate through this list, append the values to tbl.Titles
	for _,ele := range columnList {
		tbl.Titles = append(tbl.Titles,ent{ele})
		tbl.ToAdd = append(tbl.ToAdd,ent{ele})
	}
	tbl.Name = displayName
	tbl.BackName = tableName
	tbl.Type = typ
	toGet := ""
	//write query ti graphql
	for _,ele := range columnList {
		toGet += ele +"\n"
	}
	toGet += "ID\n"

	boilerPlate := "{\"query\":\"query MyQuery {" + tableName + "{\n" + toGet + "}" + "}\",\"variables\":{}}"
	
	body2 := makeQuery(boilerPlate)
	var result2 map[string]interface{}
	json.Unmarshal([]byte(body2), &result2)
	rowData := result2["data"].(map[string]interface{})[tableName].([]interface{})
	for _,ele := range rowData {
		var tmp rowEnt
		for _,col := range columnList {

			x := ele.(map[string]interface{})[col]
			
			switch ty := x.(type) {
			case float64:
				//x is a float64
				tmp.Column = append(tmp.Column,ent{fmt.Sprintf("%d",int(ty))})
				
			case nil:
				fmt.Println("NIL")
			default:
				//x is a string
				tmp.Column = append(tmp.Column,ent{x.(string)})
			} 
			
		}
		tmp.ID = fmt.Sprintf("%d",int(ele.(map[string]interface{})["ID"].(float64)))
		tmp.Type = tableName
		tbl.Rows = append(tbl.Rows,tmp)
		
	}
	
	tlTab = topLevel{Tbls:tbl,TblNames:getTables()}
	

	return tlTab
}

func getTables() []tableNames{


	toSend := "{\"query\":\"query MyQuery {\\n  Types {\\n    Name\\n    Tables {\\n      DisplayName\\n      Name\\n    }\\n  }\\n}\\n\",\"variables\":{}}"
	body := makeQuery(toSend)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	typesAndTabs := result["data"].(map[string]interface{})["Types"].([]interface{})
	var tblTotal []tableNames

	for _,ele := range typesAndTabs {
		Type := ele.(map[string]interface{})["Name"].(string)
		tablesIn := ele.(map[string]interface{})["Tables"].([]interface{})
		var tables tblRow
		for _,table := range tablesIn {
			tableName := table.(map[string]interface{})["DisplayName"].(string)
			officName := table.(map[string]interface{})["Name"].(string)
			tables.Column = append(tables.Column,tblTit{Type,tableName,officName})
		}
		tblTotal = append(tblTotal,tableNames{Type,tables})
	}
	return tblTotal
}

func setRow(typ string, tName string, keys []string, values []string) {
	toSend := fmt.Sprintf("{\"query\":\"query MyQuery {\\n  Tables(where: {Name: {_eq: \\\"%s\\\"}}) {\\n    Columns {\\n      Type\\n      Name\\n    }\\n    Name\\n  }\\n}\\n\",\"variables\":{}}",tName)
	body := makeQuery(toSend)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	columns := result["data"].(map[string]interface{})["Tables"].([]interface{})[0].(map[string]interface{})["Columns"].([]interface{})
	toAdd := ""
	for _,ele := range columns {
		colName := ele.(map[string]interface{})["Name"]
		colType := ele.(map[string]interface{})["Type"]
		//find index of colname in keys
		for ind,key := range keys {
			if key == colName {
				// fmt.Println(keys[ind],colName,values[ind],colType)
				if(colType == "string") {
					toAdd += colName.(string) + ":" + "\\\"" + values[ind] + "\\\","
				} else if(colType == "int") {
					toAdd += colName.(string) + ":" + values[ind] + ","
				}
				break
			}
		}
	}

	toAdd = toAdd[:len(toAdd)-1]
	fmt.Println(toAdd)
	boilerPlate := fmt.Sprintf("{\"query\":\"mutation MyMutation {  insert_%s(objects: {%s}) {    affected_rows  }}\",\"variables\":{}}",tName,toAdd)
	fmt.Println(boilerPlate)
	execute := makeQuery(boilerPlate)
	fmt.Println(execute)


}

func deleteRow(table string, id string) {
	//build the query, execute it
	boilerplate := fmt.Sprintf("{\"query\":\"mutation MyMutation {\\n  delete_%s(where: {ID: {_eq: %s}}) {\\n    affected_rows\\n  }\\n}\\n\",\"variables\":{}}",table,id)
	makeQuery(boilerplate)
}

func makeQuery(query string) string{

	url := "http://localhost:8080/v1/graphql"
	method := "POST"
  
	payload := strings.NewReader(query)
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)
  
	if err != nil {
	  fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
  
	res, err := client.Do(req)
	if(err != nil) {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
  
	return string(body)
}

func login(username string, password string) bool {
	return true
}

// func login(username string, password string) bool {
// 	//bind info
// 	bindusername := "CN=DBvis,OU=CC Service Accounts (Current),DC=cc,DC=nmt,DC=edu"
// 	bindpassword := "Gb2dpr6Pne9qvX%GYqA48nERqgzAVJ"

// 	//dial into server over TCP without TLS
// 	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "192.168.200.10", 389))
// 	if err != nil {
// 	    log.Fatal(err)
// 	}
// 	defer l.Close()

// 	// First bind with a read only user
// 	err = l.Bind(bindusername, bindpassword)
// 	if err != nil {
// 		fmt.Println("Can't bind")
// 		log.Println(err)
// 		return false
// 	}

// 	// Search for the given username
// 	searchRequest := ldap.NewSearchRequest(
// 	    "dc=cc,dc=nmt,dc=edu",
// 	    ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
// 	    fmt.Sprintf("(&(sAMAccountName=%s))", username),
// 	    []string{"dn"},
// 	    nil,
// 	)

// 	sr, err := l.Search(searchRequest)
// 	if err != nil {
// 		fmt.Println("Error in search")
// 		log.Println(err)
// 		return false
// 	}
	
// 	if len(sr.Entries) != 1 {
// 		for _,i := range sr.Entries {
// 			fmt.Println(i,"||",i.DN)
// 		}
// 		log.Println("User does not exist or too many entries returned")
// 		return false
// 	}

	
// 	userdn := sr.Entries[0].DN

// 	// Bind as the user to verify their password
// 	err = l.Bind(userdn, password)
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}

// 	// Rebind as the read only user for any further queries
// 	err = l.Bind(bindusername, bindpassword)
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	return true
// }