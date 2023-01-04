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

var SENSORS_ACTIVATION_BATCH int = 20

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
	fmt.Fprintf(w,addMoodResults(),	pureHappy,pureSad, pureAngry,
									existingHappy,existingSad,existingAngry)
									
	//render happy/sad
	sniffThreshold, err := strconv.ParseFloat(os.Getenv("SNIFF_THRESHOLD"),64)
	if err != nil { fmt.Fprintf(w,"!!Error in converting sniffing threhold to float64")}
	
	if pureHappy > sniffThreshold {
		fmt.Fprintf(w, addDog("happy"))
	} else {
		fmt.Fprintf(w, addDog("sad"))
	}
		
	//render API section
	fmt.Fprintf(w,addDataTitle("Sniffing threshold"))
	fmt.Fprintf(w,addDataContent("%.2f%% pure happiness"),sniffThreshold)
	fmt.Fprintf(w,addDataTitle("/activate API"))
	fmt.Fprintf(w,addDataContent("All sensors activated successfully"))
	fmt.Fprintf(w,addDataTitle("/measure API"))
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

	htmlOutput += "<p align='center'>"
	htmlOutput += "<table align='center' border='0'>"
	
	//pure mood row
	htmlOutput += "<tr>"
	htmlOutput += "<td style='font-size:30px;color:DarkGreen'>%.2f%% Happy</td>"
	htmlOutput += "<td style='font-size:30px;color:DarkRed'>%.2f%% Sad</td>"
	htmlOutput += "<td style='font-size:30px;color:DarkOrange'>%.2f%% Angry</td>"
	htmlOutput += "</tr>"

	//pre-existing row	
	htmlOutput += "<tr style='font-size:15px;color:gray'>"
	htmlOutput += "<td>(%.2f%% w/ pre-existing)</td>"
	htmlOutput += "<td>(%.2f%% w/ pre-existing)</td>"
	htmlOutput += "<td>(%.2f%% w/ pre-existing)</td>"
	htmlOutput += "</tr>"
	
	htmlOutput += "</table></p>"
	return
}

func addHeader (myHeader string) (htmlOutput string) {

    htmlOutput += "<body><p style='font-size:50px;color:navy' align='center'>"
	htmlOutput += myHeader
	htmlOutput += "</p>"
	return
}

func addDog (imgPrefix string) (htmlOutput string) {

	htmlOutput += "<p style='font-size:20px;color:purple' align='center'>"
	htmlOutput += "<img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/" + imgPrefix + "-dog.jpg' alt=''>"
	htmlOutput += "</p>"
	return
}

func addAPICallsTable () (htmlOutput string) {

	htmlOutput += "<p align='left'>"
	htmlOutput += "<table style='font-size:15px;color:gray' border='1'>"
	
	htmlOutput += "<tr>"
	htmlOutput += "<th>Sensor</th>" + "<th>Role</th>" + "<th>Current Mood</th>"+ "<th>Pre-Existing</th>"
	htmlOutput += "</tr>"

	for _, sensor := range AllSensorsData.Sensors {
  		htmlOutput += "<tr>"
		htmlOutput += "<td>" + strconv.Itoa(sensor.Id) + "</td>"
		htmlOutput += "<td>" + sensor.Role + "&nbsp;</td>"
		htmlOutput += "<td>" + sensor.Mood + "&nbsp;</td>"
		htmlOutput += "<td>" + sensor.Legacy + "</td>"
		htmlOutput += "</tr>"
	}

	htmlOutput += "</table></p></body>"
	return
}

func addDataTitle (title string) (htmlOutput string) {

	htmlOutput += "<p style='font-size:15px;color:purple' align='left'>"
	htmlOutput += title
	htmlOutput += "</p>"
	return
}

func addDataContent (content string) (htmlOutput string) {

	htmlOutput += "<p style='font-size:15px;color:gray' align='left'>"
	htmlOutput += content
	htmlOutput += "</p><BR>"
	return
}

func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	http.ListenAndServe(*addr, nil)
}
