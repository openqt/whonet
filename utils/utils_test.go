package utils

import (
	"testing"
)

func TestDeepEqual(t *testing.T) {
	data := [][]interface{}{
		{map[string]string{"a": "b"}, map[string]string{"a": "b"}, true},
		{map[string]string{"a": "b"}, map[string]string{"b": "b"}, false},
		{map[string]string{}, map[string]string{"b": "b"}, false},
		//
		{[]int{1, 2, 3}, [3]int{1, 2, 3}, true},
		{[]int{1, 2, 3}, []interface{}{[3]int{1, 2, 3}}, false},
		{[]int{2, 3}, []interface{}{[3]int{1, 2, 3}}, false},

		{[]interface{}{}, []interface{}{}, true},
		{nil, nil, true},
	}

	for _, td := range data {
		if DeepEqual(td[0], td[1]) != td[2] {
			t.Errorf("%v != %v\n", td[0], td[1])
		}
	}
}
