package fly

import (
	"fly-go/fly/executor"
)

// 所有的runner需要放在这里
var runners = []executor.Runner{
	&executor.DemoRunner1{Alias: "DemoRunner1"},
	&executor.DemoRunner2{Alias: "DemoRunner2"},
}

// GetRunners 返回所有已注册的 Runner
func GetRunners() map[string]executor.Runner {
	// TODO: 从 TaskQueue 获取已注册的 runner
	rs := make(map[string]executor.Runner, 0)
	for _, runner := range runners {
		rs[runner.Name()] = runner
	}
	return rs
}
