package task

import (
	"fmt"
	"os"
	"os/exec"
)

type Task struct {
	Name       string `json:"appName"`
	Type       string `json:"type"`
	Template   string `json:"template"`
	Background bool   `json:"background"`
}

func Daemon(c <-chan Task, m Manager) {
	for task := range c {

		task.ExecShell(m)
		m.FinishTask()
	}
}

func (task Task) ExecShell(m Manager) string {
	m.SetTask(task.Name)
	defer m.FinishTask()

	msg := fmt.Sprintf("执行器状态:%s\n\n", m.ToString())

	var fullCommand = fmt.Sprintf("./command/%s %s", task.Template, task.Name)
	var cmd = exec.Command("bash", "-c", fullCommand)
	output, _ := cmd.Output()

	msg += string(output)

	task.writeLog(msg)
	return msg
}

func (task Task) writeLog(content string) {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", os.ModePerm)
		check(err)
	}

	f, err := os.Create("log/" + task.Name + ".log")
	defer f.Close()

	check(err)

	_, err = f.WriteString(content)
	check(err)
}
