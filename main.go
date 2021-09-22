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
)

type Config struct {
	Token  string
	Config []AppConfig
}

type AppConfig struct {
	AppName  string
	Type     string
	Template string
}

var configs Config

var c = make(chan AppConfig, 10)

var currentTask = "空闲等待中..."

func main() {
	loadConfig()

	// handler是异步执行的
	http.HandleFunc("/", handleWebHook)
	go doShellCommand()

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

func handleWebHook(w http.ResponseWriter, request *http.Request) {
	var URI = request.URL.RequestURI()
	log.Println("接受请求: " + URI)

	value := strings.Split(URI, "/")
	// "/Token/Type/AppName" -> ["", "Token", "Type", "AppName"]
	if len(value) != 4 {
		log.Print("参数数量不正确, 忽略请求")
		return
	}

	var token = value[1]
	if token != configs.Token {
		log.Print("Token错误, 忽略请求")
		return
	}

	var appType = value[2]
	var appName = value[3]

	for _, config := range configs.Config {
		if config.AppName == appName && config.Type == appType {
			if appName == "System" {
				writeMessage(w, fmt.Sprintf("执行器状态:%s\n\n", currentTask))
				writeMessage(w, execShell(config))
			} else {
				c <- config
				writeDone(w)
			}
			return
		}
	}

	log.Printf("未注册的操作 --> App:%s Type:%s", appName, appType)
}

func doShellCommand() {
	for {
		config := <-c
		currentTask = fmt.Sprintf("执行任务中(%s)", config.AppName)
		execShell(config)
		currentTask = "空闲等待中..."
	}
}

func execShell(config AppConfig) string {
	log.Printf("开始执行请求 -->  App:%s Type:%s", config.AppName, config.Type)

	var fullCommand = fmt.Sprintf("./command/%s %s", config.Template, config.AppName)

	var cmd = exec.Command("bash", "-c", fullCommand)

	output, _ := cmd.Output()
	msg := string(output)
	fileLog := OpenLog(config.AppName)
	fileLog.LogOnce(fmt.Sprintf("执行指令: %s\n执行过程中的输出:\n%s", fullCommand, msg))
	return msg
}

func writeDone(w http.ResponseWriter) {
	_, err := fmt.Fprint(w, "Accepted.")
	if err != nil {
		return
	}
}

func writeMessage(w http.ResponseWriter, msg string) {
	_, err := fmt.Fprint(w, msg)
	if err != nil {
		return
	}
}
