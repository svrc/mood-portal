package main


import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var SENSORS_BATCH int = 20
var AGRRESSIVE_HAPPY_THRESHOLD float32 = 0.1
var MILD_HAPPY_THRESHOLD float32 = 0.9

type Sensor struct {
	Id int `json:"id"`
	Role string `json:"role"`
	Mood string `json:"mood"`
	Legacy string `json:"legacy"`
}

type AllSensors struct {
	Sensors []*Sensor
}

var AllSensorsData AllSensors

func handler(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, addHeader("DevX Mood Analyzer"))

	//process APIs calls
	if processSensorActivation() != "success" {
		return
	}
	
	if processSensorsMeasurement() != "success" {
		return
	}



	happyRatio: calculateHappyRatio()

	//render dog section
	sniffLevel := os.Getenv("SNIFF_LEVEL")
	if sniffLevel == "2" {

		if happyRatio > AGRRESSIVE_HAPPY_THRESHOLD {
			fmt.Fprintf(w, addHappyDog())
		} else {
			fmt.Fprintf(w, addSadDog())
		}
		fmt.Fprintf(w,addDataTitle("mood sniffing level"))
		fmt.Fprintf(w,addDataContent("2 (agressive)"))
	}
	
	if sniffLevel == "1" {
		if happyRatio > MILD_HAPPY_THRESHOLD {
			fmt.Fprintf(w, addHappyDog())
		} else {
			fmt.Fprintf(w, addSadDog())
		}
		fmt.Fprintf(w,addDataTitle("mood sniffing level"))
		fmt.Fprintf(w,addDataContent("1 (mild)"))
	}
		
	
	//render API section
	fmt.Fprintf(w,addDataTitle("/activate"))
	fmt.Fprintf(w,addDataContent("All sensors activated successfully"))
	fmt.Fprintf(w,addDataTitle("/measure"))
	fmt.Fprintf(w,addDataContent(createResultsTable()))
}

func processSensorActivation() (status string) {

	tlsConfig := &http.Transport{
	 	TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	 }


	tlsClient := &http.Client{Transport: tlsConfig}
	for i := 0; i < SENSORS_BATCH ; i++ {
		response, err := tlsClient.Get(os.Getenv("SENSORS_ACTIVATE_API"))	
		if err != nil { 
			status = "Error in calling activate API: " + err.Error()
		} 	 	
		defer response.Body.Close()
	}
	status = "success"
	return
}

func processSensorsMeasurement() (status string) {
	
	tlsConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}


	tlsClient := &http.Client{Transport: tlsConfig}

	response, err := tlsClient.Get(os.Getenv("SENSORS_MEASURE_API"))	 

	if err != nil { 
		status = "Error in calling measure API: " + err.Error()
	} 	 	

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body) 	

	if err != nil { 	
		status = "Error in reading measure results: " + err.Error()
	}

	json.Unmarshal(responseData, &AllSensorsData.Sensors)

	status = "success"
	return
}

func calculateHappyRatio () (happyRatio float32){
	
	numHappy := 0
	for _, sensor := range AllSensorsData.Sensors {
		if sensor.Mood == "happy" {
			numHappy++
		}
	}
	happyRatio=numHappy/SENSORS_BATCH
	return
}

func createResultsTable () (htmlOutput string) {

	htmlOutput += "<table border='1'>"
	
	htmlOutput += "<tr style='color:grey' align='center'>"
	htmlOutput += "<th>Sensor</th>" + "<th>Role</th>" + "<th>Current Mood</th>"+ "<th>Pre-Existing</th>"
	htmlOutput += "</tr>"

	for _, sensor := range AllSensorsData.Sensors {
  		htmlOutput += "<tr style='color:grey' align='left'>"
		htmlOutput += "<td>" + strconv.Itoa(sensor.Id) + "</td>"
		htmlOutput += "<td>" + sensor.Role + "&nbsp;</td>"
		htmlOutput += "<td>" + sensor.Mood + "&nbsp;</td>"
		htmlOutput += "<td>" + sensor.Legacy + "</td>"
		htmlOutput += "</tr>"
	}

	htmlOutput += "</table>"
	return
}

func addHeader (myHeader string) (htmlOutput string) {

    htmlOutput += "<H1><font color='navy'>"
	htmlOutput += myHeader
	htmlOutput += "</font></H1>"
	return
}

func addSadDog () (htmlOutput string) {

	htmlOutput += "<H2><font color='red'>"
	htmlOutput += "Your overall mood is not great. We hope it will get better."
	htmlOutput += "</font>"
	htmlOutput += "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>"
	htmlOutput += "</H2>"
	return
}

func addHappyDog () (htmlOutput string) {

	htmlOutput += "<H2><font color='green'>"
	htmlOutput += "Your overall mood is happy. Keep it that way!"
	htmlOutput += "</font>"
	htmlOutput += "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>"
	htmlOutput += "</H2>"
	return
}

func addDataTitle (title string) (htmlOutput string) {

	htmlOutput += "<BR><BR>"
	htmlOutput += "<font color='purple'>"
	htmlOutput += title
	htmlOutput += "</font><BR>"
	return
}

func addDataContent (content string) (htmlOutput string) {

	htmlOutput += "<font color='gray'>" 
	htmlOutput += content
	htmlOutput += "</font>"
	return
}

func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	http.ListenAndServe(*addr, nil)
}
