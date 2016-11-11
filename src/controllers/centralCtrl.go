package controllers

import (
	"text/template"
	"net/http"
	"fmt"
	"log"
	"time"
	"strconv"
	"strings"
	"os"
	"bufio"
	"path/filepath"
	
	"github.com/gorilla/websocket"
)

func Register(templates *template.Template) {
	
	pc := new(projectsController)
	pc.template = templates.Lookup("projects.html")
	
	http.HandleFunc("/", pc.serveHome)
	http.HandleFunc("/ws", serveWs)
	
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("ServeWS: ",r.URL.Path)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err != nil {
		lastMod = time.Unix(0, n)
	}
	
	dataFile := r.FormValue("dataFile")
	go writer(ws, lastMod, dataFile)
	reader(ws)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	basePath, err := filepath.Abs("./templates")
	basePath = filepath.ToSlash(basePath)
	if err != nil {
		fmt.Println("Could not find path for templates -", basePath)
		panic(err)
	}
	path := basePath + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else {
		contentType = "text/plain"
	}
	
	f, err := os.Open(path)
	
	if err == nil {
		defer f.Close()
		w.Header().Add("Content Type", contentType)
		
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}