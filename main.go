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
	Role string `json:"role"`
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
	
	fmt.Fprintf(w, openPage("DevX Mood Analyzer"))
    //fmt.Fprintf(w, "<H1><font color='navy'>Welcome to the DevX Mood Analyzer </font></H1><H2>")

	//display happy or sad dog
	if !beHappy { 
		fmt.Fprintf(w, sadMood())
		//fmt.Fprintf(w, "<font color='red'>")
		//fmt.Fprintf(w,"Your overall mood is not great. We hope it will get better.")
		//fmt.Fprintf(w, "</font>")
		//fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>")
		//fmt.Fprintf(w, "</H2>")
		//fmt.Fprintf(w, "<BR><font color='brown'>Aggressive mood sniffing algorithm</font><BR>")
	} else { 
		fmt.Fprintf(w, happyMood())
		//fmt.Fprintf(w, "<font color='green'>")
		//fmt.Fprintf(w,"Your mood is always happy. Good for you!")
		//fmt.Fprintf(w, "</font>")
		//fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>")
		//fmt.Fprintf(w, "</H2>")
		//fmt.Fprintf(w, "<BR><font color='brown'>Mild mood sniffing algorithm</font><BR>")
	}	
	
	//sensors activation
	fmt.Fprintf(w, "<BR><BR>")
	fmt.Fprintf(w, "<font color='purple'>/activate</font><BR>")
	fmt.Fprintf(w, processSensorActivation(10))
	
	//sensors measurements
	fmt.Fprintf(w, "<BR><BR>")
	fmt.Fprintf(w, "<font color='purple'>/measure</font>")
	fmt.Fprintf(w, processSensorsMeasurement())

	fmt.Fprintf(w, closePage())
}

func processSensorActivation(numSensors int) (htmlOutput string) {

	for i := 0; i < numSensors; i++ {
		response, err := http.Get(ACTIVATE_SENSORS_API)	
		if err != nil { 
			htmlOutput = "ERROR! in calling activate API"
			return 
		} 	 	
		defer response.Body.Close()
	}
	
	htmlOutput += "<font color='gray'>" + "Succefully activated sensors." + "</font>"
	return
}

func processSensorsMeasurement() (htmlOutput string) {
	
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

	htmlOutput += "<table border='1'>"
	
	htmlOutput += "<tr style='color:grey' align='center'>"
	htmlOutput += "<th>Sensor</th>" + "<th>Role</th>" + "<th>Mood</th></tr>"
	htmlOutput += "</tr>"

	for _, sensor := range allSensors.Sensors {
  		htmlOutput += "<tr style='color:grey' align='left'>"
		htmlOutput += "<td>" + strconv.Itoa(sensor.Id) + "</td>"
		htmlOutput += "<td>" + sensor.Role + "</td>"
		htmlOutput += "<td>" + sensor.Mood + "</td>"
		htmlOutput += "</tr>"
	}

	htmlOutput += "</table>"
	return 
}

func openPage (myHeader string) (htmlOutput string) {

	htmlOutput += "<head>"
    htmlOutput += "<meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>"
	htmlOutput += "<link href='tanzu.css' rel='stylesheet'>"
    htmlOutput += "</head>"
	htmlOutput += "<body><div class='container-fluid'>"
	htmlOutput += "<div class='row'><div class='jumbotron'>"
	htmlOutput += "<h1>" + myHeader + "</h1>"
    htmlOutput += "</div><div class='row'>"
	return
}

func sadMood () (htmlOutput string) {

	htmlOutput += "<div class='col-sm-6'><form class='form-horizontal'>"
	htmlOutput += "<div class='form-group'><div class='col-sm-offset-2 col-sm-4'>"
	htmlOutput += "<p class='panel-title'>Your overall mood is not great. We hope it will get better.</p>"
	htmlOutput += "<img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>"
	htmlOutput += "<p class='panel-body'>Mood sniffing algorithm: Aggressive</p>"
	htmlOutput += "</div></div></form></div>"
	return
}

func happyMood () (htmlOutput string) {

	htmlOutput += "<div class='col-sm-6'><form class='form-horizontal'>"
	htmlOutput += "<div class='form-group'><div class='col-sm-offset-2 col-sm-4'>"
	htmlOutput += "<p class='panel-title'>Your mood is always happy. Good for you!</p>"
	htmlOutput += "<img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>"
	htmlOutput += "<p class='panel-body'>Mood sniffing algorithm: Mild</p>"
	htmlOutput += "</div></div></form></div>"
	return
}

func closePage () (htmlOutput string) {

	htmlOutput += "</body>"
	return

}

func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
