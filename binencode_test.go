package binencode

import "testing"

func TestEncoding(t *testing.T) {
	test := map[string]interface{}{
		"test":            "I'm a string",
		"number":          "12",
		"asdasdasd":       []byte("asdkjfdbalfdj"),
		"yet another key": []byte("asdkjfdbaasasdasdfagagrgakambo;b   abnnlib&iafr93hlfdj"),
	}

	encoded, _ := EncodePayload(test)

	if (len(encoded)-len(test)*4-4)%8 != 0 {
		t.Errorf("expected 8bit padded payload, got %d bytes", len(encoded))
	}
}

func TestEncodingUnsupported(t *testing.T) {
	test := map[string]interface{}{
		"test": 12,
	}

	encoded, err := EncodePayload(test)

	if err == nil {
		t.Error("expected invalid payload error, got nil")
	}
	if len(encoded) != 0 {
		t.Error("Expected empty byte array")
	}
}

func TestDecoding(t *testing.T) {
	test := map[string]interface{}{
		"test":            "I'm a string",
		"number":          "12",
		"stuff":           []byte("whatever"),
		"yet another key": "some more data that have to encoded so that we can test the whole thing",
	}

	encoded, _ := EncodePayload(test)

	decoded, _ := DecodePayload(encoded)

	if len(decoded) == 0 {
		t.Error("expected success, got 0 fields")
	}

	for k, v := range decoded {
		value := test[k]

		switch typ := value.(type) {
		case string:
			if typ != string(v) {
				t.Errorf("expected %s, got %s", test[k], string(v))
			}
		case []byte:

			for i, b := range typ {
				if b != v[i] {
					t.Errorf("bytes %d differ: expected %b, got %b", i, b, v[i])
				}
			}
		}
	}
}
