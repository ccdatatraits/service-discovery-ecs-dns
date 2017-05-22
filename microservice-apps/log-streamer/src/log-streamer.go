package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "io/ioutil"
  "log"
  "net/http"
  "os/exec"
  "strconv"
  "strings"
)

// Spec represents the Log response object.
type LogResponse struct {
  Value       string `json:"value"`
  ContainerId string `json:"containerid"`
  InstanceId  string `json:"instanceid"`
}

// Spec represents the Healthcheck response object.
type Healthcheck struct {
  Status string `json:"status"`
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc(/*"/tracker.gif?{request}"*/"/{full_request}", LogHandler).Methods("GET")
  http.HandleFunc("/health", HealthHandler)
  http.Handle("/", r)
  log.Fatal(http.ListenAndServe(":8081", nil))
}

// Handler to process the healthcheck
func HealthHandler(res http.ResponseWriter, req *http.Request) {
  hc := Healthcheck{Status: "OK"}
  if err := json.NewEncoder(res).Encode(hc); err != nil {
    log.Panic(err)
  }
}

// Handler to process the GET HTTP request. Returns the time formatted
func LogHandler(res http.ResponseWriter, req *http.Request) {
  requestURI := req.RequestURI
  log.Printf("RequestURI is %s\n", requestURI)
  /*vars := mux.Vars(req)
  request := vars["request"]
  log.Printf("request is %s\n", request)*/

  result := strconv.QuoteToASCII(requestURI)

  logReponse := LogResponse{Value: result, ContainerId: GetContainerId(), InstanceId: GetInstanceId()}
  if err := json.NewEncoder(res).Encode(logReponse); err != nil {
    log.Panic(err)
  }
}

// Call the remote web service and return the result as a string
func GetHttpResponse(url string) string {
  response, err := http.Get(url)
  defer response.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Fatal(err)
  }
  return string(contents)

}

// Get the ContainerID if exists
func GetContainerId() string {
  //return ""
  cmd := "cat /proc/self/cgroup | grep \"docker\" | sed s/\\\\//\\\\n/g | tail -1"
  out, err := exec.Command("bash", "-c", cmd).Output()
  if err != nil {
    log.Printf("Container Id err is %s\n", err)
    return ""
  }
  log.Printf("The container id is %s\n", out)
  return strings.TrimSpace(string(out))
}

// Get the Instance ID if exists
func GetInstanceId() string {
  //return ""
  cmd := "curl"
  cmdArgs := []string{"-s", "http://169.254.169.254/latest/meta-data/instance-id"}
  out, err := exec.Command(cmd, cmdArgs...).Output()
  if err != nil {
    log.Printf("Instance Id err is %s\n", err)
    return ""
  }
  log.Printf("The instance id is %s\n", out)
  return string(out)
}
