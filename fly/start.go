package fly

import (
	"time"
)

func Start() {
	println("fly start")
	Init()
	for {
		for _, t := range tm.TM {
			// go func(t Task) {
			// Check if the task can run before executing
			ok := t.TimeIsUp()
			if !ok {
				return
			}
			// Run the task and update its status
			if err := t.Run(); err != nil {
				t.Runner.Status = StatusError
				t.Runner.Msg = err.Error()
			} else {
				t.Runner.Status = StatusSuccess
			}
			// Update the task's last runtime and next runtime after execution
			// Update the task in the database
			t.Update()
			// }(task)
		}
		// Sleep for a short duration before checking the tasks again
		time.Sleep(1 * time.Second)
	}
}
