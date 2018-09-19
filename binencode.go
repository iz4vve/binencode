// Package binencode contains utilities to encode a generic map[string]interface{}
// to a binary 8-byte aligned payload as described in the docs of `EncodePayload`
package binencode

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// encodePayload encodes in binary format the payload for the RabbitMQ message.
// The encoding is structured as a LittleEndian representation of:
//
// - number of fields (int32)
// [for each key-value pair]:
// - length of key (int32)
// - key (8-byte aligned)
// - length of value (int32)
// - value (8-byte aligned)
//
// TODO: insert link to the wiki explaining the binary encoding.
//
// Currently only strings, slices and 4-byte integers are supported as values,
// and only strings as keys.
func EncodePayload(fields map[string]interface{}) ([]byte, error) {

	b := new(bytes.Buffer)
	offset := 4

	binary.Write(b, binary.LittleEndian, int32(len(fields)))

	for key, value := range fields {

		k := []byte(key)
		binary.Write(b, binary.LittleEndian, int32(len(k)))
		for _, c := range k {
			binary.Write(b, binary.LittleEndian, c)
		}

		currLen := b.Len()
		offset = (offset+len(key)+16)&^7 - 4

		paddingKey := make([]byte, offset-currLen)
		binary.Write(b, binary.LittleEndian, paddingKey)

		switch t := value.(type) {
		case string:
			binary.Write(b, binary.LittleEndian, int32(len(t)))
			offset = (offset+len(t)+16)&^7 - 4
			v := []byte(t)
			for _, c := range v {
				binary.Write(b, binary.LittleEndian, c)
			}
		case []byte:
			binary.Write(b, binary.LittleEndian, int32(len(t)))
			binary.Write(b, binary.LittleEndian, t)
			offset = (offset+len(t)+16)&^7 - 4
		default:
			return []byte{}, errors.New("invalid payload")
		}

		currLen = b.Len()

		paddingValue := make([]byte, offset-currLen)
		binary.Write(b, binary.LittleEndian, paddingValue)
	}

	return b.Bytes(), nil
}

// decodePayload decodes a binary payload encoded by encodePayload into
// a map of strings to bytes.
// Currently only strings and slices of bytes are supported in the encoding function,
// but if the encoded payload is generated by an external source, the value can contain
// other objects.
func DecodePayload(payload []byte) (map[string][]byte, error) {
	var m = make(map[string][]byte)

	offset := 4

	numFields := binary.LittleEndian.Uint32(payload[:4])

	for i := 0; i < int(numFields); i++ {
		keyLen := binary.LittleEndian.Uint32(payload[offset : offset+4])
		key := payload[offset+4 : offset+4+int(keyLen)]
		offset = (offset+16+int(keyLen))&^7 - 4

		valueLen := binary.LittleEndian.Uint32(payload[offset : offset+4])
		value := payload[offset+4 : offset+4+int(valueLen)]
		offset = (offset+16+int(valueLen))&^7 - 4

		m[string(key)] = value
	}

	return m, nil
}
