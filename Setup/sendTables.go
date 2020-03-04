package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "log"
  "os"

)

func main() {
	file, err := os.Open("description.tbl")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	b, err := ioutil.ReadAll(file)
	lines := strings.Split(string(b),"\n")
	
	ID := 0
	for it,el := range lines {
		fmt.Println("Index",it)
		lineSpl := strings.Split(el,",")
		if((it+1) % 2)==1 {
			numCols,_ := strconv.Atoi(lineSpl[1])
			//every odd value, Parse and create a table
			ID = (addTable(lineSpl[0],numCols,lineSpl[2]))
			fmt.Println(ID)
		}else {
			//every even value, Parse and add columns for each entry
			
			for _,elSpl := range lineSpl {
				fmt.Println("Adding column",elSpl)
				addColumns(elSpl,false,ID)
			}
		}
		

	}

}

func addTable(name string,col int,DisplayName string) int{
		//parse through description.tbl, for every odd entry add a table entry, and retrieve the ID of the newly added
	//then, Add the columns for the even rows

	url := "http://localhost:8080/v1/graphql"
	method := "POST"
	toSend := fmt.Sprintf("{\"query\":\"mutation MyMutation {\\n  insert_Tables(objects: {Name: \\\"%s\\\", NumberOfColumns: %d, DisplayName: \\\"%s\\\"}) {\\n    affected_rows\\n    returning {\\n      ID\\n    }\\n  }\\n}\\n\",\"variables\":{}}",name,col,DisplayName)
	
	payload := strings.NewReader(toSend)

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
	//parse through body, get ID
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	ID := result["data"].(map[string]interface{})["insert_Tables"].(map[string]interface{})["returning"].([]interface{})[0].(map[string]interface{})["ID"]
	

	// fmt.Println(string(body))
	intID := int(ID.(float64))
	return intID
}

func addColumns(name string,hidden bool,tableID int) {
	url := "http://localhost:8080/v1/graphql"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("{\"query\":\"mutation MyMutation {\\n  insert_Columns(objects: {Hidden: %t, Name: \\\"%s\\\", tableID: %d}) {\\n    affected_rows\\n  }\\n}\\n\",\"variables\":{}}",hidden,name,tableID))
  
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
  
	fmt.Println(string(body))
}