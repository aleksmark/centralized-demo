package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "net/http"
    "os"
)

var (
	listenPort string
)

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

    http.HandleFunc("/login", login)
    http.HandleFunc("/vote", vote)

    http.ListenAndServe(":"+listenPort, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
    username, _ := r.URL.Query()["username"]
    password, _ := r.URL.Query()["password"]

    myLog("info", "Login attempt for user "+username[0]+" with password "+password[0])

    if username[0] == "admin" && password[0] == "admin" {
        w.WriteHeader(http.StatusOK)

        myLog("info", "User "+username[0]+" logged in with password "+password[0])
    } else {
        w.WriteHeader(http.StatusForbidden)

        myLog("error", "Invalid password for user "+username[0])
    }
}

func vote(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)

    myLog("info", "User X voted")
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
