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

type Sensor struct {
	Id int `json:"id"`
	Role string `json:"role"`
	Mood string `json:"mood"`
	Legacy string `json:"legacy"`
}

type AllSensors struct {
	Sensors []*Sensor
}

func handler(w http.ResponseWriter, r *http.Request) {
	
	//conrtol the mood sniffing algorithm intensity
	beHappy := false

	fmt.Fprintf(w, addHeader("DevX Mood Analyzer"))
    
	if !beHappy { 
		fmt.Fprintf(w, sadMood())
	} else { 
		fmt.Fprintf(w, happyMood())
	}	
	
	//sensors activation
	fmt.Fprintf(w,addDataTitle("/activate"))
	fmt.Fprintf(w,addDataContent(processSensorActivation(10)))
	
	//sensors measurements
	fmt.Fprintf(w,addDataTitle("/measure"))
	fmt.Fprintf(w,addDataContent(processSensorsMeasurement()))
}

func processSensorActivation(numSensors int) (htmlOutput string) {

	tlsConfig := &http.Transport{
	 	TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	 }


	tlsClient := &http.Client{Transport: tlsConfig}
	for i := 0; i < numSensors; i++ {
		response, err := tlsClient.Get(os.Getenv("SENSORS_ACTIVATE_API"))	
		if err != nil { 
			htmlOutput = err
			return 
		} 	 	
		defer response.Body.Close()
	}
	
	htmlOutput += "All sensors activated successfully"
	return
}

func processSensorsMeasurement() (htmlOutput string) {
	
	tlsConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}


	tlsClient := &http.Client{Transport: tlsConfig}

	response, err := tlsClient.Get(os.Getenv("SENSORS_MEASURE_API"))	 

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
	htmlOutput += "<th>Sensor</th>" + "<th>Role</th>" + "<th>Current Mood</th>"+ "<th>Pre-Existing</th>"
	htmlOutput += "</tr>"

	for _, sensor := range allSensors.Sensors {
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

func sadMood () (htmlOutput string) {

	htmlOutput += "<H2><font color='red'>"
	htmlOutput += "Your overall mood is not great. We hope it will get better."
	htmlOutput += "</font>"
	htmlOutput += "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>"
	htmlOutput += "</H2>"
	htmlOutput += addDataTitle("mood sniffing algorithm")
	htmlOutput += addDataContent("Agressive")
	return
}

func happyMood () (htmlOutput string) {

	htmlOutput += "<H2><font color='green'>"
	htmlOutput += "Your overall mood is happy. Keep it that way!"
	htmlOutput += "</font>"
	htmlOutput += "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>"
	htmlOutput += "</H2>"
	htmlOutput += addDataTitle("mood sniffing algorithm")
	htmlOutput += addDataContent("Mild")
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