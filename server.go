package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"fmt"
)

var (
	_homepageTemplate *template.Template
  knt               int
	err               error
	client 						mqtt.Client
)

func main() {
	//Configure Logging
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetLevel(log.TraceLevel)

	//Start Server
	log.Info("------------Server Started------------")
	err = loadPageTemplates()
	if err != nil {
		log.Error(err)
	}
	http.HandleFunc("/", homePage)
	http.HandleFunc("/solidcolor", solidColor)
	go func() {
		err = http.ListenAndServe(":80", nil)
		if err != nil {
			log.Error("Http Error:", err)
		}
	}()

	//Run mqtt
	opts := mqtt.NewClientOptions()
	opts.SetClientID("go-simple")
	opts.AddBroker("localhost:1883")
	//create and start a client using the above ClientOptions
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	//subscribe
	if token := client.Subscribe("rgb", 0, mqttMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	select {

	}
}

func mqttMessageHandler(client mqtt.Client, msg mqtt.Message) {
  fmt.Printf("TOPIC: %s\n", msg.Topic())
  fmt.Printf("MSG: %s\n", msg.Payload())
}

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Info("Homepage Hit")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := make(map[string]interface{})
	data["Title"] = "LED Color Select"
	err := _homepageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func solidColor(w http.ResponseWriter, r *http.Request) {
	log.Info("solidColor Hit")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	m := r.URL.Query()
	rgb := HexToRGB(m["selectedcolor"])
	log.Info("Selected Color:", rgb)
	sRGB := intArrayToString(rgb)
	log.Trace("rgb string:", sRGB)
	token := client.Publish("rgb",0,false,sRGB)
	token.Wait()

}

//This function converts the returned color of form #AABBCC to rgb(1,2,3)
func HexToRGB(hex []string) []int {
	var r, g, b string
	var red, green, blue int
	var result []int
	log.Info("HexToRGB entered")
	log.Trace("Hex:", hex)
	log.Trace("Hex[0]:", hex[0])
	hex[0] = hex[0][1:]
	log.Info("Hex[0][1:]", hex[0])
	r = hex[0][:2]
	g = hex[0][2:4]
	b = hex[0][4:]
	log.Trace("hex r: ", r)
	log.Trace("hex g: ", g)
	log.Trace("hex b: ", b)
	red = hex2int(r)
	green = hex2int(g)
	blue = hex2int(b)
	log.Trace("red: ", red)
	log.Trace("green: ", green)
	log.Trace("blue: ", blue)
	result = append(result, red)
	result = append(result, green)
	result = append(result, blue)
	return result
}

//This function is used by HexToRGB()
func hex2int(hexStr string) int {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return int(result)
}

//This function converts an int array into a string representing rgb values to me used in the mqtt message
func intArrayToString(input []int) string {
	result := "("
	for i, v := range input {
		if i <= 1 {
			result = result + strconv.Itoa(v) + ","
		} else {
			result = result + strconv.Itoa(v) + ")"
		}
}
	return result
}



//This function is used to load the HTML template
func loadPageTemplates() error {
	var err error
	_homepageTemplate, err = template.ParseFiles("web/index.html")
	if err != nil {
		return err
	}
	return err
}
