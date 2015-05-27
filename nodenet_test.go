package nodenet_test

import (
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestBuildComFromConfig(t *testing.T) {
	nodenet.BuildComFromConfig("example/nodenet.conf.example")

}
