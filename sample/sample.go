package sample

import (
	"fmt"
	"os"

	"github.com/gogo/gateway"
	"github.com/golang/protobuf/proto"
)

var Marshaler = gateway.JSONPb{}

// ProtoToJson returns the object as the Dispatcher Server does through gRPC as a Gateway
// the v proto.Message is critical, I have not found out yet why but that's some sugar that
// somehow changes the way the marshaler expects the message to be, else the output is different
func protoToJson(v proto.Message, pretty bool) ([]byte, error) {
	if pretty {
		Marshaler.Indent = "  "
	}
	Marshaler.EmitDefaults = false
	out, err := Marshaler.Marshal(v)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func WriteToSample(b []byte, name string, f *os.File) error {
	if _, err := f.WriteString(fmt.Sprintf("\n// tag::%s[]\n", name)); err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	if _, err := f.WriteString(fmt.Sprintf("\n// end::%s[]\n", name)); err != nil {
		return err
	}
	return nil
}

func CreateSample(v proto.Message, name string, f *os.File, pretty bool) error {
	out, err := protoToJson(v, pretty)
	if err != nil {
		return err
	}

	err = WriteToSample(out, name, f)
	if err != nil {
		return err
	}

	return nil
}

func CreateJsonSample(v proto.Message, name string, f *os.File) error {
	return CreateSample(v, name, f, true)
}

func CreateJsonSampleRaw(v proto.Message, name string, f *os.File) error {
	return CreateSample(v, name, f, false)
}

func Setup(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}
