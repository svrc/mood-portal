package main


import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

var ALWAYS_HAPPY = false


func handler(w http.ResponseWriter, r *http.Request) {
	
	log.Println(r.RemoteAddr, r.Method, r.URL.String())
	
        fmt.Fprintf(w, "<H1><font color='navy'>Welcome to the DevX Mood Analyzer </font></H1><H2>")

	if ALWAYS_HAPPY == false {
		fmt.Fprintf(w, "<font color='red'>")
		fmt.Fprintf(w,"Your overall mood is not great. We hope it will get better.")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>")
	} else {
		fmt.Fprintf(w, "<font color='green'>")
		fmt.Fprintf(w,"Your mood is always happy. Good for you!")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>")
	}
	fmt.Fprintf(w, "</H2>")
	
	//call api to write sensor data backend-api and display sensor data
	for i := 1; i < 5; i++ {
    		http.Get("http://devx-mood-backend.dekt-apps.serving.dekt.io/write")
	}
	//call api to read sensor data and display it
	fmt.Fprintf(w, "<BR><BR>Mood sensors raw data: ")
	response, err := http.Get("http://devx-mood-backend.dekt-apps.serving.dekt.io/sensors-data")
	if err != nil {
		fmt.Fprintf(w,"ERROR! in calling API")
	} else {
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Fprintf(w,"ERROR! in reading body")
		} else {
			fmt.Fprintf(w,string(responseData))
		}	
	}
}


func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
