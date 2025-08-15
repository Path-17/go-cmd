package go_cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Basic structure:
// Map commands to functions inside of a Map()
// Initialize the Map() with an init() function, pass in a function to it
// Each command is a table lookup, if not found, throw error
// Each command is a struct, with a function, map of params, help message
// Each command must use a specific flag, like --server= , or --port=
// ---> replace all = with spaces
// ---> two types of params, boolean (flag present), and specific (typed in value)

// Global consts
const CMD_MAIN string = "main"

type CMD_HELP_TYPE bool

// Helper function to easily get runtime typing since types are not first class for some reason...
func CmdTypeOf[T any]() reflect.Type {
	var zero T
	return reflect.TypeOf(zero)
}

// Helper function for param parsing
// Pass in the slice of listCmd that doesn't include the command name
// ex: for command "foo --bar=test", with listCmd [foo, --bar, test]
// --> pass in [--bar, test] with the correct specificCmd using list[1:]
func parseParams(listCmd []string, specificCmd CmdCommand) (map[string]CmdParam, error) {
	parsedParams := make(map[string]CmdParam)
	for it := 0; it < len(listCmd); it += 1 {
		// if the paramName doesn't exist in command description, return err
		word := listCmd[it]
		parameterType, ok := specificCmd.Parameters[word]
		if !ok {
			return nil, fmt.Errorf(prettyErrorFormatString, "The parameter \""+word+"\" doesn't exist.")
		}
		// can rename the variable, know that it is a valid param name
		parameterName := word
		// if help flag encountered, exit early and print the help message
		if parameterType.ParamType == CmdTypeOf[CMD_HELP_TYPE]() {
			parsedParams[parameterName] = CmdParam{
				ParamValue: "true",
				ParamType: CmdTypeOf[CMD_HELP_TYPE](),
			}
		// boolean params (flags) don't need values
		} else if parameterType.ParamType == CmdTypeOf[bool]() {
			parsedParams[parameterName] = CmdParam{
				ParamValue: "true",
				ParamType: CmdTypeOf[bool](),
			}
			// all other params are treated as strings
		} else {
			// iterate it to get next word, should be value as current param is not bool type
			it += 1
			// check if the index exists first
			// ---> if not error with "no value supplied for param"
			if it >= len(listCmd) {
				return nil, fmt.Errorf(prettyErrorFormatString, "The string parameter \""+parameterName+"\" was not provided a value.")
			}
			// check if the next index is another param instead of a value
			parameterValue := listCmd[it]
			parsedParams[parameterName] = CmdParam{
				ParamValue: parameterValue,
				ParamType: CmdTypeOf[bool](),
			}
		}
	}

	return parsedParams, nil
}

// Local global only written to by CmdInit
// var registeredCommands map[string]CmdCommand

// Base format string for error messages
// takes a location and a error message, debugging
const errorLocationFormatString = "ERROR in %s: %s"

// Just print the error with basic formatting, user facing
const prettyErrorFormatString = "ERROR: %s"

// Error wrapping format string
const errorWrapFormatString = "ERROR: %w"

// Description of a command, built up by the user and passed in to init inside of a map
// Parameters field is used for easy error handling,
type CmdCommand struct {
	CmdName     string
	Handler     func(map[string]CmdParam) error
	Parameters  map[string]CmdParam
	HelpMessage string
}

type CmdParam struct {
	ParamType reflect.Type
	ParamHelp string
	ParamValue string
}

type CmdApp struct {
	AppName            string
	registeredCommands map[string]CmdCommand
	helpMessage        string
}

func (cmd CmdCommand) CmdHelp() {

	fmt.Println("Params:")

	for paramName, param := range cmd.Parameters {
		fmt.Printf("  %-20s %s\n", paramName, param.ParamHelp)
	}
}

func (app CmdApp) AppHelp() {

	fmt.Printf("%s\n",app.helpMessage)
	if len(app.registeredCommands) == 1 {
		fmt.Printf("\nUsage: %s [params]\n\n", os.Args[0])
		app.registeredCommands[CMD_MAIN].CmdHelp()
	} else {
		for _, v := range app.registeredCommands {
			fmt.Printf("\nUsage: %s %s [params]\n\n", os.Args[0], v.CmdName)
			v.CmdHelp()
		}
	}
}

func (app CmdApp) GetregisteredCommands() map[string]CmdCommand {
	return app.registeredCommands
}

func CmdInitApp(appHelp string) CmdApp {
	return CmdApp{
		registeredCommands: make(map[string]CmdCommand),
		helpMessage:        appHelp,
	}
}

func (app CmdApp) RegisterCommand(cmd CmdCommand) {
	app.registeredCommands[cmd.CmdName] = cmd
}

// Given a command string (un trimmed, no cleanups etc.) run the associated handler from registeredCommands
// Expects the first "word" to be the key to look up in registeredCommands, case sensitive
// Expects parameters to have --<param> value || --<param>=value format
func (app CmdApp) ProcessCommand(rawCmd string) error {
	var err error
	var ok bool

	// clean up the string first
	trimmedCmd := strings.TrimSpace(rawCmd)
	// allow '=' or spaces to be used in commands, treated the same
	normalizedCmd := strings.ReplaceAll(trimmedCmd, "=", " ")
	listCmd := strings.Split(normalizedCmd, " ")
	// Check if help was called

	// if there is only one command registered + default help command, just run the only command without needing to specify it's name
	var specificCmd CmdCommand
	if len(app.registeredCommands) == 1 {
		specificCmd = app.registeredCommands[CMD_MAIN]
	} else {
		// listCmd[0] is the command name if there is more than one command registered
		// find the associated command
		specificCmd, ok = app.registeredCommands[listCmd[0]]
		if !ok {
			return fmt.Errorf(prettyErrorFormatString, "The command "+listCmd[0]+" doesn't exist.")
		}
	}

	var parsedParams map[string]CmdParam
	if len(app.registeredCommands) == 1 {
		parsedParams, err = parseParams(listCmd, specificCmd)
	} else {
		parsedParams, err = parseParams(listCmd[1:], specificCmd)
	}
	if err != nil {
		return err
	}

	// check if help is flagged, if it is print help and return
	for _, param := range parsedParams {
		if param.ParamType == CmdTypeOf[CMD_HELP_TYPE]() {

			app.AppHelp()

			return nil
		}
	}

	// run the command if not nil
	if specificCmd.Handler == nil {
		return fmt.Errorf(prettyErrorFormatString, "Handler function for command not specified (nil)")
	}
	err = specificCmd.Handler(parsedParams)

	return err

}
