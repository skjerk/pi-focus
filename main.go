package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

// Define the result struct
// status can be "success" or "error"
// value is the actual focus value
var result struct {
	status string
	value  int64
}

// Send Frontpage to the client
func frontpageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.gohtml"))
	focus := getFocusValue()
	autofocus := getAutoFocusValue()
	// Setup initial data for the template
	data := struct {
		Focus     int64
		Autofocus bool
	}{
		Focus:     focus,
		Autofocus: autofocus,
	}
	tmpl.Execute(w, data)
}

// apiHandler handles REST requests from the frontend
func apiHandler(w http.ResponseWriter, r *http.Request) {
	var value int64
	log.Println("apiHandler")

	if r.Method == "GET" {
		value = getFocusValue()
	}

	if r.Method == "POST" {
		r.ParseForm()
		// Handle Autofocus
		autofocus := r.Form.Get("autofocus")
		if autofocus != "" {
			if autofocus == "true" {
				enableAutofocus()
			} else if autofocus == "false" {
				disableAutofocus()
			}
		}
		// Handle focus
		focus := r.Form.Get("focus")
		if focus != "" {
			arg, err := strconv.ParseInt(string(focus), 10, 64)
			if err != nil {
				arg = getFocusValue()
			}
			value = setFocusValue(arg)
		}
	}
	result := struct {
		value int64
	}{
		value: value,
	}
	jresult, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jresult)
}

// getFocusValue returns the current focus value from the webcam
func getFocusValue() int64 {
	myCommand := "v4l2-ctl --get-ctrl=focus_absolute"
	output, _ := exec.Command("/bin/sh", "-c", myCommand).Output()
	value, err := strconv.ParseInt(strings.TrimSpace(string(output[16:])), 10, 64)
	if err != nil {
		log.Println("Error converting value to string in getFocusValue")
	}
	return value
}

// setFocusValue set the webcam focus to the specified value
// min=0 max=250 step=5 default=0 value=35
func setFocusValue(value int64) int64 {
	myCommand := fmt.Sprintf("v4l2-ctl --set-ctrl=focus_absolute=%d", value)
	err := exec.Command("/bin/sh", "-c", myCommand).Run()
	if err != nil {
		log.Println(err)
	}
	return value
}

// Enable Autofocus
// default=1 value=0
func enableAutofocus() {
	myCommand := fmt.Sprintf("v4l2-ctl --set-ctrl=focus_auto=1")
	err := exec.Command("/bin/sh", "-c", myCommand).Run()
	if err != nil {
		log.Println(err)
	}
}

// Enable Manual focus
func disableAutofocus() {
	myCommand := fmt.Sprintf("v4l2-ctl --set-ctrl=focus_auto=0")
	err := exec.Command("/bin/sh", "-c", myCommand).Run()
	if err != nil {
		log.Println(err)
	}
}

// getAutoFocusValue returns the current autofocus value from the webcam
func getAutoFocusValue() bool {
	myCommand := "v4l2-ctl --get-ctrl=focus_auto"
	output, err := exec.Command("/bin/sh", "-c", myCommand).Output()
	value, err := strconv.ParseInt(strings.TrimSpace(string(output[12:])), 10, 64)
	if err != nil {
		log.Println("Error converting value to string in getFocusValue")
	}
	if value == 0 {
		return false
	} else {
		return true
	}
}

func main() {
	http.HandleFunc("/", frontpageHandler)
	http.HandleFunc("/api", apiHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	log.Fatal(http.ListenAndServe(":1080", nil))
}
