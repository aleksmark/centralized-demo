package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "html/template"
    "net/http"
    "os"
)

var (
    listenPort  string
    backendHost = os.Getenv("BACKEND_HOST")
    backendPort = os.Getenv("BACKEND_PORT")
)

type PageVariables struct {
    Hostname string
    Message  string
    Env      string
}

type LogRecord struct {
    Level       string `json:"level"`
    Message     string `json:"message"`
    Environment string `json:"environment"`
    Application string `json:"application"`
    Hostname    string `json:"hostname"`
}

func main() {
    flag.StringVar(&listenPort, "listen-port", "80", "server listen port")
    flag.Parse()

    myLog("info", "Server is starting on port "+listenPort)

    http.HandleFunc("/index", serveIndexPage)
    http.HandleFunc("/login", serveLoginPage)
    http.HandleFunc("/vote", serveVotePage)

    http.ListenAndServe(":"+listenPort, nil)
}

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index page")

    myLog("info", "User landed on index page")
}

func serveLoginPage(w http.ResponseWriter, r *http.Request) {
    username, _ := r.URL.Query()["username"]
    password, _ := r.URL.Query()["password"]

    resp, _ := http.Get("http://" + backendHost + ":" + backendPort + "/login?username=" + username[0] + "&password=" + password[0])

    if resp.StatusCode == 200 {
        renderTemplate(w, "Welcome User "+username[0])

        myLog("info", "User "+username[0]+" logged in")
    } else {
        renderTemplate(w, "Wrong username and password")

        myLog("error", "Invalid login attempt")
    }
}

func serveVotePage(w http.ResponseWriter, r *http.Request) {
    resp, _ := http.Get("http://" + backendHost + ":" + backendPort + "/vote")

    if resp.StatusCode == 200 {
        renderTemplate(w, "You voted successfully")

        myLog("info", "User voted")
    } else {
        renderTemplate(w, "Vote failed")

        myLog("error", "User vote failed")
    }
}

func renderTemplate(w http.ResponseWriter, message string) {
    indexPageVars := PageVariables{
        Hostname: getHostname(),
        Env:      os.Getenv("ENV"),
        Message:  message,
    }

    t, err := template.ParseFiles("index.html")
    if err != nil {
        myLog("error", "Template parsing error: "+err.Error())
    }

    err = t.Execute(w, indexPageVars)
    if err != nil {
        myLog("error", "Template executing error: "+err.Error())
    }

    serialized, _ := json.Marshal(indexPageVars)
    myLog("info", "Template rendered with data: "+string(serialized))
}

func getHostname() string {
    name, err := os.Hostname()
    if err != nil {
        panic(err)
    }

    return name
}

func myLog(level string, message string) {
    logRecord := LogRecord{
        Level:       level,
        Message:     message,
        Environment: os.Getenv("ENV"),
        Application: os.Getenv("APP"),
        Hostname:    getHostname(),
    }

    serialized, _ := json.Marshal(logRecord)

    fmt.Println(string(serialized))
}
