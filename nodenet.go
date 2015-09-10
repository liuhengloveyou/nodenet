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

func BuildFromConfig(fileName string) {
	r, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&Config); err != nil {
		panic(err)
	}

	for i := 0; i < len(Config.Components); i++ {
		_, err := NewComponent(Config.Components[i].Name, Config.Components[i].InType, Config.Components[i].InConf)
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(Config.Groups); i++ {
		l := len(Config.Groups[i].Members)
		coms := make([]*Component, l)
		for j := 0; j < l; j++ {
			coms[j] = components[Config.Groups[i].Members[j]]
			if coms[j] == nil {
				panic("No component name's " + Config.Groups[i].Members[j])
			}
		}

		NewGroup(Config.Groups[i].Name, Config.Groups[i].Dispense, coms)
	}

	graphs = Config.Graphs
}

func SendMsgToNext(name string, comsg *Message) (err error) {
	if name == "" {
		return fmt.Errorf("SendTo where?")
	}

	/*
		group := GetGroupComponentByName(name)
		if com == nil {
			com := components[name]
		} else {

		}
		if com == nil {
			return fmt.Errorf("Get component nil: %s", name)
		}

		msg, _ := comsg.Marshal()
		log.Println(com.Name, "SendMsgToNext:", string(msg))

		return com.in.SendMessage(msg)

	*/

	return nil
}

func GetGroupComponentByName(name string) *ComponentGroup {
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
