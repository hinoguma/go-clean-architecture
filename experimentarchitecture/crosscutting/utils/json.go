package utils

import "fmt"

type JsonString string

func (v JsonString) String() string {
	return string(v)
}

/*****************************
	Json Builder
 *****************************/

type Json struct {
	data interface{}
}

func (js *Json) Cursor() *JsonCursor {
	return NewJsonCursor(*js)
}

type JsonKeyType string

const (
	JsonKeyTypeKey   JsonKeyType = "key"
	JsonKeyTypeIndex JsonKeyType = "index"
)

type JsonKey struct {
	Type  JsonKeyType
	key   string
	index int
}

func (key JsonKey) IsKey() bool {
	return key.Type == JsonKeyTypeKey
}
func (key JsonKey) IsIndex() bool {
	return key.Type == JsonKeyTypeIndex
}

type JsonValue struct {
	value interface{}
}

type JsonCursor struct {
	json    Json
	pointer []JsonKey
	now     *JsonValue
}

func NewJsonCursor(js Json) *JsonCursor {
	data := js.data
	return &JsonCursor{
		json:    Json{data: data},
		pointer: make([]JsonKey, 0),
	}
}

func (c *JsonCursor) Key(key string) *JsonCursor {
	c.pointer = append(c.pointer, JsonKey{
		Type: JsonKeyTypeKey,
		key:  key,
	})

	if c.now == nil {
		c.now = &JsonValue{
			value: c.json.data,
		}
	}
	m, ok := c.now.value.(map[string]interface{})
	if ok {
		if v, exists := m[key]; exists {
			c.now.value = v
		}
	}

	return c
}
func (c *JsonCursor) Index(index int) *JsonCursor {
	c.pointer = append(c.pointer, JsonKey{
		Type:  JsonKeyTypeIndex,
		index: index,
	})
	return c
}
func (c *JsonCursor) ResetPointer() *JsonCursor {
	c.pointer = make([]JsonKey, 0)
	return c
}

func (c *JsonCursor) Value() interface{} {
	if c.json.data == nil {
		return nil
	}
	prev := c.json.data
	lastIndex := len(c.pointer) - 1
	for i, key := range c.pointer {
		if i == lastIndex {
			if key.IsKey() {
				m, ok := prev.(map[string]interface{})
				if !ok {
					return nil
				}
				return m[key.key]
			} else if key.IsIndex() {
				arr, ok := prev.([]interface{})
				if !ok {
					return nil
				}
				if key.index < 0 || key.index >= len(arr) {
					return nil
				}
				return arr[key.index]
			}
			return nil
		}
		if key.IsKey() {
			prevMap, ok := prev.(map[string]interface{})
			if !ok {
				return nil
			}
			prev = prevMap[key.key]
		} else if key.IsIndex() {
			prevArr, ok := prev.([]interface{})
			if !ok {
				return nil
			}
			if key.index < 0 || key.index >= len(prevArr) {
				return nil
			}
			prev = prevArr[key.index]
		}

	}
	return nil
}

func (c *JsonCursor) GetJson() Json {
	return c.json
}

func (c *JsonCursor) SetByPointer(value interface{}) *JsonCursor {
	if c.json.data == nil {
		c.json.data = make(map[string]interface{})
	}
	previous := c.json.data
	lastIndex := len(c.pointer) - 1
	for i, key := range c.pointer {
		if key.IsKey() {
			prevMap := previous.(map[string]interface{})
			if i == lastIndex {
				prevMap[key.key] = value
				break
			} else {
				if c.pointer[i+1].IsKey() {
					_, ok := prevMap[key.key].(map[string]interface{})
					if !ok {
						prevMap[key.key] = make(map[string]interface{})
					}
				} else if c.pointer[i+1].IsIndex() {
					_, ok := prevMap[key.key].([]interface{})
					if !ok {
						prevMap[key.key] = make([]interface{}, 0)
					}
				}
				previous = prevMap[key.key]
			}
		}
		if key.IsIndex() {
			prevArr := previous.([]interface{})
			for j := len(prevArr); j <= key.index; j++ {
				prevArr = append(prevArr, nil)
			}
			if i == lastIndex {
				prevArr[key.index] = value
				break
			} else {
				if c.pointer[i+1].IsKey() {
					_, ok := prevArr[key.index].(map[string]interface{})
					if !ok {
						prevArr[key.index] = make(map[string]interface{})
					}
				} else if c.pointer[i+1].IsIndex() {
					_, ok := prevArr[key.index].([]interface{})
					if !ok {
						prevArr[key.index] = make([]interface{}, 0)
					}
				}
				previous = prevArr[key.index]
			}
		}
	}
	return c
}

