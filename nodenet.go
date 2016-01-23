package nodenet

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	components map[string]*Component      = make(map[string]*Component)
	groups     map[string]*ComponentGroup = make(map[string]*ComponentGroup)
	graphs     map[string][]string        = make(map[string][]string)
)

var Config struct {
	Name   string              `json:"name"`
	Graphs map[string][]string `json:"graphs"`
	Groups []struct {
		Name     string   `json:"name"`
		Dispense string   `json:"dispense"`
		Members  []string `json:"members"`
	} `json:"groups"`
	Components []struct {
		Name   string      `json:"name"`
		InType string      `json:"intype"`
		InConf interface{} `json:"inconf"`
	} `json:"components"`
}

func BuildFromConfig(fileName string) error {
	r, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&Config); err != nil {
		return err
	}

	// 组件
	for i := 0; i < len(Config.Components); i++ {
		_, err := NewComponent(Config.Components[i].Name, Config.Components[i].InType, Config.Components[i].InConf)
		if err != nil {
			return err
		}
	}

	// 组
	for i := 0; i < len(Config.Groups); i++ {
		l := len(Config.Groups[i].Members)
		coms := make([]*Component, l)

		for j := 0; j < l; j++ {
			coms[j] = components[Config.Groups[i].Members[j]]
			if coms[j] == nil {
				return fmt.Errorf("No component name's [%s]", Config.Groups[i].Members[j])
			}
			coms[j].Group = Config.Groups[i].Name
		}

		groups[Config.Groups[i].Name] = NewGroup(Config.Groups[i].Name, Config.Groups[i].Dispense, coms)
	}

	// 图
	graphs = Config.Graphs
	for gk, gv := range graphs {
		for _, name := range gv {
			if GetComponentByName(name) == nil && GetGroupByName(name) == nil {
				return fmt.Errorf("No component or group names' [%s] in graph %s", name, gk)
			}
		}
	}

	return nil
}

func SendMsgToNext(msg *Message) (err error) {
	if msg == nil {
		return fmt.Errorf("Send message nil.")
	}

	next := msg.TopGraph()
	if next == "" {
		return fmt.Errorf("No graph found.")
	}

	return SendMsgToComponent(next, msg)
}

func SendMsgToComponent(name string, msg *Message) (err error) {
	var com *Component = nil
	group := GetGroupByName(name) // 发向一个组吗?
	if group == nil {
		com = components[name]
	} else {
		com = group.GetNode(msg.DispenseKey)
	}
	if com == nil {
		return fmt.Errorf("Get component nil: %s", name)
	}

	msgb := msg.Encode()

	retry := 2 // 失败重试一次
	for err = com.in.SendMessage(msgb); retry > 0 && err != nil; retry -= 1 {
		// retry
	}

	return

}

func GetGroupByName(name string) *ComponentGroup {
	if group, ok := groups[name]; ok {
		return group
	}

	return nil
}

func GetComponentByName(name string) *Component {
	if component, ok := components[name]; ok {
		return component
	}

	return nil
}

func GetGraphByName(name string) []string {
	return graphs[name]
}
