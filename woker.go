package nodenet

import (
	"strings"
)

var workers map[string]ComponentHandler = make(map[string]ComponentHandler)

func SetWorker(nodeType string, handler ComponentHandler) {
	workers[nodeType] = handler
}

func GetWorkerByType(nodeType string) ComponentHandler {
	if worker, ok := workers[nodeType]; ok {
		return worker
	}

	return nil
}

func GetWorkerByName(name string) ComponentHandler {
	t := strings.Split(name, ".")
	if len(t) < 1 {
		return nil
	}

	return GetWorkerByType(t[0])
}
