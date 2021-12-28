package main


import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"sort"
)

var ALWAYS_HAPPY = true


func handler(w http.ResponseWriter, r *http.Request) {
	
	log.Println(r.RemoteAddr, r.Method, r.URL.String())
	
        fmt.Fprintf(w, "<H1><font color='navy'>Welcome to the DevX Mood Analyzer </font></H1><H2>")

	if ALWAYS_HAPPY == false {
	
		fmt.Fprintf(w, "<font color='maroon'>Your mood sensors' current data:</font>")
		fmt.Fprintf(w, "</font><BR><BR>")

		response, err := http.Get("http://devx-mood-backend.dekt-apps.serving.dekt.io/sensors-data")
	
		if err != nil {
			fmt.Fprintf(w,"ERROR! in calling API")
		} else {
			defer response.Body.Close()
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Fprintf(w,"ERROR! in reading body")
			} else {
				fmt.Fprintf(w, "<font color='gray'>")
				jsonStr := string(responseData)
				ok, resTable := JSON2HtmlTable(jsonStr, []string{"title2", "title1"}, []string{"title1"})
				if ok {
					fmt.Println(resTable)
				} else {
					fmt.Fprintf(w,jsonStr)
				}
				fmt.Fprintf(w, "</font><BR><BR>")
				fmt.Fprintf(w, "<font color='red'>")
				fmt.Fprintf(w,"Your overall mood is not great. We hope it will get better.")
				fmt.Fprintf(w, "</font>")
				fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/sad-dog.jpg' alt=''>")
			}
		}
	} else {
		fmt.Fprintf(w, "<font color='green'>")
		fmt.Fprintf(w,"Your mood is always happy. Good for you!")
		fmt.Fprintf(w, "</font>")
		fmt.Fprintf(w, "<BR><BR><img src='https://raw.githubusercontent.com/dektlong/devx-mood/main/happy-dog.jpg' alt=''>")
	}
	fmt.Fprintf(w, "</H2>")
}

// JSON2HtmlTable convert json string to html table string
func JSON2HtmlTable(jsonStr string, customTitles []string, rowSpanTitles []string) (bool, string) {
	htmlTable := ""
	jsonArray := []map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonStr), &jsonArray)
	if err != nil || 0 == len(jsonArray) {
		fmt.Println("invalid json string")
		return false, htmlTable
	}

	titles := customTitles
	if nil == customTitles || 0 == len(customTitles) { // if custom titles are not provided, use json keys as titles
		titles = getKeys(jsonArray[0])
	}

	if nil != rowSpanTitles && 0 != len(rowSpanTitles) { // if sort keys are provided, sort json array
		for tid, title := range rowSpanTitles {
			swapped := true
			for swapped {
				swapped = false
				for i := 0; i < len(jsonArray)-1; i++ {
					va, oka := jsonArray[i][title].(string)
					vb, okb := jsonArray[i+1][title].(string)
					if !oka || !okb {
						swapped = false
						break
					}
					if strings.Compare(va, vb) > 0 {
						if tid != 0 {
							va, _ := jsonArray[i][rowSpanTitles[tid-1]].(string)
							vb, _ := jsonArray[i+1][rowSpanTitles[tid-1]].(string)
							if va != vb {
								continue
							}
						}
						tmp := jsonArray[i]
						jsonArray[i] = jsonArray[i+1]
						jsonArray[i+1] = tmp
						swapped = true
					}
				}
			}
		}
	}
	// convert table headers
	if 0 == len(titles) {
		fmt.Println("json is not supported")
	}
	tmp := []string{}
	for _, title := range titles {
		tmp = append(tmp, fmt.Sprintf("<th>%s</th>", title))
	}
	thCon := strings.Join(tmp, "")

	// convert table cells
	segs := map[string][]int{}
	initSeg := []int{0, len(jsonArray)}
	for i, key := range rowSpanTitles {
		seg := initSeg
		for j:=1; j<len(jsonArray); j++ {
			if jsonArray[j][key] != jsonArray[j-1][key] {
				inSlice := false
				for _, k := range seg {
					if k == j {
						inSlice = true
					}
				}
				if !inSlice {
					seg = append(seg, j)
				}
			}
		}
		sort.Ints(seg)
		segs[rowSpanTitles[i]] = seg
		if i < len(rowSpanTitles) - 1 {
			segs[rowSpanTitles[i+1]] = segs[key]
			initSeg = segs[key]
		}
	}
	rows := []string{}
	for i, jsonObj := range jsonArray {
		tmp = []string{}
		for _, key := range titles {
			seg := segs[key]
			if seg != nil && len(seg) != 0 {
				if 0 == i {
					cell := fmt.Sprintf(`<td rowspan="%d">%v</td>`, seg[1], jsonObj[key])
					tmp = append(tmp, cell)
				} else {
					for n, j := range seg {
						if j == i {
							rowspan := 1
							if n < len(seg)-1 {
								rowspan = seg[n+1] - seg[n]
							}
							cell := fmt.Sprintf(`<td rowspan="%d">%v</td>`, rowspan, jsonObj[key])
							tmp = append(tmp, cell)
						}
					}
				}
			} else {
				cell := fmt.Sprintf("<td>%v</td>", jsonObj[key])
				tmp = append(tmp, cell)
			}
			//cell := fmt.Sprintf("<td>%v</td>", jsonObj[key])
			//tmp = append(tmp, cell)
		}
		tdCon := strings.Join(tmp, "")
		row := fmt.Sprintf("<tr>%s</tr>", tdCon)
		rows = append(rows, row)
	}
	trCon := strings.Join(rows, "")

	htmlTable = fmt.Sprintf(`<table border="1" cellpadding="1" cellspacing="1">%s%s</table>`,
		fmt.Sprintf("<thead>%s</thead>", thCon), fmt.Sprintf("<tbody>%s</tbody>", trCon))
	return true, htmlTable
}

func getKeys(jsonObj map[string]interface{}) []string {
	keys := make([]string, 0, len(jsonObj))
	for k := range jsonObj {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	
	http.HandleFunc("/", handler)

	var addr = flag.String("addr", ":8080", "addr to bind to")
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
