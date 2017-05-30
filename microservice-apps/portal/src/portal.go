package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Spec represents the ErrorResponse response object.
type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{full_request}", LogHandler).Methods("GET")
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir(GetHtmlFileDir())))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler to process the GET HTTP request. Calls the Log service
func LogHandler(res http.ResponseWriter, req *http.Request) {
  requestURI := req.RequestURI
  log.Printf("RequestURI is %s\n", requestURI)
  result := strconv.QuoteToASCII(requestURI)

  var log_endpoint string
  var err error
  if log_endpoint, err = GetLogEndpoint(); err != nil {
    log.Fatal(err)
  }
  url := "http://" + log_endpoint + "/" + result
  log.Printf("URL is %v\n", url)

  req, err = http.NewRequest("GET", url, nil)

  cli := &http.Client{}
  response, err := cli.Do(req)
  defer response.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

	if response.StatusCode != 200 {
		log.Printf("Got Error response: %s\n", response.Status)
		errorMsg := ErrorResponse{Error: response.Status}
		if err := json.NewEncoder(res).Encode(errorMsg); err != nil {
			log.Panic(err)
		}
	} else {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Got Successful response: %s\n", string(contents))
		res.Write(contents)
	}
}

// get the log service endpoint from DNS
func GetLogEndpoint() (string, error) {
	var addrs []*net.SRV
	var err error
	if _, addrs, err = net.LookupSRV("", "", "log-streamer.servicediscovery.internal"); err != nil {
		return "", err
	}
	for _, addr := range addrs {
		return strings.TrimRight(addr.Target, ".") + ":" + strconv.Itoa(int(addr.Port)), nil
	}
	return "", errors.New("No record found")
}

/*// get the HTML file
func GetHtmlFileDir() string {
	html_file_dir := os.Getenv("HTML_FILE_DIR")
	if len(html_file_dir) > 0 {
		return html_file_dir
	}
	log.Println("env variables HTML_FILE_DIR not found. Returning default values: ./public/")
	return "./public/"
}*/
