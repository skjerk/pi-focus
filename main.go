package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Send Frontpage to the client
func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html"))
	tmpl.Execute(w, "data goes here")
}

// apiHandler handles REST requests from the frontend
func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("apiHandler")
	if r.Method == "GET" {
		log.Println("GET")
		log.Println(r.URL.Query())
		vars, ok := r.URL.Query()["focus"]
		if !ok || len(vars[0]) < 1 {
			log.Println("Url Param 'focus' is missing")
			return
		}
		value = getFocusValue()
		return {"value": value}
	}
	if r.Method == "POST" {
		log.Println("GET")
		log.Println(r.URL.Query())
		vars, ok := r.URL.Query()["focus"]
		if !ok || len(vars[0]) < 1 {
			log.Println("Url Param 'focus' is missing")
			return
		}
		value := vars[0]
		myCommand := fmt.Sprintf("sudo v4l-ctl --set-ctrl=focus_absolute=%d", value)
		cmd := exec.Command("/bin/sh", "-c", myCommand)
		//cmd.Stdin = "raspberry" // Password for sudo
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
	//w.Write(response)
}

func getFocusValue() int64 {
	myCommand := "sudo v4l-ctl --get-ctrl=focus_absolute"
	cmd := exec.Command("/bin/sh", "-c", myCommand)
	//cmd.Stdin = "raspberry" // Password for sudo
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return
	}

	return value
}

func main() {
	http.HandleFunc("/", frontpageHandler)
	http.HandleFunc("/api/", apiHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
