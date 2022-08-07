package main


import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type Sensor struct {
	Id int `json:"id"`
	Team string `json:"team"`
	Mood string `json:"mood"`
}

type AllSensors struct {
	Sensors []*Sensor
}

var ACTIVATE_SENSORS_API string = "http://mood-sensors.dev.dekt.io/activate"
var MEASURE_SENSORS_API string = "http://mood-sensors.dev.dekt.io/measure"

func handler(w http.ResponseWriter, r *http.Request) {
	
	//conrtol the mood sniffing algorithm intensity
	beHappy := false

	log.Println(r.RemoteAddr, r.Method, r.URL.String())
	
    fmt.Fprintf(w, "<H1><font color='navy'>Welcome to the DevX Mood Analyzer </font></H1><H2>")

	//display happy or sad dog
	if !beHappy { 
		fmt.Fprintf(w, "<font color='red'>")
		fmt.Fprintf(w,"Your overall mood is not great. We hope it will get better.")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>")
		fmt.Fprintf(w, "</H2>")
		fmt.Fprintf(w, "<BR><font color='brown'>Aggressive mood sniffing algorithm</font><BR>")
	} else { //always happy
		fmt.Fprintf(w, "<font color='green'>")
		fmt.Fprintf(w,"Your mood is always happy. Good for you!")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>")
		fmt.Fprintf(w, "</H2>")
		fmt.Fprintf(w, "<BR><font color='brown'>Mild mood sniffing algorithm</font><BR>")
	}	
	
	//sensors activation
	fmt.Fprintf(w, "<BR><BR>")
	fmt.Fprintf(w, "<font color='purple'>/activate</font><BR>")
	fmt.Fprintf(w, "<font color='gray'>")
	fmt.Fprintf(w, processSensorActivation(10))
	fmt.Fprintf(w, "</font>")
	
	//sensors measurements
	fmt.Fprintf(w, "<BR><BR>")
	fmt.Fprintf(w, "<font color='purple'>/measure</font>")
	fmt.Fprintf(w, "<font color='gray'>")
	fmt.Fprintf(w, processSensorsMeasurement())
	fmt.Fprintf(w, "</font>")
}

func processSensorActivation(numSensors int) (htmlOutput string) {

	for i := 0; i < numSensors; i++ {
		err := http.Get(ACTIVATE_SENSORS_API)	
		if err != nil { 
			htmlOutput = "ERROR! in calling activate API"
		return 
		} 	 	
	}
	
	htmlOutput += "Succefully activated " + strconv.Itoa(numSensors) + " sensors."
	return
}

func processSensorsMeasurement() (htmlOutput string) {
	
	err := http.Get(MEASURE_SENSORS_API)	 

	if err != nil { 
		htmlOutput = "ERROR! in calling measure API"
		return 
	} 	 	

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body) 	

	if err != nil { 	
		htmlOutput = "ERROR! in reading response from measure API"
		return
	}

	var allSensors AllSensors
	json.Unmarshal(responseData, &allSensors.Sensors)

	htmlOutput += "<table>"
	
	htmlOutput += "<tr>"
	htmlOutput += "<th>Sensor ID</th>" + "<th>Team</th>" + "<th>Mood</th></tr>"
	htmlOutput += "</tr>"

	for _, sensor := range allSensors.Sensors {
  		htmlOutput += "<tr>"
		htmlOutput += "<td>" + strconv.Itoa(sensor.Id) + "</td>"
		htmlOutput += "<td>" + sensor.Team + "</td>"
		htmlOutput += "<td>" + sensor.Mood + "</td>"
		htmlOutput += "</tr>"
	}

	htmlOutput += "</table>"
	return 
}
	
func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
