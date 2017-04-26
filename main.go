/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
	"flag"
	"net/http"
	"text/template"
    "mime"
    "path/filepath"
    "time"
    "crypto/sha1"
    "crypto/tls"
    "net"
    "syscall"
)

var _ = net.Dial
var _ = time.Now

var build string
var version string

var addr = flag.String("addr", ":6969", "http service address")
var regenerate = flag.Bool("regen", false, "regenerate cert&key")

var homeTemplate *template.Template
var infoTemplate *template.Template

func serveHome(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        Log("POST ON HOME", r)
		http.Error(w, "Method not allowed", 405)
		return
	}
    cookie, err := r.Cookie("session")
    if err == nil && cookie.Value != "" {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        homeTemplate.Execute(w, r.Host)
    } else {
        http.Redirect(w, r, "/login.html", http.StatusFound)
    }
}

func serveAsset(w http.ResponseWriter, r *http.Request, path string) {
    if path == "" {
        path = r.URL.Path[1:]
    }
    etag := r.URL.Path[1:]+"|"+build
    w.Header().Set("Cache-Control", "public, max-age:240")
    w.Header().Set("ETag", etag)
    if match := r.Header.Get("If-None-Match"); match != "" {
        if match == etag {
            w.WriteHeader(http.StatusNotModified)
            return
        }
    }
    res, _ := Asset(path)
    ext := filepath.Ext(path)
    ext = mime.TypeByExtension(ext)
    w.Header().Set("Content-Type", ext)
    Write(w, r, res)
}

type Info struct {
    Title string
    Info string
}

func serveInfo(w http.ResponseWriter, r *http.Request, title string, info string){
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    infoTemplate.Execute(w, Info{"Pink Samurai - "+title, info})
}

func serve(w http.ResponseWriter, r *http.Request) {
    Log("SERVE", r.URL.Path)
    path := r.URL.Path[1:]
    
    w.Header().Add("Strict-Transport-Security",
        "max-age=63072000; includeSubDomains; preload")
    
    if path == "" {
        serveHome(w, r)
        
    }else if path == "build" ||
        path == "buildnumber"{
        w.Header().Set("Content-Type", "text/plain")
        Write(w, r, []byte(build))
        
    }else if path == "version" ||
        path == "versionnumber"{
        w.Header().Set("Content-Type", "text/plain")
        Write(w, r, []byte(version))
        
    } else if path == "info.html" {
        serveInfo(w, r, r.URL.Query().Get("title"), r.URL.Query().Get("info"))
        
    } else if path == "login.html" {
        if r.PostFormValue("username") != "" &&
            r.PostFormValue("username") != "Tentacule-Sama" &&
            r.PostFormValue("password") != "" {
            Log("LOGIN", r.PostFormValue("username") + " " +
                r.PostFormValue("password"))

            if sha1.Sum([]byte(r.PostFormValue("password"))) ==
                members.members[r.PostFormValue("username")] {

                expiration := time.Now().Add(1 * time.Hour)

                h := members.Add(r.PostFormValue("username"),
                        r.PostFormValue("password"))

                cookie := http.Cookie{
                    Name: "session",
                    Value: h,
                    Expires: expiration}
                http.SetCookie(w, &cookie)

                cookie = http.Cookie{
                    Name: "username",
                    Value: members.hashs[h].Username,
                    Expires: expiration}
                http.SetCookie(w, &cookie)

                http.Redirect(w, r, "/", http.StatusFound)

            }else{
                if _, ok := members.members[r.PostFormValue("username")]; ok {
                    serveInfo(w, r, "Bad Password", "Bad password for " +
                        r.PostFormValue("username") +
                        "<a href=\"/\">Retry</a>")
                    
                } else {
                    members.Add(r.PostFormValue("username"),
                        r.PostFormValue("password"))
                    members.Save()
                    serveInfo(w, r, "User registered", "User registered : " +
                        r.PostFormValue("username") + 
                        "<a href=\"/\">Login</a>")
                }
            }

        } else {
            serveAsset(w, r, "login.html")
            
        }

    } else if _, err := Asset(path); err == nil {
        serveAsset(w, r, "")

    } else {
        serve404(w, r)
    }
}

func main() {

    syscall.MustLoadDLL("kernel32.dll").MustFindProc("Beep").Call(750, 300)
    
    var err error
    flag.Parse()
    var PSServeMux *http.ServeMux = http.NewServeMux()
    PSCertificate, _ := tls.LoadX509KeyPair("server.crt", "server.key")
    var PSServer = &http.Server {
        Addr:           ":6969",
        Handler:        PSServeMux,
        ReadTimeout:    15*time.Second,
        WriteTimeout:   15*time.Second,
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{
                PSCertificate,
            },
        },
    }

    NewMemInfoDumper()
    
    PSServer.Addr = *addr
    PSServer.SetKeepAlivesEnabled(true)
    
    res, err := Asset("home.html")
    if err != nil {
        panic(err)
    }
    homeTemplate = template.Must(template.New("").Parse(string(res)))

    res, err = Asset("info.html")
    if err != nil {
        panic(err)
    }
    infoTemplate = template.Must(template.New("").Parse(string(res)))
    
    println("Pink Samurai v" + version + "b" + build)
    if *regenerate {
        genCrtAndKey()
    }
    
    hub := newHub()

    go hub.run()
    PSServeMux.HandleFunc("/", serve)
    PSServeMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    go func(){
        err = PSServer.ListenAndServeTLS("", "")
        if err != nil {
            Log("ListenAndServeTLS: ", err)
            panic(err)
        }
    }()
    select{}
}

