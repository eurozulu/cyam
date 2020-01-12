package main

import (
	"fmt"
	"github.com/eurozulu/cyam"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	args := NewArguments(os.Args[1:])
	if args.HasFlag("h") || args.HasParameter("help") {
		showUse()
		os.Exit(0)
	}
	if len(args.Commands) < 2 {
		showUse()
		os.Exit(1)
	}

	pt, err := loadMatchPattern(args.Commands[1])
	if err != nil {
		panic(err)
	}
	if args.HasFlag("v") {
		fmt.Printf("searching for key using pattern: %s\n", pt.String())
	}

	var yo cyam.YamlObject
	if args.Commands[0] == "-" {
		yo, err = loadYamlStream()
	} else {
		yo, err = loadYamlFile(args.Commands[0])
	}
	if err != nil {
		panic(err)
	}

	wk := cyam.NewWalker(pt, os.Stdout)
	wk.IncludeKey = (args.HasFlag("k") || args.HasParameter("keys"))

	wk.Walk(yo)
}

func loadMatchPattern(p string) (cyam.Pattern, error) {
	return cyam.NewPathPattern(p)
}

func loadYamlFile(p string) (cyam.YamlObject, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var yo cyam.YamlObject
	err = yaml.NewDecoder(f).Decode(&yo)
	if err != nil {
		return nil, err
	}
	return yo, nil
}

func loadYamlStream() (cyam.YamlObject, error) {
	var yo cyam.YamlObject
	err := yaml.NewDecoder(os.Stdin).Decode(&yo)
	if err != nil {
		return nil, err
	}
	return yo, nil
}

func showUse() {
	fmt.Printf("%s <yaml file name> <yaml path>\n", os.Args[0])
	fmt.Println("\tyaml file name\t\tThe path to the yalm file")
	fmt.Println("\tyaml path\t\tThe yaml search path")
	fmt.Println("\t\tone.two.three	\tfinds will return the value of 'one.two.three' if present")
	fmt.Println("\t\t*.two.three	\tfinds will return the values of all properties ending with '.two.three' if present")
	fmt.Println("\t\tone.two.*		\tfinds will return the key names and values of all properties starting with 'one.two.' if present")
	fmt.Println("\t\tone.*three	\tfinds will return the key names and values of all properties starting with 'one.'. which end with 'three' if present")
	fmt.Println("\t\tone.*.three	\tfinds will return the key names and values of all properties named 'three' with any parent starting with 'one.' if present")
	fmt.Println()
	fmt.Println("a single * indicates a single key name: so one.*.three = one.two.three, one.2.three and one.haha.three")
	fmt.Println("but NOT one.two.twohalf.three, one.one.onw.three")
	fmt.Println("double ** indicates any key names: so one.**.three = one.two.three, one.2.two.three and one.haha.hehe.hoho.three")
	fmt.Println("but NOT one.three.one, three.one.one.one")
	fmt.Println()
	fmt.Println("Objects/maps inside arrays are included in ** searches.")
	fmt.Println("Arrays may be searched using [n] where n is a specific index or * for all indexes.")

}
