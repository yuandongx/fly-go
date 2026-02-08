package fly

import (
	"time"
)

func Start() {
	tm, err := Init()
	if err != nil {
		println("Failed to initialize application:", err.Error())
		return
	}
	for {
		tm.RunAllTask()
		// Sleep for a short duration before checking the tasks again
		time.Sleep(1 * time.Second)
	}
}
