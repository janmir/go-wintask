package tasker

import (
	"fmt"
	"testing"
)

var (
	tasker = New()
)

func TestQuery(t *testing.T) {
	output := tasker.Query("adobe", false)
	fmt.Printf("%+v\n", output)
}
