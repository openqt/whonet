package utils

import (
	"testing"
)

func TestInt(t *testing.T) {
	data := map[string]int{
		"i42e":  42,
		"i-42e": -42,
		"i0e":   0,
	}

	bc := NewBencode()
	for s, v := range data {
		if bc.Encode(v) != s {
			t.Errorf("Encode int(%v) %v != %v\n", v, s, bc.encodeInt(v))
		}

		if bc.Decode(s) != v {
			t.Errorf("Decode int(%v) %v != %v\n", v, s, bc.Decode(s))
		}
	}
}

func TestString(t *testing.T) {
	data := map[string]string{
		"4:spam":  "spam",
		"0:":      "",
		"5:barbb": "barbb",
	}

	bc := NewBencode()
	for s, v := range data {
		if bc.Encode(v) != s {
			t.Errorf("String %s != %s\n", s, bc.Encode(v))
		}

		if bc.Decode(s) != v {
			t.Errorf("String %s != %s\n", bc.Decode(s), v)
		}
	}
}

func TestList(t *testing.T) {
	data := map[string][]interface{}{
		"li1ei2ei3ee":             {1, 2, 3},
		"l4:spami42ee":            {"spam", 42},
		"l3:fool4:spam2:okei42ee": {"foo", []string{"spam", "ok"}, 42},
	}

	bc := NewBencode()
	for s, v := range data {
		if bc.Encode(v) != s {
			t.Errorf("List %s != %s\n", s, bc.Encode(v))
		}

		//if reflect.DeepEqual(bc.Decode(s), v) {
		//	t.Errorf("List %v != %v\n", bc.Decode(s), v)
		//}
	}
}

//func TestDict(t *testing.T) {
//	bc := Bencode{}
//
//	data := map[string]map[string]interface{}{
//		"d3:bar4:spam3:fooi42ee": {"bar": "spam", "foo": 42},
//	}
//
//	for s, v := range data {
//		if bc.EncodeDict(v) != s {
//			t.Errorf("List %s != %s\n", s, bc.EncodeDict(v))
//		}
//	}
//}
