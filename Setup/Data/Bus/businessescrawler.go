package main

import (
	"fmt"
	"github.com/gocolly/colly" 
	"os"
	"strings"
)

type business struct {
	city string
	name string
	address string
	phone string
	category string
	
}

func main() {
	c := colly.NewCollector()
	g := colly.NewCollector()
	var found []business
	i := 0
	j := 0
	k := 0
	
	f, err := os.Create("./" + "nmcities.csv")
	if err != nil {
		fmt.Println("Problem creating file...")
	}
	defer f.Close()
	var temp business;
	var city string;
	g.OnHTML(".vcard", func(e *colly.HTMLElement) {
		
		
		//business title
		temp.name = strings.ReplaceAll(e.ChildText("div.fn.org"),","," ")
		
		//address
		//phone number
		
		temp.phone = strings.ReplaceAll(e.ChildText(".tel"),","," ")
		temp.city = city
		found = append(found, temp)
		//category
		i++
	})
	
	g.OnHTML(".adr", func(e *colly.HTMLElement) {
		tempStr := ""
		for _, element := range e.ChildTexts("span") {
			tempStr += element + " "
		}
		found[j].address = strings.ReplaceAll(tempStr,","," ")
		j++	
	})
	
	g.OnHTML(".category", func(e *colly.HTMLElement) {
		tempStr := ""
		for _, element := range e.ChildTexts("span") {
			tempStr += element + ";"
		}
		found[k].category = strings.ReplaceAll(tempStr,","," ")
		k++

	})

	g.OnHTML("#page_content > article:nth-child(2) > nav", func(e * colly.HTMLElement) {
		links := ( e.ChildAttrs("a","href"))
		toVisit := links[len(links)-1] 
		fmt.Println(toVisit)
		g.Visit("https://us-business.info" + toVisit)
	})

	c.OnHTML(".locations > ul > li", func(e *colly.HTMLElement) {
		
		if(len(e.Text) > 1) {
			city = e.Text
			g.Visit("https://us-business.info/directory" + e.ChildAttr("a","href")[2:])
		
		}

	})
	


	c.Visit("https://us-business.info/directory/nm")
	//g.Visit("https://us-business.info/directory/ruidoso-nm/")
	outString := "City" + "," + "Name" + "," + "Address" + "," + "Phone Number" + "," + "Category"
	_, err = f.WriteString(outString + "\n") 
	if err != nil {
		fmt.Println("Problem with writing...")
	}
	f.Sync()
	for _,element := range found {
		//fmt.Println(element)
		outString := element.city + "," + element.name + "," + element.address + "," + element.phone + "," + element.category
		_, err := f.WriteString(outString + "\n") 
		if err != nil {
			fmt.Println("Problem with writing...")
		}
		f.Sync()

	}

}


/*

*/