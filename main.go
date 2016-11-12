/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
    "mime"
    "path/filepath"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTemplate = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
        res, err := Asset(r.URL.Path[1:])
        if err != nil {
            //http.ServeFile(w, r, r.URL.Path[1:])
            http.Error(w, "Not found", 404)
        } else {
            ext := filepath.Ext(r.URL.Path[1:])
            ext = mime.TypeByExtension(ext)
            w.Header().Set("Content-Type", ext+"; charset=utf-8")
            w.Write(res)
        }
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    homeTemplate = template.Must(template.ParseFiles("home.html"))
	homeTemplate.Execute(w, r.Host)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
