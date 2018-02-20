package gogomarshal

import (
	"io"

	"errors"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Proto is a Marshaller which marshals/unmarshals into/from serialize proto bytes
type Proto struct{}

// ContentType always returns "application/octet-stream".
func (*Proto) ContentType() string {
	return "application/octet-stream"
}

// Marshal marshals "value" into Proto
func (*Proto) Marshal(value interface{}) ([]byte, error) {
	message, ok := value.(proto.Message)
	if !ok {
		return nil, errors.New("unable to marshal non proto field")
	}
	return proto.Marshal(message)
}

// Unmarshal unmarshals proto "data" into "value"
func (*Proto) Unmarshal(data []byte, value interface{}) error {
	message, ok := value.(proto.Message)
	if !ok {
		return errors.New("unable to unmarshal non proto field")
	}
	return proto.Unmarshal(data, message)
}

// NewDecoder returns a Decoder which reads proto stream from "reader".
func (marshaller *Proto) NewDecoder(reader io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(value interface{}) error {
		buffer, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		return marshaller.Unmarshal(buffer, value)
	})
}

// NewEncoder returns an Encoder which writes proto stream into "writer".
func (marshaller *Proto) NewEncoder(writer io.Writer) runtime.Encoder {
	return runtime.EncoderFunc(func(value interface{}) error {
		buffer, err := marshaller.Marshal(value)
		if err != nil {
			return err
		}
		_, err = writer.Write(buffer)
		if err != nil {
			return err
		}

		return nil
	})
}
