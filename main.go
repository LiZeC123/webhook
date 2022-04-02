package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/LiZeC123/webhook/task"
)

var config task.Config
var c = make(chan task.Task, 10)
var manager task.Manager

func main() {
	config.Load()
	manager.Init()
	go task.Daemon(c, manager)

	fmt.Println(config)

	// handler是异步执行的
	http.HandleFunc("/", handleWebHook)

	err := http.ListenAndServe(":3080", nil)
	if err != nil {
		log.Panic(err)
	}
}

func handleWebHook(w http.ResponseWriter, request *http.Request) {
	var URI = request.URL.RequestURI()
	log.Println("接受请求: " + URI)

	req, err := parseRequest(URI)
	if err != nil {
		log.Printf("%s, 忽略此请求\n", err.Error())
		return
	}

	m, err := config.Match(req)
	if err != nil {
		log.Printf("%s --> App:%s Type:%s", err.Error(), req.Name, req.Type)
	} else {
		log.Printf("开始执行请求 -->  App:%s Type:%s", req.Name, req.Type)
	}

	var msg string
	if m.Background {
		// 后台任务发送信息异步执行
		c <- m
		msg = "Accepted."
	} else {
		// 前台任务直接在当前线程执行
		msg = m.ExecShell(manager)
	}

	writeMessage(w, msg)
}

func parseRequest(URI string) (task.Task, error) {
	// "/Token/Type/AppName" -> ["", "Token", "Type", "AppName"]
	value := strings.Split(URI, "/")
	if len(value) != 4 {
		return task.Task{}, errors.New("参数数量不正确")
	}

	var token = value[1]
	if token != config.Token {
		return task.Task{}, errors.New("Token不匹配")
	}

	return task.Task{Name: value[3], Type: value[2]}, nil
}

func writeMessage(w http.ResponseWriter, msg string) {
	_, err := fmt.Fprint(w, msg)
	if err != nil {
		return
	}
}
