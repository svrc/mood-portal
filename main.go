package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

var IS_HAPPY = false

var OVERRIDE_BACKEND_API = true

func handler(w http.ResponseWriter, r *http.Request) {
	
	log.Println(r.RemoteAddr, r.Method, r.URL.String())
	
        fmt.Fprintf(w, "<H1><font color='navy'>Welcome to DevX on K8s demo</font></H1>")

	fmt.Fprintf(w, "<H2><font color='maroon'>What is your mood today?</font>")
	
	fmt.Fprintf(w, "</font><BR><BR>")
	
	if OVERRIDE_BACKEND_API == false {
	
		response, err := http.Get("http://devx-mood-backend.dekt-apps.serving.dekt.io/sensors-data")

		responseData, err := ioutil.ReadAll(response.Body)
		
		fmt.Fprintf(w,string(responseData))
	}
	
	if IS_HAPPY == false {
		fmt.Fprintf(w, "<font color='red'>")
		fmt.Fprintf(w,"Your current mood is sad. We hope you have a better day.")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>")
		
    } else {
		fmt.Fprintf(w, "<font color='green'>")
		fmt.Fprintf(w,"Your current mood is happy. Have an awsome rest of your day!")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>")
	}
	fmt.Fprintf(w, "</H2>")
}


func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
