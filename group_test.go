package nodenet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestGroupGetNext(t *testing.T) {
	com1, err := nodenet.NewComponent("com1", "rpc", nil)
	fmt.Println(com1, err)

	com2, err := nodenet.NewComponent("com2", "rpc", nil)
	fmt.Println(com2, err)

	g := nodenet.NewGroup("g1", "polling", []*nodenet.Component{com1, com2})
	fmt.Println(g)

	fmt.Println(nodenet.GroupGetNext("g1").Name)
	fmt.Println(nodenet.GroupGetNext("g1").Name)
	fmt.Println(nodenet.GroupGetNext("g1").Name)
	fmt.Println(nodenet.GroupGetNext("g1").Name)
	fmt.Println(nodenet.GroupGetNext("g1").Name)
}
