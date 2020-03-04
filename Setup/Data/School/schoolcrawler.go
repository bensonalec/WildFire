package main

import (
	"fmt"
	"github.com/gocolly/colly" 
	//"strings"
)

func main() {
	c := colly.NewCollector()
	g := colly.NewCollector()
	g.OnHTML("body",func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.OnHTML("#maincontent > div > div", func(e *colly.HTMLElement) {
		//fmt.Println(e.ChildAttr("a","href"))
		fmt.Println("FOUND")
		
	})
	c.Visit("https://www.greatschools.org/new-mexico/schools/?gradeLevels=e&view=table")
}