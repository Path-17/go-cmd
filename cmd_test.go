package go_cmd

import(
	"testing"
	"reflect"
	"fmt"
)

// type CmdCommand struct {
// 	CommandName string
// 	Handler     func(map[string]string) (string, error)
// 	Parameters  map[string]reflect.Type
// 	HelpMessage string
// }

func Test_CmdProcess(t *testing.T){
	cmdMap := make(map[string]CmdCommand)
	// add a dummy command
	cmdMap["test"] = CmdCommand {
		CommandName: "test",
		Handler: func(map[string]string) error {
			fmt.Println("Poggers")
			return nil
		},
		Parameters: map[string]reflect.Type{
			"--foo": TypeOf[string](),
			"--bar": TypeOf[bool](),
		},
		HelpMessage: "testhelp",
	}
	CmdInit(cmdMap)
	err := CmdProcess("test --foo 1 --bar")
	if err != nil{
		fmt.Print(err)
	}
}
