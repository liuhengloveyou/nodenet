package nodenet_test

import (
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestBuildFromConfig(t *testing.T) {
	nodenet.BuildFromConfig("example/nodenet.conf.example")
}
