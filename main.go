package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	TYPE_GITHUB string = "GithubHook"
	TYPE_USER   string = "UserHook"

	STATUS_DONE string = "Done"
	STATUS_RUN  string = "Running"
)

type AppConfig struct {
	AppName string
	Type    string
	WorkDir string
	Cmd     []string
}

type AppStatus struct {
	AppName string
	Status  string
	Time    string
}

var configs []AppConfig
var status = make(map[string]*AppStatus)

func main() {
	loadConfig()
	initStatus()
	// handler是异步执行的
	http.HandleFunc("/", handleWebHook)

	err := http.ListenAndServe(":3080", nil)
	if err != nil {
		log.Panic(err)
	}
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panic(err)
		}
	}(file)

	content, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(content, &configs)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Load Config: %v\n", configs)
}

func initStatus() {
	for _, config := range configs {
		status[config.AppName] = &AppStatus{AppName: config.AppName, Status: STATUS_DONE, Time: nowString()}
	}
}

func handleWebHook(writer http.ResponseWriter, request *http.Request) {
	var URI = request.URL.RequestURI()
	log.Println("Receive Request: " + URI)
	if URI == "/" {
		writeIndexFile(writer, request)
		return
	}

	value := strings.Split(URI, "/")
	// "/A/B" -> ["", "A", "B"]
	if len(value) != 3 {
		writeError(writer, request, "参数数量不正确")
		return
	}

	var appType = value[1]
	var appName = value[2]
	for _, config := range configs {
		if config.AppName == appName && config.Type == appType {
			appStatus := status[config.AppName]
			setRunning(appStatus)
			execShell(config.WorkDir, config.Cmd)
			setDone(appStatus)
		}
	}

	for i := 0; i < len(configs); i++ {
		if configs[i].AppName == appName && configs[i].Type == appType {


			return
		}
	}

	msg := fmt.Sprintf("Undefined Request: appName= %s appType=%s", appName, appType)
	writeError(writer, request, msg)
}

func writeIndexFile(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%-15s\t%-10s\t%-15s\n", "AppName", "Status", "Time")

	for _, s := range status {
		_, _ = fmt.Fprintf(w, "%-15s\t%-10s\t%-15s\n", s.AppName, s.Status, s.Time)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, msg string) {
	_, _ = fmt.Fprintf(w, msg)
}

func execShell(workDir string, cmd []string) {
	var fullCommand = "cd " + workDir + ";"
	for i := 0; i < len(cmd); i++ {
		fullCommand = fullCommand + cmd[i] + ";"
	}
	log.Printf("Do Command: %s", fullCommand)

	_ = exec.Command("bash", "-c", fullCommand).Run()
}

func nowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func setRunning(appStatus *AppStatus)  {
	appStatus.Status = STATUS_RUN
	appStatus.Time = nowString()
}

func setDone(appStatus *AppStatus)  {
	appStatus.Status = STATUS_DONE
	appStatus.Time = nowString()
}