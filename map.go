package po

import (
	"sync"
)

type Map sync.Map

func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	return (*sync.Map)(m).Load(key)
}

func (m *Map) Store(key, value interface{}) {
	(*sync.Map)(m).Store(key, value)
}

func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	return (*sync.Map)(m).LoadOrStore(key, value)
}

func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	return (*sync.Map)(m).LoadAndDelete(key)
}

func (m *Map) Delete(key interface{}) {
	(*sync.Map)(m).Delete(key)
}

func (m *Map) Range(f func(key, value interface{}) bool) {
	(*sync.Map)(m).Range(f)
}

func (m *Map) Len() int {
	i := 0
	m.Range(func(_, _ interface{}) bool {
		i++
		return true
	})
	return i
}

func (m *Map) Keys() []interface{} {
	mapLen := m.Len()
	if mapLen == 0 {
		return nil
	}
	keys := make([]interface{}, 0, mapLen)
	m.Range(func(key, _ interface{}) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

func (m *Map) Values() []interface{} {
	mapLen := m.Len()
	if mapLen == 0 {
		return nil
	}
	values := make([]interface{}, 0, mapLen)
	m.Range(func(_, value interface{}) bool {
		values = append(values, value)
		return true
	})
	return values
}

func (m *Map) LoadString(key interface{}) (string, bool) {
	v, ok := m.Load(key)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (m *Map) LoadInt(key interface{}) (int, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int), true
}

func (m *Map) LoadUint16(key interface{}) (uint16, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(uint16), true
}

func (m *Map) LoadUint32(key interface{}) (uint32, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(uint32), true
}

func (m *Map) LoadUint64(key interface{}) (uint64, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(uint64), true
}

func (m *Map) LoadInt64(key interface{}) (int64, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int64), true
}

func (m *Map) LoadBytes(key interface{}) ([]byte, bool) {
	v, ok := m.Load(key)
	if !ok {
		return nil, false
	}
	return v.([]byte), true
}

func (m *Map) LoadFloat32(key interface{}) (float32, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(float32), true
}

func (m *Map) LoadFloat64(key interface{}) (float64, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(float64), true
}
