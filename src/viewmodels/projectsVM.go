package viewmodels

import (

)

type Projects struct {
	Title string
	Active string
	Projects []Project
}

type Project struct {
	Id string
	Title string
	Host string
	Data string
	EnvCount int
	Environments []Environment 
}

type Environment struct {
	Id string
	Title string
	Host string
	Data string
	AppCount int
	Applications []Application
}

type Application struct {
	Id string
	Title string
	Host string
	Data string
	LastMod string
	HtmlPath string
	DataFile string
}

func GetProjects() Projects {
	result := Projects{
		Title: "Pulse",
		Active: "",
	}

	
	projectResult := []Project {
		platformProject,
	}

	result.Projects = projectResult
	
	return result
}
