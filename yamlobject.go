package cyam

type YamlObject map[string]interface{}
type YamlArray  []interface{}

func (yf YamlObject) IsPresent(k string) bool {
	_, ok := yf[k]
	return ok
}

func (yf YamlObject) GetString(k string) string {
	v, ok := yf[k]
	if !ok {
		return ""
	}
	sv, ok := v.(string)
	if !ok {
		return ""
	}
	return sv
}

func (yf YamlObject) GetInt(k string) int {
	v, ok := yf[k]
	if !ok {
		return 0
	}
	si, ok := v.(int)
	if !ok {
		return 0
	}
	return si
}

func (yf YamlObject) GetObject(k string) YamlObject {
	v, ok := yf[k]
	if !ok {
		return nil
	}
	vo, ok := v.(map[string] interface{})
	return vo
}


func (yf YamlObject) GetArray(k string) []interface{} {
	v, ok := yf[k]
	if !ok {
		return nil
	}
	va, ok := v.([]interface{})
	if !ok {
		return nil
	}
	return va
}
