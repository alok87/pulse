package controllers

import (
	"time"
	"fmt"
	"os"
	"io/ioutil"
	
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		//fmt.Println("Wating to read from Websocket reader")
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time, dataFile string) {
	//fmt.Println("Insider writer")
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-fileTicker.C:
			var p []byte
			var err error
			
			p, lastMod, err = readFileIfModified(lastMod, dataFile, "writerFunc")

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
				//fmt.Println("Wrote to websocket",dataFile)
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func readFileIfModified(lastMod time.Time, filename string, what string) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil { // if os stats results error then	
	    fmt.Println("os stat returned error for reading file", err)
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) { // if the file was not modified when compared to lastMod
		//fmt.Println("File not updated, nothing to do: ", what)
		//fmt.Println("Modtime ",fi.ModTime())
		//fmt.Println("LastModPasssed ",lastMod)
		//fmt.Println("")
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename) // else if the file was modified return the modified time and new data
	if err != nil {
		fmt.Println("err reading file ",filename)
		return nil, fi.ModTime(), err
	}
		//fmt.Println("File updated, data returned",what)
		//fmt.Println("Modtime ",fi.ModTime())
		//fmt.Println("LastModPasssed ",lastMod)
		//fmt.Println("")
	return p, fi.ModTime(), nil
}
