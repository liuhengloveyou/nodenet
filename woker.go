package nodenet

import (
	"strings"
)

var workers map[string]MessageHandler = make(map[string]MessageHandler)

func SetWorker(nodeType string, handler MessageHandler) {
	workers[nodeType] = handler
}

func GetWorkerByType(nodeType string) MessageHandler {
	if worker, ok := workers[nodeType]; ok {
		return worker
	}

	return nil
}

// 节点名应该配置成: 节点业务类型(组名)-节点名
func GetWorkerByName(name string) MessageHandler {
	t := strings.Split(name, "-")
	if len(t) < 1 {
		return nil
	}

	return GetWorkerByType(t[0])
}
