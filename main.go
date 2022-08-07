package main


import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
	fmt.Fprintf(w, "<font color='purple'>/activate:	</font>")
	fmt.Fprintf(w, "<font color='gray'>")
	for i := 1; i < 11; i++ {
		fmt.Fprintf(w,processSensorActivation())
	}

	fmt.Fprintf(w,"<BR><BR><BR><BR>")
	
	for i := 1; i < 11; i++ {
		http.Get(ACTIVATE_SENSORS_API)
	}
	
	fmt.Fprintf(w, "[{\"sensorsStatus\":\"activated\"}]")
	fmt.Fprintf(w, "</font>")
	
	//sensors measurements
	fmt.Fprintf(w, "<BR><BR>")
	fmt.Fprintf(w, "<font color='purple'>/measure: </font>")
	fmt.Fprintf(w, "<font color='gray'>")
	
	fmt.Fprintf(w,processSensorsMeasurement())

	fmt.Fprintf(w,"<BR><BR><BR><BR>")
	response, err := http.Get(MEASURE_SENSORS_API)
	if err != nil {
		fmt.Fprintf(w,"ERROR! in calling measure API")
	} 
	
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(w,"ERROR! in reading response from measure API")
	}

	fmt.Fprintf(w,string(responseData))
	
	fmt.Fprintf(w, "</font>")
	
}

func processSensorActivation () (htmlOutput string) {
	
	response, err := http.Get(ACTIVATE_SENSORS_API)	
	if err != nil { 
		htmlOutput = "ERROR! in calling activate API"
		return 
	} 	 	
		
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body) 	
	if err != nil { 	
		htmlOutput = "ERROR! in reading response from activate API"
		return
	}
	
	htmlOutput += string(responseData)
	return
}

func processSensorsMeasurement () (htmlOutput string) {
	
	response, err := http.Get(MEASURE_SENSORS_API)	 

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
	//htmlOutput += "<tr><th><b>Sensor ID</b></th><th><b>Team</b></th><th><b>Mood</b></th></tr>"
	htmlOutput += "<tr><th><b>Team</b></th><th><b>Mood</b></th></tr>"

	for _, sensor := range allSensors.Sensors {
  		//htmlOutput += "<tr><th>"
  		//htmlOutput += sensor.Id
		htmlOutput += "</th><th>"
  		htmlOutput += sensor.Team
		htmlOutput += "</th><th>"
  		htmlOutput += sensor.Mood
		htmlOutput += "</th></tr>"
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
