package main

import (
	"fmt"
	"testing"
)

const testCommandOne = "cmdOne"
const testCommandTwo = "cmdTwo"
const testCommandThree = "cmdThree"

const testFlagOne = "flag_one"
const testFlagTwo = "flag_two"
const testFlagThree = "flag_three"

const testParamOne = "param_one"
const testParamTwo = "param_two"
const testParamThree = "param_three"

const testParamValueOne = "one"
const testParamValueTwo = "two"
const testParamValueThree = "three"

var allCmds = []string{testCommandOne, testCommandTwo, testCommandThree}
var allFlags = []string{testFlagOne, testFlagTwo, testFlagThree}

var allParamKeys = []string{testParamOne, testParamTwo, testParamThree}
var allParamValues = []string{testParamValueOne, testParamValueTwo, testParamValueThree}
var allParams = func() []Parameter {
	ap := make([]Parameter, len(allParamKeys))
	for i, k := range allParamKeys {
		ap[i] = Parameter{
			Key:   k,
			Value: allParamValues[i],
		}
	}
	return ap
}()


var allArgs = func() []string {
	aa := make([]string, 0, len(allCmds) + len(allFlags) + len(allParams))
	aa = append(aa, allCmds...)
	for _, f := range allFlags {
		aa = append(aa, fmt.Sprintf("-%s", f))
	}
	for _, pm := range allParams {
		aa = append(aa, fmt.Sprintf("--%s", pm.Key))
		aa = append(aa, pm.Value)
	}
	return aa
} ()

func TestNewArguments(t *testing.T) {
	args := NewArguments([]string{testCommandOne})
	if len(args.Commands) != 1 {
		t.Fatalf("expected %d commands, found %d", 1, len(args.Commands))
	}
	if len(args.Flags) != 0 {
		t.Fatalf("expected %d flags, found %d", 0, len(args.Flags))
	}
	if len(args.Params) != 0 {
		t.Fatalf("expected %d parameters, found %d", 0, len(args.Params))
	}
	if args.Commands[0] != testCommandOne {
		t.Fatalf("expected %s command, found %s", testCommandOne, args.Commands[0])
	}


	args = NewArguments(allArgs)
	if len(args.Commands) != len(allCmds) {
		t.Fatalf("expected %d commands, found %d", len(allCmds), len(args.Commands))
	}
	if len(args.Flags) != len(allFlags) {
		t.Fatalf("expected %d flags, found %d", len(allFlags), len(args.Flags))
	}
	if len(args.Params) != len(allParams) {
		t.Fatalf("expected %d parameters, found %d", len(allParams), len(args.Params))
	}

	if args.Commands[0] != testCommandOne {
		t.Fatalf("expected %s command, found %s", testCommandOne, args.Commands[0])
	}

	if args.Flags[0] != testFlagOne {
		t.Fatalf("expected %s flag, found %s", testFlagOne, args.Flags[0])
	}

	if args.Params[0].Key != testParamOne {
		t.Fatalf("expected %s parameter, found %s", testFlagOne, args.Params[0].Key)
	}

	if args.Params[0].Value != testParamValueOne {
		t.Fatalf("expected %s parameter value, found %s", testParamValueOne, args.Params[0].Value)
	}
}

func TestArguments_HasFlag(t *testing.T) {
	args := NewArguments(nil)
	exFalse := args.HasFlag(testFlagOne)
	if exFalse {
		t.Fatalf("expdcted false for hasFlag with %s, on empty argument set, found true", testFlagOne)
	}

	args = NewArguments([]string{ "-" + testFlagOne})
	if len(args.Flags) != 1 {
		t.Fatalf("expected flag count of %d, found %d", 1, len(args.Flags))
	}
	if !args.HasFlag(testFlagOne) {
		t.Fatalf("expected true for hasFlag with %s, on argument set containing that flag, found false", testFlagOne)
	}
	if args.HasFlag("") {
		t.Fatalf("expected false for hasFlag with empty key, found true")
	}

	args = NewArguments(allArgs)
	if len(args.Flags) != 3 {
		t.Fatalf("expected flag count of %d, found %d", 3, len(args.Flags))
	}
	if args.HasFlag("") {
		t.Fatalf("expected false for hasFlag with empty key, found true")
	}
	if args.HasFlag(testParamOne) {
		t.Fatalf("expected false for hasFlag with test parameter key, found true")
	}

	for _, f := range allFlags {
		if !args.HasFlag(f) {
			t.Fatalf("expected true for hasFlag with %s, on argument set containing that flag, found false", f)
		}
	}
}


func TestArguments_HasParameter(t *testing.T) {
	args := NewArguments(nil)
	exFalse := args.HasParameter(testParamOne)
	if exFalse {
		t.Fatalf("expected false for HasParameter with %s, on empty argument set, found true", testParamOne)
	}

	args = NewArguments([]string{ "--" + testParamOne})
	if len(args.Params) != 1 {
		t.Fatalf("expected parameter count of %d, found %d", 1, len(args.Params))
	}
	if !args.HasParameter(testParamOne) {
		t.Fatalf("expected true for HasParameter with %s, on argument set containing that parameter, found false", testParamOne)
	}
	if args.HasParameter("") {
		t.Fatalf("expected false for HasParameter with empty key, found true")
	}

	args = NewArguments(allArgs)
	if len(args.Params) != 3 {
		t.Fatalf("expected parameter count of %d, found %d", 3, len(args.Params))
	}
	if args.HasParameter("") {
		t.Fatalf("expected false for HasParameter with empty key, found true")
	}
	if args.HasParameter(testFlagOne) {
		t.Fatalf("expected false for HasParameter with test flag key, found true")
	}

	for _, p := range allParams {
		if !args.HasParameter(p.Key) {
			t.Fatalf("expected true for HasParameter with %s, on argument set containing that parameter, found false", p.Key)
		}
		v, ok := args.GetParameter(p.Key)
		if !ok {
			t.Fatalf("expected true for GetParameter with %s, on argument set containing that parameter, found false", p.Key)
		}
		if v != p.Value {
			t.Fatalf("expected %s value for GetParameter with %s, on argument set containing that parameter, found %s", p.Value, p.Key, v)
		}
	}
}
