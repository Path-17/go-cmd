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
	cmdApp := CmdInitApp("Welcome to cmdApp testing")
	// add a dummy command
	cmdApp.RegisterCommand(CmdCommand{
		CmdName: CMD_MAIN,
		Handler: func(params map[string]CmdParam) error {
			fmt.Println("Poggers1")
			return nil
		},
		Parameters: map[string]CmdParam{
			"--help": {
				ParamType: CmdTypeOf[CMD_HELP_TYPE](),
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
	})
	cmdApp.RegisterCommand(CmdCommand{
		CmdName: "test",
		Handler: func(params map[string]CmdParam) error {
			fmt.Println("Poggers2")
			return nil
		},
		Parameters: map[string]CmdParam{
			"--help": {
				ParamType: CmdTypeOf[CMD_HELP_TYPE](),
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
	})
	err := cmdApp.ProcessCommand("test --help --foo=test --bar")
	// TODO: add support for quotes!
	// err := CmdProcess("--help --foo='test' --bar")
	if err != nil {
		fmt.Print(err)
	}
}
