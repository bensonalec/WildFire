package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type topLevel struct {
	Tbls table
	TblNames []tableNames
}

type topUpload struct {
	TblNames []tableNames
	Msg string
}

type tableNames struct {
	DBName string
	TblNames tblRow
}

type table struct {
	Name string
	BackName string
	Page int
	Last int
	Type string
	Titles []doubEnt
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
type doubEnt struct {
	Name string
	Cont string
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

func getTable(typ string,tableName string,limit int, pageNum int, sortParam string) topLevel {
	//parse sort param by comma, first split is the column to sort by, and the second split is the index of what sort
	//retrieve the column names for that tableName
	offset := limit * (pageNum-1)
	

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
		tbl.Titles = append(tbl.Titles,doubEnt{Name:tableName,Cont:ele,Type:typ})
		tbl.ToAdd = append(tbl.ToAdd,ent{ele})
	}
	tbl.Name = displayName
	tbl.BackName = tableName
	tbl.Page = pageNum + 1
	tbl.Last = pageNum - 1
	tbl.Type = typ
	toGet := ""
	//write query ti graphql
	for _,ele := range columnList {
		toGet += ele +"\n"
	}
	toGet += "ID\n"

	//here, set up sort logic
	sortSplit := strings.Split(sortParam,",")
	colSort := sortSplit[0]
	sortIndex,_ := strconv.Atoi(sortSplit[1])
	sortParameters := ""
	//switch on sortIndex
	switch sortIndex {
	case 0:
		//no sort
		sortParameters = ""
	case 1:
		//ascending sort
		sortParameters = fmt.Sprintf("order_by: {%s: asc},",colSort)
	case 2:
		//descending sort
		sortParameters = fmt.Sprintf("order_by: {%s: desc},",colSort)
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%d.\n", sortIndex)
	}



	boilerPlate := fmt.Sprintf("{\"query\":\"query MyQuery {%s(%slimit: %d, offset: %d){\n%s}" + "}\",\"variables\":{}}",tableName,sortParameters,limit,offset,toGet)
	
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

func searchTable(typ string,tableName string,limit int, pageNum int,searchTerm string) topLevel {
	

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
		tbl.Titles = append(tbl.Titles,doubEnt{tableName,ele,typ})
		tbl.ToAdd = append(tbl.ToAdd,ent{ele})
	}
	tbl.Name = displayName
	tbl.BackName = tableName
	tbl.Page = pageNum + 1
	tbl.Last = pageNum - 1
	tbl.Type = typ
	toGet := ""
	//write query ti graphql
	for _,ele := range columnList {
		toGet += ele +"\n"
	}
	toGet += "ID\n"
	//build search query {Name: {_ilike: "%Aiko%"}}, {Address: {_ilike: "%Sands%"}}
	searchBoilerplate := "{_ilike: \\\"%" + searchTerm +"%\\\"}"
	searchCriteria := ""
	//iterate through all columsn that aren't integer type 
	for _,ele := range columnList{
		if(ele != "Attendees") {
			newCrit := fmt.Sprintf("{%s : %s}",ele,searchBoilerplate)
			searchCriteria += newCrit + ","
	
		}
	}
	searchCriteria = searchCriteria[:len(searchCriteria)-1]
	// fmt.Println(searchCriteria)
	boilerPlate := fmt.Sprintf("{\"query\":\"query MyQuery {%s(where: {_or: [%s]}){\n%s}" + "}\",\"variables\":{}}",tableName,searchCriteria,toGet)

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
	boilerPlate := fmt.Sprintf("{\"query\":\"mutation MyMutation {  insert_%s(objects: {%s}) {    affected_rows  }}\",\"variables\":{}}",tName,toAdd)
	makeQuery(boilerPlate)


}

func ToCSV(tableName string, searchTerm string) string{
	//basically, given the tableName and a search query, get these rows from the database, and export this data to a csv
	//so, basic strucutre is:
	//1: perform search on given table (can reuse logic from search)
	csvOutput := ""
	toSend := fmt.Sprintf("{\"query\":\"query MyQuery {\\n  Tables(where: {Name: {_eq: \\\"%s\\\"}}) {\\n    DisplayName\\n    Columns {\\n      Name\\n    }\\n  }\\n}\\n\",\"variables\":{}}",tableName)

	body := makeQuery(toSend)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	nameAndColumns := result["data"].(map[string]interface{})["Tables"].([]interface{})
	columns := nameAndColumns[0].(map[string]interface{})["Columns"].([]interface{})

	//2: get the columns for this
	var columnList []string
	for _,ele := range columns {
		//3: make line 1 of the csv the column names
		columnList= append(columnList,ele.(map[string]interface{})["Name"].(string))
		csvOutput += ele.(map[string]interface{})["Name"].(string) + ","
	}
	csvOutput = csvOutput[:len(csvOutput)-1]+"\n"
	toGet := ""
	//write query ti graphql
	for _,ele := range columnList {
		toGet += ele +"\n"
	}
	toGet += "ID\n"
	//build search query {Name: {_ilike: "%Aiko%"}}, {Address: {_ilike: "%Sands%"}}
	searchBoilerplate := "{_ilike: \\\"%" + searchTerm +"%\\\"}"
	searchCriteria := ""
	//iterate through all columsn that aren't integer type 
	for _,ele := range columnList{
		if(ele != "Attendees") {
			newCrit := fmt.Sprintf("{%s : %s}",ele,searchBoilerplate)
			searchCriteria += newCrit + ","
	
		}
	}
	searchCriteria = searchCriteria[:len(searchCriteria)-1]
	boilerPlate := fmt.Sprintf("{\"query\":\"query MyQuery {%s(where: {_or: [%s]}){\n%s}" + "}\",\"variables\":{}}",tableName,searchCriteria,toGet)

	body2 := makeQuery(boilerPlate)
	var result2 map[string]interface{}
	json.Unmarshal([]byte(body2), &result2)
	rowData := result2["data"].(map[string]interface{})[tableName].([]interface{})

	//4: now, get each line from the results

	for _,ele := range rowData {
		newLine := ""	
		for _,col := range columnList {

			x := ele.(map[string]interface{})[col]
			
			switch ty := x.(type) {
			case float64:
				//x is a float64
				// tmp.Column = append(tmp.Column,ent{fmt.Sprintf("%d",int(ty))})
				newLine += fmt.Sprintf("%d",int(ty)) + ","
			case nil:
				newLine += "nil" + ","
			default:
				//x is a string
				// tmp.Column = append(tmp.Column,ent{x.(string)})
				newLine += x.(string) + ","
			} 
			
		}
		csvOutput += newLine + "\n"
	}	
	return csvOutput

}

func ImportFromCSV(lines string) string{
	//It does, so csv should be in form:
	/*
	tableName
	info,info,info,info,info
	info,info,info,info,info
	*/
	//based on the tableName, retreive the expected columns, then read in this input
	//1: Retrieve the expected columns and types
	//2: Check that all lines in the input are of the right length and types
	//3: If all are valid, then make a new entry for each row
	//3: If any are invalid, notify the user the input in invalid, and on what line
	allLines := strings.Split(lines,"\n")
	if(len(allLines) < 2) {
		return "invalid input"
	} else {
		tableTitle := allLines[0]
		//based on table Title, make a request to get columns and types expected
		//columns have types and names
		payload := fmt.Sprintf("{\"query\":\"query MyQuery {\\n  Tables(where: {Name: {_eq: \\\"%s\\\"}}) {\\n    Columns {\\n      Name\\n      Type\\n    }\\n  }\\n}\\n\",\"variables\":{}}",tableTitle)
		results := makeQuery(payload)
		//now parse through results to get all the column entries
		var result map[string]interface{}
		json.Unmarshal([]byte(results), &result)
		columns := result["data"].(map[string]interface{})["Tables"].([]interface{})[0].(map[string]interface{})["Columns"].([]interface{})
		
		//get lenght of columns
		expectedLength := len(columns)
		var expectedTypes []string
		var keys []string
		for _,ele := range columns {
			obj := ele.(map[string]interface{})
			
			expectedTypes = append(expectedTypes,obj["Type"].(string))
			keys = append(keys,obj["Name"].(string))
		}
		//now, we know the expectedLength and the expectedTypes
		for lineNum,line := range allLines[1:] {
			//split the line by comma
			if(line != "") {
				lineSpl := strings.Split(line,",")
			
				//check that the line is the appropriate length
				if(len(lineSpl) != expectedLength) {
					return "invalid input on line " + strconv.Itoa(lineNum+2) + ", expecting " + strconv.Itoa(expectedLength) + " columns"
				} 
				//check that the all input is of the appropriate type
				for ind,ele := range lineSpl {
					//check that ele is of the type expected by expectedTypes[ind]
					if(expectedTypes[ind] == "string") {
						//good to go
					}else if(expectedTypes[ind] == "int") {
						//try to convert ele to int, if it fails then return invalid input on line lineNum, column ind
						if _, err := strconv.Atoi(ele); err != nil {
							return "invalid input on line " + strconv.Itoa(lineNum+2) + ", column " + strconv.Itoa(ind+1) +", expected an integer and did not receive one"
						}
					}
				}
	
			}
		}
		//here, if it hasn't returned all input is valid
		//so, we want to now make entries for these lines
		
		for _, line := range allLines[1:] {
			if(line != "") {
				lineSpl := strings.Split(line,",")
				var toAdd string
				for _,ele := range columns {
					colName := ele.(map[string]interface{})["Name"]
					colType := ele.(map[string]interface{})["Type"]
					//find index of colname in keys
					for ind,key := range keys {
						if key == colName {
							// fmt.Println(keys[ind],colName,values[ind],colType)
							if(colType == "string") {
								toAdd += colName.(string) + ":" + "\\\"" + lineSpl[ind] + "\\\","
							} else if(colType == "int") {
								toAdd += colName.(string) + ":" + lineSpl[ind] + ","
							}
							break
						}
					}
				}
				toAdd = toAdd[:len(toAdd)-1]
				boilerPlate := fmt.Sprintf("{\"query\":\"mutation MyMutation {  insert_%s(objects: {%s}) {    affected_rows  }}\",\"variables\":{}}",tableTitle,toAdd)
				
				makeQuery(boilerPlate)
			}
		}

	
	
		//second, need to iterate through all lines once again
	}
	return "Succesfully Added Rows"
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