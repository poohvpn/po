package po

type Object map[string]interface{}

func (o *Object) Set(key string, value interface{}) {
	if o == nil {
		return
	}
	if *o == nil {
		*o = make(Object)
	}
	(*o)[key] = value
	return
}

func (o Object) GetBool(key string) (value bool) {
	v, ok := o[key]
	if !ok {
		return
	}
	value, _ = v.(bool)
	return
}

func (o Object) GetInt(key string) (value int) {
	v, ok := o[key]
	if !ok {
		return
	}
	value, _ = v.(int)
	return
}

func (o Object) GetString(key string) (value string) {
	v, ok := o[key]
	if !ok {
		return
	}
	value, _ = v.(string)
	return
}
