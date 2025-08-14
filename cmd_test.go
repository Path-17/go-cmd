package go_cmd

import (
	"fmt"
	// "reflect"
	"testing"
)

// type CmdCommand struct {
// 	CommandName string
// 	Handler     func(map[string]string) (string, error)
// 	Parameters  map[string]reflect.Type
// 	HelpMessage string
// }

func Test_CmdProcess(t *testing.T) {
	cmdMap := make(map[string]CmdCommand)
	// add a dummy command
	cmdMap[CMD_MAIN] = CmdCommand{
		Handler: func(params map[string]string) error {
			fmt.Printf("%v", params)
			fmt.Println("Poggers")
			return nil
		},
		Parameters: map[string]CmdParamMetadata{
			"--help": {
				ParamType: CmdTypeOf[bool](),
				ParamHelp: "Helphelp",
			},
			"--foo": {
				ParamType: CmdTypeOf[string](),
				ParamHelp: "Helpfoo",
			},
			"--bar": {
				ParamType: CmdTypeOf[bool](),
				ParamHelp: "Helpbar",
			},
		},
		HelpMessage: "testhelp",
		HelpHandler: CmdDefaultHelp,
	}
	// cmdMap["test"] = CmdCommand{
	// 	Handler: func(params map[string]string) error {
	// 		fmt.Printf("%v", params)
	// 		fmt.Println("Poggers2")
	// 		return nil
	// 	},
	// 	Parameters: map[string]CmdParamMetadata{
	// 		"--foo2": {
	// 			ParamType: CmdTypeOf[string](),
	// 			ParamHelp: "Helpfoo",
	// 		},
	// 		"--bar2": {
	// 			ParamType: CmdTypeOf[bool](),
	// 			ParamHelp: "Helpbar",
	// 		},
	// 	},
	// 	HelpMessage: "testhelp2",
	// }
	// cmdMap["test2"] = CmdCommand {
	// 	CommandName: "test2",
	// 	Handler: func(map[string]string) error {
	// 		fmt.Println("Poggers2")
	// 		return nil
	// 	},
	// 	Parameters: map[string]reflect.Type{
	// 		"--foo": CmdTypeOf[string](),
	// 		"--bar": CmdTypeOf[bool](),
	// 	},
	// 	HelpMessage: "testhelp",
	// }
	CmdInit(cmdMap)
	err := CmdProcess("--help --foo=test --bar")
	// TODO: add support for quotes!
	// err := CmdProcess("--help --foo='test' --bar")
	if err != nil {
		fmt.Print(err)
	}
}
