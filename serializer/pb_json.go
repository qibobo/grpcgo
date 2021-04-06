package serializer

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func PB2Json(message proto.Message, filename string) error {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	jsonStr, err := marshaler.MarshalToString(message)
	if err != nil {
		return fmt.Errorf("failed to marshal proto to json %s\n", err)
	}
	err = ioutil.WriteFile(filename, []byte(jsonStr), 644)
	if err != nil {
		return fmt.Errorf("failed to write json string to file %s\n", err)
	}
	return nil
}
