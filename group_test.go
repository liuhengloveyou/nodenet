package nodenet_test

import (
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestGroup(t *testing.T) {
	com1, _ := nodenet.NewComponent("com1", "tcp", nil)
	com2, _ := nodenet.NewComponent("com2", "tcp", nil)
	com3, _ := nodenet.NewComponent("com3", "tcp", nil)
	com4, _ := nodenet.NewComponent("com4", "tcp", nil)
	com5, _ := nodenet.NewComponent("com5", "tcp", nil)

	g := nodenet.NewGroup("group1", "polling", []*nodenet.Component{com1, com2, com3, com4, com5})
	t.Log(g.GetNode(""))
	t.Log(g.GetNode(""))
	t.Log(g.GetNode(""))
	t.Log(g.GetNode(""))
	t.Log(g.GetNode(""))
	t.Log(g.GetNode(""))
	t.Log(g.GetNode(""))

	g = nodenet.NewGroup("group1", "hash", []*nodenet.Component{com1, com2, com3, com4, com5})
	t.Log("")
	t.Log(g.GetNode("a123"))
	t.Log(g.GetNode("b234"))
	t.Log(g.GetNode("c345"))
	t.Log(g.GetNode("d456"))
	t.Log(g.GetNode("e567"))
	t.Log(g.GetNode("f678"))
	t.Log(g.GetNode("g789"))
}
