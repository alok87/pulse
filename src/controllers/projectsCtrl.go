package controllers

import (
	"text/template"
	"net/http"
	"time"
	"strconv"
	"fmt"
	"bytes"
	"path/filepath"
	
	"github.com/organizationofs/pulse/src/viewmodels"
)

type projectsController struct {
	template *template.Template
}

func (this *projectsController) serveHome(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vm := viewmodels.GetProjects()
	
	var projectVM []viewmodels.Project
	
	for _, project := range vm.Projects {
		
		//Find BasePath of the project
		basePath, err := filepath.Abs(".")
		if err != nil {
			fmt.Println("Could not find path for templates -", basePath)
			panic(err)
		}
		
		//Update the project's data
		project.Host = req.Host
		project.EnvCount = len(project.Environments)
		
		//Update the environment's data
		var environmentVM []viewmodels.Environment
		for _, env := range project.Environments {
			
			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprint(project.Id,"/",env.Id))
			env.Id=buffer.String()
			
			env.Host = req.Host
			env.AppCount = len(env.Applications)
			project.EnvCount = len(env.Applications) * len(project.Environments) 
			
			//Uppdate app's data
			var applicationVM []viewmodels.Application
			for _, app := range env.Applications {
				var appFile bytes.Buffer
				dataPath := fmt.Sprint("data/",env.Id,"/",app.Id,"/",app.Id,".dat")
				htmlPath := fmt.Sprint("data/",env.Id,"/",app.Id,"/",app.Id,".html")
				fullPath := fmt.Sprint(basePath, "/", dataPath)
				fullPath = filepath.ToSlash(fullPath)
				appFile.WriteString(fullPath)
				app.DataFile = appFile.String()
				app.HtmlPath = htmlPath
				localNexusFile := fmt.Sprint(app.DataFile, ".nexus")
				p, lastMod, err := readFileIfModified(time.Time{}, app.DataFile, "serveHome: app readFile()")
				if err != nil {
					p = []byte(err.Error())
					lastMod = time.Unix(0, 0)
				}
				app.Data = string(p)
				var buffer bytes.Buffer
				buffer.WriteString(fmt.Sprint(project.Id,"/",env.Id,"/",app.Id))
				app.Id=buffer.String()
				app.LastMod = strconv.FormatInt(lastMod.UnixNano(), 16)
				app.Host = req.Host
				applicationVM = append(applicationVM, app)
				// Poll nexus to update the local data for the application
				go pollNexus(dataPath, app.DataFile, localNexusFile) 
			}
			env.Applications = applicationVM
			environmentVM = append(environmentVM, env)
		}
		project.Environments = environmentVM
		projectVM = append(projectVM, project)
	}
	vm.Projects = projectVM
	this.template.Execute(w, &vm)
}
