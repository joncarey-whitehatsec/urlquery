package urlquery

import (
	"testing"
	"strings"
	"fmt"
)

type TestParseChild struct {
	Description string `query:"desc"`
	Long        uint16 `query:" vip"`
	Height      int    `query:" ignore"`
}

type TestParseInfo struct {
	Id       int
	Name     string           `query:"name"`
	Child    TestParseChild   `query:"child"`
	ChildPtr *TestParseChild  `query:"childPtr"`
	Children []TestParseChild `query:"children"`
	Params   map[byte]int8
	status   bool
	UintPtr  uintptr
}

func TestUnmarshal_NestedStructure(t *testing.T) {
	var data = "Id=1&name=test&child[desc]=c1&child[Long]=10&childPtr[Long]=2&childPtr[Description]=b" +
		"&children[0][desc]=d1&children[1][Long]=12&children[5][desc]=d5&children[5][Long]=50&desc=rtt" +
		"&Params[120]=1&Params[121]=2&status=1&UintPtr=300"
	data = encodeSquareBracket(data)
	v := &TestParseInfo{}
	err := Unmarshal([]byte(data), v)

	if err != nil {
		t.Error(err)
	}

	if v.Id != 1 {
		t.Error("Id wrong")
	}

	if v.Name != "test" {
		t.Error("Name wrong")
	}

	if v.Child.Description != "c1" || v.Child.Long != 10 || v.Child.Height != 0 {
		t.Error("Child wrong")
	}

	if v.ChildPtr == nil || v.ChildPtr.Description != "" || v.ChildPtr.Long != 2 || v.ChildPtr.Height != 0 {
		t.Error("ChildPtr wrong")
	}

	if len(v.Children) != 6 {
		t.Error("Children's length is wrong")
	}

	if v.Children[0].Description != "d1" {
		t.Error("Children[0] wrong")
	}

	if v.Children[1].Description != "" || v.Children[1].Long != 12 {
		t.Error("Childre[1] wrong")
	}

	if v.Children[2].Description != "" || v.Children[3].Description != "" || v.Children[4].Description != "" {
		t.Error("Children[2,3,4] wrong")
	}

	if v.Children[5].Description != "d5" || v.Children[5].Long != 50 || v.Children[5].Height != 0 {
		t.Error("Children[5] wrong")
	}

	if len(v.Params) != 2 || v.Params[120] != 1 || v.Params[121] != 2 {
		t.Error("Params wrong")
	}

	if v.status != false {
		t.Error("status wrong")
	}

	if v.UintPtr != uintptr(300) {
		t.Error("UintPtr wrong")
	}
}

func TestUnmarshal_Map(t *testing.T) {
	var m map[string]string
	data := "id=1&name=ab&arr[0]=6d"
	data = encodeSquareBracket(data)
	err := Unmarshal([]byte(data), &m)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(m)

	if len(m) != 2 {
		t.Error("length is wrong")
	}
	if v1, ok1 := m["id"]; v1 != "1" || !ok1 {
		t.Error("map[id] is wrong")
	}
	if v2, ok2 := m["name"]; v2 != "ab" || !ok2 {
		t.Error("map[iname] is wrong")
	}
	if _, ok3 := m["arr%5B0%5D"]; ok3 {
		t.Error("map[arr%5B0%5D] should not be exist")
	}
}

type TestFormat struct {
	Id uint64
	B  rune `query:"b"`
}

func TestUnmarshal_UnmatchedDataFormat(t *testing.T) {
	var data = "Id=1&b=a"
	data = encodeSquareBracket(data)
	v := &TestFormat{}
	err := Unmarshal([]byte(data), v)

	if err == nil {
		t.Error("error should not be ignored")
	}
	if _, ok := err.(ErrTranslated); !ok {
		t.Errorf("error type is unexpected. %v", err)
	}
}

func TestUnmarshal_UnhandledType(t *testing.T) {
	var data = "Id=1&b=a"
	data = encodeSquareBracket(data)
	v := &map[interface{}]string{}
	err := Unmarshal([]byte(data), v)

	if err == nil {
		t.Error("error should not be ignored")
	}
	if _, ok := err.(ErrInvalidMapKeyType); !ok {
		t.Errorf("error type is unexpected. %v", err)
	}
}

type TestUnhandled struct {
	Id     int
	Params map[string]TestFormat
}

func TestUnmarshal_UnhandledType2(t *testing.T) {
	var data = "Id=1&b=a"
	data = encodeSquareBracket(data)
	v := &TestUnhandled{}
	err := Unmarshal([]byte(data), v)

	if err == nil {
		t.Error("error should not be ignored")
	}
	if _, ok := err.(ErrInvalidMapValueType); !ok {
		t.Errorf("error type is unexpected. %v", err)
	}
}

//benchmark: mock multi-layer nested structure,
//BenchmarkUnmarshal-4   	  210648	     16986 ns/op
func BenchmarkUnmarshal(b *testing.B) {
	var data = "Id=1&name=test&child[desc]=c1&child[Long]=10&childPtr[Long]=2&childPtr[Description]=b" +
		"&children[0][desc]=d1&children[1][Long]=12&children[5][desc]=d5&children[5][Long]=50&desc=rtt" +
		"&Params[120]=1&Params[121]=2&status=1&UintPtr=300"
	data = encodeSquareBracket(data)

	for i := 0; i < b.N; i++ {
		v := &TestParseInfo{}
		err := Unmarshal([]byte(data), v)
		if err != nil {
			b.Error(err)
		}
	}
}

func encodeSquareBracket(data string) string {
	data = strings.ReplaceAll(data, "[", "%5B")
	data = strings.ReplaceAll(data, "]", "%5D")
	return data
}
