package nodenet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestBuildFromConfig(t *testing.T) {
	nodenet.BuildFromConfig("example/nodenet.conf.example")
}

func TestNodenet(t *testing.T) {
	nodenet.BuildFromConfig("example/nodenet.conf.example")

	err := nodenet.GetComponentByName("com1").Run()
	fmt.Println(err)
}
