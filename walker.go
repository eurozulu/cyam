package cyam

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

type Pattern interface {
	String() string
	Match(k string) bool
}

type Walker struct {
	MatchPattern Pattern
	Out          io.Writer
	yo           *yaml.Encoder
	IncludeKey   bool
}

func NewWalker(pattern Pattern, out io.Writer) *Walker {
	yo := yaml.NewEncoder(out)
	return &Walker{
		MatchPattern: pattern,
		Out:          out,
		yo:           yo,
	}
}

func (w *Walker) Walk(o YamlObject) {
	if w.yo == nil {
		w.yo = yaml.NewEncoder(w.Out)
	}
	w.walkObject("", o)
}

func (w Walker) walkArray(p string, a YamlArray) {
	for i, ae := range a {
		cp := fmt.Sprintf("%s[%d", p, i)

		if w.MatchPattern.Match(cp) {
			w.writeOut(cp, ae)
		}
		w.walkValue(cp, ae)
	}
}

func (w Walker) walkObject(p string, o YamlObject) {
	for k, v := range o {
		cp := k
		if p != "" {
			cp = strings.Join([]string{p, k}, ".")
		}
		if w.MatchPattern.Match(cp) {
			w.writeOut(cp, v)
		}
		w.walkValue(cp, v)
	}
}

func (w Walker) walkValue(p string, v interface{}) {
	switch c := v.(type) {
	case YamlObject:
		w.walkObject(p, c)
	case map[string]interface{}:
		w.walkObject(p, c)
	case map[interface{}]interface{}:
		w.walkObject(p, w.typeObject(c))

	case YamlArray:
		w.walkArray(p, c)
	case []interface{}:
		w.walkArray(p, c)

	default: // strings, numbers and other dead ends.
		return
	}
}

func (w Walker) writeOut(k string, v interface{}) error {
	if w.IncludeKey {
		m := make(map[string]interface{})
		m[k] = v
		v = m
	}
	if err := yaml.NewEncoder(w.Out).Encode(v); err != nil {
		return err
	}
	//_, err := w.Out.Write([]byte{'\n'})
	return nil
}

func (w Walker) typeObject(m map[interface{}]interface{}) YamlObject {
	by, err := yaml.Marshal(&m)
	if err != nil {
		panic(err)
	}
	var yo YamlObject
	err = yaml.Unmarshal(by, &yo)
	if err != nil {
		panic(err)
	}
	return yo
}
