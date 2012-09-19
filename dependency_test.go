package dependency

import (
	"testing"
)

const testconfig = `
@dep bacon = crispy
@dep herp = derp
I don't even.
`

type TestStruct struct{}

func (t *TestStruct) String() string {
	return "hola manuel"
}

func NewTestStruct() (interface{}, error) {
	return &TestStruct{}, nil
}

func TestAll(t *testing.T) {
	SetConfig(testconfig, "text")
	Add("bacon", "crispy", NewTestStruct)
	Refresh()

	emptyinterface, err := Get("bacon")
        if err != nil || emptyinterface == nil {
            t.Fatalf("Failed to get the struct from constructor: %s", err)
        }
        _, ok := emptyinterface.(*TestStruct)
        if !ok {
            t.Fatal("Couldn't convert empty interface to proper type")
        }
}
