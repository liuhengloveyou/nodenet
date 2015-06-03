package nodenet

import (
	"encoding/json"
	"os"
	"time"
)

const (
	HEARTBEAT = time.Duration(3) * time.Second
)

var (
	components map[string]*Component = make(map[string]*Component)
	groups     map[string]*Group     = make(map[string]*Group)
	graphs     map[string][]string   = make(map[string][]string)
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
		Name   string                 `json:"name"`
		InType string                 `json:"intype"`
		InConf map[string]interface{} `json:"inconf"`
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
