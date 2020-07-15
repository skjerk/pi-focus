package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
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
	autofocusvalue := 1
	IPAddress := getIPAddress()

	if autofocus {
		autofocusvalue = 1
	} else {
		autofocusvalue = 0
	}

	// Setup initial data and execute the template
	data := struct {
		Focus          int64
		Autofocus      bool
		AutofocusValue int
		IPAddress      string
	}{
		Focus:          focus,
		Autofocus:      autofocus,
		AutofocusValue: autofocusvalue,
		IPAddress:      IPAddress,
	}
	tmpl.Execute(w, data)
}

// apiHandler handles REST requests from the frontend
func apiHandler(w http.ResponseWriter, r *http.Request) {
	var value int64

	if r.Method == "GET" {
		value = getFocusValue()
	}

	if r.Method == "POST" {
		r.ParseForm()

		// Set Autofocus on/off
		autofocus := r.Form.Get("autofocus")
		if autofocus != "" {
			if autofocus == "true" {
				enableAutofocus()
			} else if autofocus == "false" {
				disableAutofocus()
			}
		}
		// Set focus to absolute value
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
		Value int64
	}{
		Value: value,
	}
	jresult, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(jresult)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jresult)
}

// getFocusValue returns the current focus value from the webcam
func getFocusValue() int64 {
	myCommand := "v4l2-ctl --get-ctrl=focus_absolute"
	output, err := exec.Command("/bin/sh", "-c", myCommand).Output()
	if err != nil {
		return 0
	}
	value, err := strconv.ParseInt(strings.TrimSpace(string(output[16:])), 10, 64)
	if err != nil {
		log.Println("Error converting value to string in getFocusValue")
	}
	return value
}

// setFocusValue set the webcam focus to the specified value
// min=0 max=250
func setFocusValue(value int64) int64 {
	myCommand := fmt.Sprintf("v4l2-ctl --set-ctrl=focus_absolute=%d", value)
	err := exec.Command("/bin/sh", "-c", myCommand).Run()
	if err != nil {
		value = 0
	}
	//log.Println("setFocusValue")
	//log.Println(value)
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
	if err != nil {
		return false
	}
	value, err := strconv.ParseInt(strings.TrimSpace(string(output[12:])), 10, 64)
	if err != nil {
		log.Println("Error converting value to string in getFocusValue")
	}
	if value == 0 {
		return false
	}
	return true
}

func getIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//log.Println(ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	log.Println("No IP address fouind")
	return ""
}

// doNothing is a handler for e.g. favicon requests.
func doNothing(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/", frontpageHandler)
	http.HandleFunc("/favicon.ico", doNothing)
	http.HandleFunc("/api", apiHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	msg := fmt.Sprintf("Connect using http://%v:1080", getIPAddress())
	log.Println(msg)
	log.Fatal(http.ListenAndServe(":1080", nil))
}
