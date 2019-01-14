package bencode

import (
	"github.com/openqt/whonet/utils"
	"testing"
)

type BenTestData struct {
	Encoded string
	Value   interface{}
}

func TestInt(t *testing.T) {
	data := map[string]int{
		"i42e":  42,
		"i-42e": -42,
		"i0e":   0,
	}

	enc, dec := NewEncoder(), NewDecoder()
	for s, v := range data {
		if !utils.DeepEqual(enc.Encode(v), s) {
			t.Errorf("Encode int(%v) %v != %v\n", v, s, enc.encodeInt(v))
		}

		if !utils.DeepEqual(dec.Decode([]byte(s)), v) {
			t.Errorf("Decode int(%v) %v != %v\n", v, s, dec.Decode([]byte(s)))
		}
	}
}

func TestString(t *testing.T) {
	data := map[string]string{
		"4:spam":  "spam",
		"0:":      "",
		"5:barbb": "barbb",
	}

	enc, dec := NewEncoder(), NewDecoder()
	for s, v := range data {
		if enc.Encode(v) != s {
			t.Errorf("String %s != %s\n", s, enc.Encode(v))
		}

		if !utils.DeepEqual(dec.Decode([]byte(s)), v) {
			t.Errorf("String %s != %s\n", dec.Decode([]byte(s)), v)
		}
	}
}

func TestList(t *testing.T) {
	data := []BenTestData{
		{"li1ei2ei3ee", [3]int{1, 2, 3}},
		{"li1ei2ei3ee", []int{1, 2, 3}},
		{"l4:spami42ee", []interface{}{"spam", 42}},
		{"l3:fool4:spam2:okei42ee", []interface{}{"foo", [2]string{"spam", "ok"}, 42}},
	}

	enc, dec := NewEncoder(), NewDecoder()
	for _, td := range data {
		if enc.Encode(td.Value) != td.Encoded {
			t.Errorf("List %s != %s\n", enc.Encode(td.Value), td.Encoded)
		}

		if !utils.DeepEqual(dec.Decode([]byte(td.Encoded)), td.Value) {
			t.Errorf("List %v != %v\n", dec.Decode([]byte(td.Encoded)), td.Value)
		}
	}
}

func TestDict(t *testing.T) {
	data := []BenTestData{
		{"d3:bar4:spam3:fooi42ee", map[string]interface{}{"bar": "spam", "foo": 42}},
	}

	enc, dec := NewEncoder(), NewDecoder()
	for _, td := range data {
		if enc.Encode(td.Value) != td.Encoded {
			t.Errorf("Map %s != %s\n", td.Encoded, enc.Encode(td.Value))
		}

		if !utils.DeepEqual(dec.Decode([]byte(td.Encoded)), td.Value) {
			t.Errorf("Map %s != %s\n", td.Encoded, enc.Encode(td.Value))
		}
	}
}
