package serbench

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/samuel/go-thrift"
)

func BenchmarkDecodeThriftReflect(b *testing.B) {
	p := thrift.NewBinaryProtocol(true, false, 256)
	b.StopTimer()
	buf := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		if err := thrift.EncodeStruct(buf, p, testSt); err != nil {
			b.Fatal(err)
		}
	}
	b.StartTimer()
	s := &testStruct{}
	for i := 0; i < b.N; i++ {
		if err := thrift.DecodeStruct(buf, p, s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeThriftProtocol(b *testing.B) {
	p := thrift.NewBinaryProtocol(true, false, 256)
	b.StopTimer()
	buf := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		if err := thrift.EncodeStruct(buf, p, testSt); err != nil {
			b.Fatal(err)
		}
	}
	b.StartTimer()
	s := &testStruct{}
	for i := 0; i < b.N; i++ {
		// Begin struct
		if err := p.ReadStructBegin(buf); err != nil {
			b.Fatal(err)
		}
		for {
			ftype, fid, err := p.ReadFieldBegin(buf)
			if err != nil {
				b.Fatal(err)
			}
			if ftype == thrift.TypeStop {
				break
			}
			switch fid {
			case 1:
				if ftype != thrift.TypeString {
					b.Fatal("Wrong field type (1)")
				}
				s.String, err = p.ReadString(buf)
				if err != nil {
					b.Fatal(err)
				}
			case 2:
				if ftype != thrift.TypeI32 {
					b.Fatal("Wrong field type (2)")
				}
				s.Int, err = p.ReadI32(buf)
				if err != nil {
					b.Fatal(err)
				}
			case 3:
				if ftype != thrift.TypeList {
					b.Fatal("Wrong field type(3)")
				}
				etype, size, err := p.ReadListBegin(buf)
				if err != nil {
					b.Fatal(err)
				}
				if etype != thrift.TypeString {
					b.Fatal("Wrong list type (3)")
				}
				s.StringList = make([]string, size)
				for i := 0; i < size; i++ {
					s.StringList[i], err = p.ReadString(buf)
					if err != nil {
						b.Fatal(err)
					}
				}
				if err := p.ReadListEnd(buf); err != nil {
					b.Fatal(err)
				}
			default:
				b.Fatal("Unknown field id")
			}
			if err := p.ReadFieldEnd(buf); err != nil {
				b.Fatal(err)
			}
		}
		if err := p.ReadStructEnd(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeJson(b *testing.B) {
	b.StopTimer()
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		if err := e.Encode(testSt); err != nil {
			b.Fatal(err)
		}
	}
	b.StartTimer()
	s := &testStruct{}
	d := json.NewDecoder(buf)
	for i := 0; i < b.N; i++ {
		if err := d.Decode(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeGob(b *testing.B) {
	b.StopTimer()
	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		if err := e.Encode(testSt); err != nil {
			b.Fatal(err)
		}
	}
	b.StartTimer()
	s := &testStruct{}
	d := gob.NewDecoder(buf)
	for i := 0; i < b.N; i++ {
		if err := d.Decode(s); err != nil {
			b.Fatal(err)
		}
	}
}