func (c *JsonCursor) SetByBucket(value interface{}) *JsonCursor {
	b := bucket{items: make([]bucketItem, 0)}
	if c.json.data == nil {
		c.json.data = make(map[string]interface{})
	}
	prev := c.json.data
	lastIndex := len(c.pointer) - 1
	for i, key := range c.pointer {
		if key.IsKey() {
			mapVal, ok := prev.(map[string]interface{})
			if !ok {
				mapVal = make(map[string]interface{})
			}
			if i == lastIndex {
				mapVal[key.key] = value
			}
			current, ok := mapVal[key.key]
			if !ok {
				current = nil
			}
			prev = current
			b.AddItem(key, current)
		} else if key.IsIndex() {
			av, ok := prev.([]interface{})
			if !ok {
				av = make([]interface{}, key.index+1)
			}
			for j := len(av); j <= key.index; j++ {
				av = append(av, nil)
			}
			if i == lastIndex {
				av[key.index] = value
			}
			current := av[key.index]
			prev = current
			b.AddItem(key, current)
		}
	}
	//fmt.Println(b)
	lastIndex = len(b.items) - 1
	for i := lastIndex - 1; i >= 0; i-- {
		item := b.items[i]
		childItem := b.items[i+1]
		if childItem.key.IsKey() {
			m, ok := item.value.(map[string]interface{})
			//fmt.Println(m, ok)
			if !ok {
				m = make(map[string]interface{})
			}
			m[childItem.key.key] = childItem.value
			b.items[i].value = m
			item.value = m
		} else if childItem.key.IsIndex() {
			arr, ok := item.value.([]interface{})
			//fmt.Println(arr, ok)
			if !ok {
				arr = make([]interface{}, childItem.key.index+1)
			}
			for j := len(arr); j <= childItem.key.index; j++ {
				arr = append(arr, nil)
			}
			arr[childItem.key.index] = childItem.value
			b.items[i].value = arr
			//item.value = arr

		}

		if i == 0 {
			if item.key.IsKey() {
				m, ok := c.json.data.(map[string]interface{})
				if !ok {
					m = make(map[string]interface{})
				}
				m[item.key.key] = item.value
				c.json.data = m
			} else if item.key.IsIndex() {
				arr, ok := c.json.data.([]interface{})
				if !ok {
					arr = make([]interface{}, item.key.index+1)
				}
				for j := len(arr); j <= item.key.index; j++ {
					arr = append(arr, nil)
				}
				arr[item.key.index] = item.value
				//c.json.data = arr
			}
		}
	}

	return c
}

type bucket struct {
	items []bucketItem
}

type bucketItem struct {
	key   JsonKey
	value interface{}
}

func (b *bucket) AddItem(key JsonKey, value interface{}) {
	b.items = append(b.items, bucketItem{
		key:   key,
		value: value,
	})
}

type JsonExtractor struct {
	src interface{}
}

func NewJsonExtractor(src interface{}) *JsonExtractor {
	return &JsonExtractor{
		src: src,
	}
}

func (e *JsonExtractor) ByKey(key string) (interface{}, bool) {
	m, ok := e.src.(map[string]interface{})
	return m[key], ok
}

func (e *JsonExtractor) ByIndex(index int) (interface{}, bool) {
	if index < 0 {
		return nil, false
	}
	arr, ok := e.src.([]interface{})
	return arr, ok
}

func Test() {
	complexMap1 := map[string]map[string]string{
		"a": {
			"b": "c",
		},
	}
	//modifyMap(complexMap1)
	//fmt.Println(complexMap1)
	modifyMap2(complexMap1)
	fmt.Println(complexMap1)

}

func modifyMap(m map[string]map[string]string) {
	a := m["a"]
	a["b"] = "d"
	a = map[string]string{
		"c": "g",
	}
}

func modifyMap2(m map[string]map[string]string) {
	ab := m["a"]["b"]
	abp := &ab
	*abp = "d"

}

func modifyVal(target interface{}, val string) {
	target = val
}
