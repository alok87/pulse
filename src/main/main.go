package main

import (
	"net/http"
	"text/template"
	"os"
	"log"
	"flag"
	"path/filepath"
	"fmt"
	
	"github.com/organizationofs/pulse/src/controllers"
)

var (
	addr = flag.String("addr", ":8081", "http service address")
)	

func main() {
	flag.Parse()
	templates := populateTemplates()
	
	controllers.Register(templates)
	
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func populateTemplates() *template.Template {
	result := template.New("templates")
	
	// basePath := "templates"
	basePath, err := filepath.Abs("./templates")
	basePath = filepath.ToSlash(basePath)
	fmt.Println("Basepath is=",basePath)
	if err != nil {
		fmt.Println("Could not find path for templates -", basePath)
		panic(err)
	}
	templateFolder, err := os.Open(basePath)
	if err != nil {
		fmt.Println("Could not open templates - ", basePath)
		panic(err)
	}
	defer templateFolder.Close()
	
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	
	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath + "/" + pathInfo.Name())
		}
	}
	
	result.ParseFiles(*templatePaths...)
	
	return result
}
