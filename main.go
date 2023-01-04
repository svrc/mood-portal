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

var SENSORS_ACTIVATION_BATCH int = 10

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

	//process APIs calls and analytics
	if processSensorActivation() != "success" {
		return
	}
	
	if processSensorsMeasurement() != "success" {
		return
	}
	pureHappy,existingHappy,pureSad,existingSad,pureAngry,existingAngry := moodAnalysis()

	//render results section
	fmt.Fprintf(w,addMoodResults(),	pureHappy,existingHappy,
									pureSad,existingSad,
									pureAngry,existingAngry)

	//render happy/sad
	sniffThreshold, err := strconv.ParseFloat(os.Getenv("SNIFF_THRESHOLD"),64)
	if err != nil { fmt.Fprintf(w,"!!Error in converting sniffing threhold to float64")}
	
	if pureHappy > sniffThreshold {
		fmt.Fprintf(w, addDog(true),sniffThreshold)
	} else {
		fmt.Fprintf(w, addDog(false),sniffThreshold)
	}
		
	//render API section
	fmt.Fprintf(w,addDataTitle("/activate"))
	fmt.Fprintf(w,addDataContent("All sensors activated successfully"))
	fmt.Fprintf(w,addDataTitle("/measure"))
	fmt.Fprintf(w,addDataContent(addAPICallsTable()))
}

func processSensorActivation() (status string) {

	tlsConfig := &http.Transport{
	 	TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	 }


	tlsClient := &http.Client{Transport: tlsConfig}
	for i := 0; i < SENSORS_ACTIVATION_BATCH ; i++ {
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

func moodAnalysis () (	float64, float64, //pure happy, pre-existing happy
						float64, float64, //pure sad, pre-existing sad
						float64, float64) { //pure angry, pre-existing angry
	
	var pureHappy,existingHappy,pureSad,existingSad,pureAngry,existingAngry float64 = 0.0,0.0,0.0,0.0,0.0,0.0
	var totalMeasurements float64 = float64(len(AllSensorsData.Sensors))
	
	for _, sensor := range AllSensorsData.Sensors {
		if sensor.Mood == "happy" {
			if sensor.Legacy == "none" {
				pureHappy++
			} else {
				existingHappy++
			}
		} else if sensor.Mood == "sad" {
			if sensor.Legacy == "none" {
				pureSad++
			} else {
				existingSad++
			}
		} else if sensor.Mood == "angry" {
			if sensor.Legacy == "none" {
				pureAngry++
			} else {
				existingAngry++
			}
		} else { 
			//error
		}
	}
	
	return	(pureHappy/totalMeasurements)*100,(existingHappy/totalMeasurements)*100,
			(pureSad/totalMeasurements)*100,(existingSad/totalMeasurements)*100,
			(pureAngry/totalMeasurements)*100,(existingAngry/totalMeasurements)*100
	
}

func addMoodResults () (htmlOutput string) {

	htmlOutput += "<H2><table border='0'>"
	
	htmlOutput += "<tr style='color:DarkGreen' align='left'>"
	htmlOutput += "<td>>Happy:</td>"
	htmlOutput += "<td>%.2f%%</td>"
	htmlOutput += "<td><small>(%.2f%% w/ pre-existing)</small></td>"
	
	htmlOutput += "<tr style='color:DarkRed' align='left'>"
	htmlOutput += "<td>>Sad:</td>"
	htmlOutput += "<td>%.2f%%</td>"
	htmlOutput += "<td><small>(%.2f%% w/ pre-existing)</small></td>"

	htmlOutput += "<tr style='color:DarkOrange' align='left'>"
	htmlOutput += "<td>>Angery:</td>"
	htmlOutput += "<td>%.2f%%</td>"
	htmlOutput += "<td><small>(%.2f%% w/ pre-existing) </small></td>"
	
	htmlOutput += "<tr></table></H2>"
	return
}

func addHeader (myHeader string) (htmlOutput string) {

    htmlOutput += "<H1><font color='navy'>"
	htmlOutput += myHeader
	htmlOutput += "</font></H1>"
	return
}

func addDog (happy bool) (htmlOutput string) {

	htmlOutput += "<BR><BR>"

	if happy {
		htmlOutput += "<img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>"

	} else {
		htmlOutput += "<img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>"
	}
	
	htmlOutput += "&nbsp;&nbsp;&nbsp;"
	htmlOutput += "<font color='navy'>Sniffing threshold: %.2f%%</font>"
	return
}

func addAPICallsTable () (htmlOutput string) {

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
