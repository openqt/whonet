package utils

import (
	"testing"
)

func TestInt(t *testing.T) {
	bc := Bencode{}

	data := map[string]int{
		"i42e":  42,
		"i-42e": -42,
		"i0e":   0,
	}

	for s, v := range data {
		if bc.EncodeInt(v) != s {
			t.Errorf("Int(%v) %v != %v\n", v, s, bc.EncodeInt(v))
		}
	}
}

func TestBytes(t *testing.T) {
	bc := Bencode{}

	data := map[string]string{
		"4:spam":  "spam",
		"0:":      "",
		"5:barbb": "barbb",
	}

	for s, v := range data {
		if bc.EncodeString(v) != s {
			t.Errorf("Bytes %s != %s\n", s, bc.EncodeString(v))
		}
	}
}

func TestList(t *testing.T) {
	bc := Bencode{}

	data := map[string][]interface{}{
		"l4:spami42ee":            {"spam", 42},
		"l3:fool4:spam2:okei42ee": {"foo", []string{"spam", "ok"}, 42},
	}

	for s, v := range data {
		if bc.EncodeList(v) != s {
			t.Errorf("List %s != %s\n", s, bc.EncodeList(v))
		}
	}
}

func TestDict(t *testing.T) {
	bc := Bencode{}

	data := map[string]map[string]interface{}{
		"d3:bar4:spam3:fooi42ee": {"bar": "spam", "foo": 42},
	}

	for s, v := range data {
		if bc.EncodeDict(v) != s {
			t.Errorf("List %s != %s\n", s, bc.EncodeDict(v))
		}
	}
}
