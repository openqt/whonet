package utils

import (
	"fmt"
	"sort"
)

/* Encoding algorithm
https://en.wikipedia.org/wiki/Bencode

Bencode uses ASCII characters as delimiters and digits.

- An integer is encoded as i<integer encoded in base ten ASCII>e.
Leading zeros are not allowed (although the number zero is still represented as "0").
Negative values are encoded by prefixing the number with a hyphen-minus.
The number 42 would thus be encoded as i42e, 0 as i0e, and -42 as i-42e.
Negative zero is not permitted.

- A byte string (a sequence of bytes,
not necessarily characters) is encoded as <length>:<contents>.
The length is encoded in base 10, like integers, but must be non-negative (zero is allowed);
the contents are just the bytes that make up the string.
The string "spam" would be encoded as 4:spam.
The specification does not deal with encoding of characters outside the ASCII set;
to mitigate this, some BitTorrent applications explicitly communicate the
encoding (most commonly UTF-8) in various non-standard ways.
This is identical to how netstrings work,
except that netstrings additionally append a comma suffix after the byte sequence.

- A list of values is encoded as l<contents>e.
The contents consist of the bencoded elements of the list, in order, concatenated.
A list consisting of the string "spam" and the number 42 would be encoded as: l4:spami42ee.
Note the absence of separators between elements,
and the first character is the letter 'l', not digit '1'.

- A dictionary is encoded as d<contents>e.
The elements of the dictionary are encoded each key immediately followed by its value.
All keys must be byte strings and must appear in lexicographical order.
A dictionary that associates the values 42 and "spam" with the keys "foo" and "bar",
respectively (in other words, {"bar": "spam", "foo": 42}),
would be encoded as follows: d3:bar4:spam3:fooi42ee.

There are no restrictions on what kind of values may be stored in lists and dictionaries;
they may (and usually do) contain other lists and dictionaries.
This allows for arbitrarily complex data structures to be encoded.

*/

type Bencode struct {
}

// 编码整数
func (code *Bencode) EncodeInt(i int) string {
	return fmt.Sprintf("i%de", i)
}

// 编码字符串
func (code *Bencode) EncodeString(s string) string {
	return fmt.Sprintf("%d:%s", len(s), string(s))
}

// 编码字符串数组
func (code *Bencode) EncodeStringList(ls []string) string {
	val := ""
	for _, i := range ls {
		val += code.EncodeString(i)
	}
	return "l" + val + "e"
}

// 编码整数数组
func (code *Bencode) EncodeIntList(ls []int) string {
	val := ""
	for _, i := range ls {
		val += code.EncodeInt(i)
	}
	return "l" + val + "e"
}

// 按类型编码
func (code *Bencode) typedValue(val interface{}) string {
	switch v := val.(type) {
	case int:
		return code.EncodeInt(v)
	case string:
		return code.EncodeString(v)
	case []string:
		return code.EncodeStringList(v)
	case []int:
		return code.EncodeIntList(v)
	case []interface{}:
		return code.EncodeList(v)
	default:
		fmt.Printf("Value %v (type %T) not recognized.", v, v)
	}
	return ""
}

// 编码一般数组
func (code *Bencode) EncodeList(ls []interface{}) string {
	val := ""
	for _, i := range ls {
		val += code.typedValue(i)
	}
	return "l" + val + "e"
}

// 编码一般字典
func (code *Bencode) EncodeDict(dt map[string]interface{}) string {
	val := ""

	var keys []string
	for k := range dt {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val += code.EncodeString(k) + code.typedValue(dt[k])
	}

	return "d" + val + "e"
}
