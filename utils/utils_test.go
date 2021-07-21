package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
	s := struct{ Test string }{Test: "test"}
	t.Run("Hash is always same", func(t *testing.T) {
		x := Hash(s)
		if hash != x {
			t.Errorf("Expected %s, But got %s", hash, x)
		}
	})
	t.Run("Hash is encode hex.", func(t *testing.T) {
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Error("Hash should be encoded by hex.")
		}
	})
}

func ExampleHash() {
	s := struct{ Test string }{Test: "test"}
	x := Hash(s)
	fmt.Println(x)
	//Output: e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746
}

func TestToBytes(t *testing.T) {
	s := "test"
	b := ToBytes(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("ToBytes should return a slice of bytes got %s", k)
	}
}

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	tests := []test{
		{input: "0:1:2", sep: ":", index: 1, output: "1"},
		{input: "0:1:2", sep: ":", index: 10, output: ""},
		{input: "0:1:2", sep: "/", index: 0, output: "0:1:2"},
	}
	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if tc.output != got {
			t.Errorf("Expected %s, got %s", tc.output, got)
		}
	}
}

func TestHandleErr(t *testing.T) {
	oldFn := logFn
	defer func() {
		logFn = oldFn
	}()
	called := false
	err := errors.New("Test")
	logFn = func(v ...interface{}) {
		called = true
	}
	HandleErr(err)
	if !called {
		t.Error("HandleErr function isn't called")
	}
}

func TestFromByte(t *testing.T) {
	type testStruct struct {
		Test string
	}
	var restored testStruct
	ts := testStruct{"Test"}
	b := ToBytes(ts)
	FromBytes(&restored, b)
	if !reflect.DeepEqual(ts, restored) {
		t.Errorf("FromBytes() should restore struct.")
	}
}

func TestToJSON(t *testing.T) {
	type testStruct struct {
		Test string
	}
	s := testStruct{"Test"}
	b := ToJSON(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("Expected %v got %v", reflect.Slice, k)
	}
	var restored testStruct
	json.Unmarshal(b, &restored)
	if !reflect.DeepEqual(s, restored) {
		t.Errorf("ToJSON should encode JSON correctly.")
	}
}
