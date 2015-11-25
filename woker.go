package nodenet

import (
	"fmt"
	"reflect"
)

type Worker struct {
	Message reflect.Type   // 可以处理的消息类型
	Handler MessageHandler // 消息处理函数
}

var workers map[string]*Worker = make(map[string]*Worker)

func RegisterWorker(name string, message interface{}, worker MessageHandler) {
	if w, ok := workers[name]; ok {
		panic(fmt.Errorf("registering duplicate worker for: %s. %v => %v", name, w, worker))
	}

	if reflect.TypeOf(message).Kind() != reflect.Struct {
		panic("Only struct.")
	}

	workers[name] = &Worker{Message: reflect.TypeOf(message), Handler: worker}
}

// 节点名应该配置成: 节点业务类型(组名)-节点名
func GetWorkerByName(name string) *Worker {
	return workers[name]
}
