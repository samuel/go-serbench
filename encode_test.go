package serbench

import (
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/samuel/go-thrift"
)

type nullWriter int

func (n nullWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

type testStruct struct {
	String     string   `thrift:"1,required"`
	Int        int32    `thrift:"2,required"`
	StringList []string `thrift:"3,required"`
}

var testSt = &testStruct{
	String:     "string",
	Int:        123,
	StringList: []string{"foo", "bar"},
}

func BenchmarkEncodeThriftReflect(b *testing.B) {
	s := testSt
	w := nullWriter(0)
	p := thrift.NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		if err := thrift.EncodeStruct(w, p, s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeThriftProtocol(b *testing.B) {
	s := testSt
	w := nullWriter(0)
	p := thrift.NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		// Begin struct
		if err := p.WriteStructBegin(w, ""); err != nil {
			b.Fatal(err)
		}
		// Field 1
		if err := p.WriteFieldBegin(w, "", thrift.TypeString, 1); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteString(w, s.String); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteFieldEnd(w); err != nil {
			b.Fatal(err)
		}
		// Field 2
		if err := p.WriteFieldBegin(w, "", thrift.TypeI32, 2); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteI32(w, s.Int); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteFieldEnd(w); err != nil {
			b.Fatal(err)
		}
		// Field 3
		if err := p.WriteFieldBegin(w, "", thrift.TypeList, 3); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteListBegin(w, thrift.TypeString, len(s.StringList)); err != nil {
			b.Fatal(err)
		}
		for _, st := range s.StringList {
			if err := p.WriteString(w, st); err != nil {
				b.Fatal(err)
			}
		}
		if err := p.WriteListEnd(w); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteFieldEnd(w); err != nil {
			b.Fatal(err)
		}
		if err := p.WriteFieldStop(w); err != nil {
			b.Fatal(err)
		}
		// End struct
		if err := p.WriteStructEnd(w); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeJson(b *testing.B) {
	s := testSt
	w := nullWriter(0)
	e := json.NewEncoder(w)
	for i := 0; i < b.N; i++ {
		if err := e.Encode(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeJsonMarshal(b *testing.B) {
	s := testSt
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeGob(b *testing.B) {
	s := testSt
	w := nullWriter(0)
	e := gob.NewEncoder(w)
	for i := 0; i < b.N; i++ {
		if err := e.Encode(s); err != nil {
			b.Fatal(err)
		}
	}
}
