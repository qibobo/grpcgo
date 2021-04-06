package serializer

import (
	"fmt"
	"io/ioutil"

	"github.com/qibobo/grpcgo/models"
	"google.golang.org/protobuf/proto"
)

func PB2File(message proto.Message, filename string) error {

	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal proto message %s\n", err)
	}
	err = ioutil.WriteFile(filename, data, 644)
	if err != nil {
		return fmt.Errorf("failed to write proto data to file %s\n", err)
	}
	return nil
}

func File2PB(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read proto bin file %s\n", err)
	}
	p := &models.Person{}
	err = proto.Unmarshal(data, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal proto %s\n", err)
	}
	return nil
}
