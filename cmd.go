package go_cmd

import (
	"fmt"
	"maps"
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

// Helper function to easily get runtime typing since types are not first class for some reason...
func TypeOf[T any]() reflect.Type {
	var zero T
	return reflect.TypeOf(zero)
}

// Local global only written to by CmdInit
var registeredCommands map[string]CmdCommand

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
	CommandName string
	Handler     func(map[string]string) error
	Parameters  map[string]reflect.Type
	HelpMessage string
}

// Always run as configuration step
func CmdInit(cmdMap map[string]CmdCommand) {
	registeredCommands = make(map[string]CmdCommand)
	clear(registeredCommands)
	maps.Copy(registeredCommands, cmdMap)
}

// Given a command string (un trimmed, no cleanups etc.) run the associated handler from registeredCommands
// Expects the first "word" to be the key to look up in registeredCommands, case sensitive
// Expects parameters to have --<param> value || --<param>=value format
func CmdProcess(rawCmd string) error {
	// clean up the string first
	trimmedCmd := strings.TrimSpace(rawCmd)
	// allow '=' or spaces to be used in commands, treated the same
	normalizedCmd := strings.ReplaceAll(trimmedCmd, "=", " ")
	listCmd := strings.Split(normalizedCmd, " ")

	// listCmd[0] is the command name
	specificCmd, ok := registeredCommands[listCmd[0]]
	if !ok {
		return fmt.Errorf(prettyErrorFormatString, "The command "+listCmd[0]+" doesn't exist.")
	}

	// make sure that the params are valid
	// bool params have no values
	// regular params have one value
	parsedParams := make(map[string]string)
	for it := 1; it < len(listCmd); it += 1 {
		// if the paramName doesn't exist in command description, return err
		word := listCmd[it]
		parameterType, ok := specificCmd.Parameters[word]
		if !ok {
			return fmt.Errorf(prettyErrorFormatString, "The parameter "+word+" doesn't exist.")
		}
		// can rename the variable, know that it is a valid param name
		parameterName := word
		// boolean params (flags) don't need values
		if parameterType == TypeOf[bool]() {
			parsedParams[parameterName] = "true"
		// all other params are treated as strings
		} else {
			// iterate it to get next word, should be value as param is not bool type
			it += 1
			// check if the index exists first
			// error with "no value supplied for param"
			if it >= len(listCmd) {
				return fmt.Errorf(prettyErrorFormatString, "The string parameter \"" + parameterName + "\" was not provided a value.")
			}
			parameterValue := listCmd[it]
			parsedParams[parameterName] = parameterValue
		}
	}

	// run the command if not nil
	if specificCmd.Handler == nil {
		return fmt.Errorf(prettyErrorFormatString, "Handler function for command \"" + specificCmd.CommandName + "\" not specified (nil)")
	}
	err := specificCmd.Handler(parsedParams)
	if err != nil {
		return fmt.Errorf(errorWrapFormatString, err)
	}

	return nil

}
