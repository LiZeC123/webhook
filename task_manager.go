package main

import (
	"fmt"
	"time"
)

type TaskManager struct {
	currentTaskName string
	currentTaskTime int64
}

func (taskManager *TaskManager) Init() {
	taskManager.FinishTask()
}

func (taskManager *TaskManager) SetTask(taskName string) {
	taskManager.currentTaskName = taskName
	taskManager.currentTaskTime = time.Now().Unix()
}

func (taskManager *TaskManager) FinishTask() {
	taskManager.currentTaskName = ""
	taskManager.currentTaskTime = time.Now().Unix()
}

func (taskManager TaskManager) GetFormatString() string {
	var elapsedTime = time.Now().Unix() - taskManager.currentTaskTime
	if taskManager.currentTaskName == "" {
		return fmt.Sprintf("空闲等待中(已等待%d秒)", elapsedTime)
	} else {
		return fmt.Sprintf("任务[%s]执行中(已耗时%d秒)", taskManager.currentTaskName, elapsedTime)
	}
}
