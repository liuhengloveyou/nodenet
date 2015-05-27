package nodenet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestComponent(t *testing.T) {
	com1, err := nodenet.NewComponent("com1", nil)
	fmt.Println(com1, err)
	com1.Run()
}
