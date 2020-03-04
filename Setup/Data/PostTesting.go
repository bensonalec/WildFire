package main

import (
	// "fmt"
	// "strings"
	// "net/http"
	// "io/ioutil"
)

// func setBusinessRow(name string ,url string ,email string,city string,zip string,address string, phone string) {

// 	queryurl := "http://localhost:8080/v1/graphql"
// 	method := "POST"

// 	payload := strings.NewReader(fmt.Sprintf("{\"query\":\"mutation MyMutation {\\n  insert_Business(objects: {Name: \\\"%s\\\", Address: \\\"%s\\\", City: \\\"%s\\\", Email: \\\"%s\\\", Phone: \\\"%s\\\", URL: \\\"%s\\\", Zip: \\\"%s\\\"}) {\\n    affected_rows\\n  }\\n}\\n\",\"variables\":{}}",name,address,city,email,phone,url,zip))
// 	//name,address,city,email,phone,url,zip
// 	client := &http.Client {
// 	}
// 	req, err := http.NewRequest(method, queryurl, payload)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	req.Header.Add("Content-Type", "application/json")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(body))
// }