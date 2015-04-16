package cloudnet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/cloudnet"
)

func TestComponent(t *testing.T) {
	com1, err := cloudnet.NewComponent("com1", "rpc", ":1234")
	fmt.Println(com1, err)
	com1.Run()
}
