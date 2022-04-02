package task

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token string
	Tasks []Task
}

func (config *Config) Load() {
	file, err := os.Open("config.json")
	defer file.Close()
	check(err)

	content, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(content, config)
	check(err)
}

func (config *Config) Match(req Task) (Task, error) {
	for _, task := range config.Tasks {
		if task.Name == req.Name && task.Type == req.Type {
			return task, nil
		}
	}

	return Task{}, errors.New("未注册的操作")
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
