package nodenet

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	groups map[string][]*Component = make(map[string][]*Component)
	graphs map[string][]string     = make(map[string][]string)
)

var Config struct {
	Name       string              `json:"name"`
	Graphs     map[string][]string `json:"graphs"`
	Groups     []Group             `json:"groups"`
	Components []struct {
		Name string                 `json:"name"`
		In   map[string]interface{} `json:"in"`
	} `json:"components"`
}

type Group struct {
	Name     string   `json:"name"`
	Dispense string   `json:"dispense"`
	Members  []string `json:"members"`
}

func BuildComFromConfig(fileName string) {
	r, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&Config); err != nil {
		panic(err)
	}

	for i := 0; i < len(Config.Components); i++ {
		_, err := NewComponent(Config.Components[i].Name, Config.Components[i].In)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(components)

	for i := 0; i < len(Config.Groups); i++ {
		members := Config.Groups[i].Members
		for j := 0; j < len(members); j++ {
			comp := components[members[j]]
			if comp != nil {
				groups[Config.Groups[i].Name] = append(groups[Config.Groups[i].Name], comp)
			}
		}
	}
	fmt.Println(groups)

	graphs = Config.Graphs
	fmt.Println(graphs)
}
