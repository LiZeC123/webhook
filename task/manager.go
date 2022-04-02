package task

import (
	"fmt"
	"time"
)

type Manager struct {
	name string
	time int64
}

func (m *Manager) Init() {
	m.FinishTask()
}

func (m *Manager) SetTask(name string) {
	m.name = name
	m.time = time.Now().Unix()
}

func (m *Manager) FinishTask() {
	m.name = ""
	m.time = time.Now().Unix()
}

func (m *Manager) ToString() string {
	var elapsedTime = time.Now().Unix() - m.time
	if m.name == "" {
		return fmt.Sprintf("空闲等待中(已等待%d秒)", elapsedTime)
	} else {
		return fmt.Sprintf("任务[%s]执行中(已耗时%d秒)", m.name, elapsedTime)
	}
}
