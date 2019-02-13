package main

import (
  "html/template"
  "net/http"
  "strings"
  "strconv"
	log "github.com/sirupsen/logrus"
)

var (
	_homepageTemplate     *template.Template
  err                   error
)

func main() {
  //Configure Logging
  Formatter := new(log.TextFormatter)
  Formatter.TimestampFormat = "02-01-2006 15:04:05"
  Formatter.FullTimestamp = true
  log.SetFormatter(Formatter)
  log.SetLevel(log.TraceLevel)

  log.Info("------------Server Started------------")
  err = loadPageTemplates()
  if err != nil {
    log.Error(err)
  }
  http.HandleFunc("/", homePage)
  http.HandleFunc("/solidcolor", solidColor)
  err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Error("Http Error:", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
  log.Info("Homepage Hit")
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  data := make(map[string]interface{})
  data["Title"] = "Title"
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
}

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

func hex2int(hexStr string) int {
 // remove 0x suffix if found in the input string
 cleaned := strings.Replace(hexStr, "0x", "", -1)

 // base 16 for hexadecimal
 result, _ := strconv.ParseUint(cleaned, 16, 64)
 return int(result)
}

//Helpers
func loadPageTemplates() error {
	var err error
	_homepageTemplate, err = template.ParseFiles("web/index.html")
	if err != nil {
		return err
	}
	return err
}
