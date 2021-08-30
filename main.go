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

// GithubHook UserHook

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

func main() {
	loadConfig()

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
			go execShell(config)
			writeDone(w)
			return
		}
	}

	log.Printf("未注册的操作 --> App:%s Type:%s", appName, appType)
}

func execShell(config AppConfig) {
	log.Printf("开始执行请求 -->  App:%s Type:%s", config.AppName, config.Type)

	var fullCommand = fmt.Sprintf("./command/%s %s", config.Template, config.AppName)

	var cmd = exec.Command("bash", "-c", fullCommand)

	output, _ := cmd.Output()

	fileLog := OpenLog(config.AppName)
	fileLog.LogOnce(fmt.Sprintf("执行指令: %s\n执行过程中的输出:\n%s", fullCommand, string(output)))

}

func writeDone(w http.ResponseWriter) {
	_, err := fmt.Fprint(w, "Accepted.")
	if err != nil {
		return
	}
}
