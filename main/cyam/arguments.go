package main

import (
	"fmt"
	"strings"
)

type Arguments struct {
	Flags    []string    // -x Arguments
	Params   []Parameter // --x <value> Arguments
	Commands []string    // all others
}

func (args Arguments) HasParameter(k string) bool {
	for _, p := range args.Params {
		if strings.EqualFold(p.Key, k) {
			return true
		}
	}
	return false
}

func (args Arguments) GetParameter(k string) (string, bool) {
	for _, p := range args.Params {
		if strings.EqualFold(p.Key, k) {
			return p.Value, true
		}
	}
	return "", false
}

func (args Arguments) HasFlag(k string) bool {
	for _, f := range args.Flags {
		if strings.EqualFold(f, k) {
			return true
		}
	}
	return false
}

type Parameter struct {
	Key   string
	Value string
}
func (p Parameter) String() string {
	v := p.Value
	if v != "" {
		v = " " + v
	}
	return fmt.Sprintf("--%s%s", p.Key, v)
}
func NewArguments(args []string) *Arguments {
	arguments := Arguments{}
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			key := args[i][2:]
			val := ""
			if i + 1 < len(args) {
				val = args[i + 1]
				i++
			}
			arguments.Params = append(arguments.Params, Parameter{
				Key:   key,
				Value: val,
			})
		} else if strings.HasPrefix(args[i], "-") {
			// dash on its own counts as a command (for using std input instead of a filename)
			if len(args[i]) > 1 {
				arguments.Flags = append(arguments.Flags, args[i][1:])
			} else {
				arguments.Commands = append(arguments.Commands, args[i])
			}
		} else {
			arguments.Commands = append(arguments.Commands, args[i])
		}
	}
	return &arguments
}
