package main

import(
	"log"
	"os"
	"strings"
	"io/ioutil"
	"fmt"
	"strconv"
	"time"
	"net/http"

)

func main() {
	// busCSV()
	schoolCSV()
}

func setBusinessRow(name string ,url string ,email string,city string,zip string,address string, phone string) {

	queryurl := "http://localhost:8080/v1/graphql"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("{\"query\":\"mutation MyMutation {\\n  insert_Business(objects: {Name: \\\"%s\\\", Address: \\\"%s\\\", City: \\\"%s\\\", Email: \\\"%s\\\", Phone: \\\"%s\\\", URL: \\\"%s\\\", Zip: \\\"%s\\\"}) {\\n    affected_rows\\n  }\\n}\\n\",\"variables\":{}}",name,address,city,email,phone,url,zip))
	//name,address,city,email,phone,url,zip
	client := &http.Client {
	}
	req, err := http.NewRequest(method, queryurl, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
func busCSV() {
	//first, need to parse lines in the csv
	file, err := os.Open("./data/businesses.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	b, err := ioutil.ReadAll(file)
	lines := strings.Split(string(b),"\n")
	//iterate through increments of 16382
	//first, 0 - 16382 , then 16382 - 2*16382, etc
	start := 0
	end := 16382
	for true {	
		if(end < len(lines)) {
			for ind,ele := range lines[start:end] {
				if(len(ele) > 0 && ind != 0) {
					//City,name,address,phone,sector
					info := strings.Split(ele,",")
					city := info[0]
					name := info[1]
					address := info[2]
					addSpl := strings.Split(address," ")
					zip := addSpl[len(addSpl)-2]
					finalAdd := ""
					for i := 0; i < len(addSpl)-2;i++ {
						finalAdd += addSpl[i] + " "
					}
					
					phone := info[3]
					sectors := info[4]
		
					fmt.Println(city,name,finalAdd,zip,phone,sectors)
					//check if sector already exists
						//query by name of sector for id
						//if result is null, then make a new sector, and retrieve its id
					setBusinessRow(name,"UNKNOWN","UNKNOWN",city,zip,finalAdd,phone)
					//add the business, retrieve it's id
					//then, for all sectors you have, insert into the busToSec table the business id and the sector id
					fmt.Println(ind)
				}
				time.Sleep(1 * time.Millisecond)
			}	
	
		} else {
			for ind,ele := range lines[start:] {
				if(len(ele) > 0 && ind != 0) {
					//City,name,address,phone,sector
					info := strings.Split(ele,",")
					city := info[0]
					name := info[1]
					address := info[2]
					addSpl := strings.Split(address," ")
					zip := addSpl[len(addSpl)-2]
					finalAdd := ""
					for i := 0; i < len(addSpl)-2;i++ {
						finalAdd += addSpl[i] + " "
					}
					
					phone := info[3]
					sectors := info[4]
		
					fmt.Println(city,name,finalAdd,zip,phone,sectors)
					//check if sector already exists
						//query by name of sector for id
						//if result is null, then make a new sector, and retrieve its id
					setBusinessRow(name,"UNKNOWN","UNKNOWN",city,zip,finalAdd,phone)
					//add the business, retrieve it's id
					//then, for all sectors you have, insert into the busToSec table the business id and the sector id
					fmt.Println(ind)
				}
			}	

		}

		start += 16382
		end += 16382
		
	}

}

func setSchoolRow(name string, address string, city string, url string, email string, phoneNumber string, numberOfStudents int) {


	queryurl := "http://localhost:8080/v1/graphql"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("{\"query\":\"mutation MyMutation {\\n  insert_Schools(objects: {Address: \\\"%s\\\", Attendees: %d, City: \\\"%s\\\", Email: \\\"%s\\\", Name: \\\"%s\\\", Phone: \\\"%s\\\", URL: \\\"%s\\\"}) {\\n    affected_rows\\n  }\\n}\\n\",\"variables\":{}}",address,numberOfStudents,city,email,name,phoneNumber,url))
	//address,numberOfStudents,city,email,name,phoneNumber,url
	client := &http.Client {
	}
	req, err := http.NewRequest(method, queryurl, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))}


func schoolCSV() {
	file, err := os.Open("./data/schools.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	b, err := ioutil.ReadAll(file)
	lines := strings.Split(string(b),"\n")
	for _,ele := range lines {
		if(ele != "\n") {
			fmt.Println(ele)
			info := strings.Split(ele,",")
			city := info[2]
			name := info[0]
			address := info[1]
			if(len(info) > 3) {
				enrollment,_ := strconv.Atoi(info[3])
				setSchoolRow(name,address,city,"","","",enrollment)
			}else {
				setSchoolRow(name,address,city,"","","",0)
			}
		}
	}

}