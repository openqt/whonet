package bencode

import (
	"fmt"
	"reflect"
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
not necessarily characters) is encoded as <Len>:<contents>.
The Len is encoded in base 10, like integers, but must be non-negative (zero is allowed);
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

type Encoder struct {
}

//////////////////////////////////////////////////////////////////////////////////////////
//
//  编解码函数
//
//////////////////////////////////////////////////////////////////////////////////////////
// 创建一个新的编码对象
func NewEncoder() *Encoder {
	return &Encoder{}
}

// 编码为文本
func (enc *Encoder) Encode(val interface{}) string {
	return enc.encode(reflect.ValueOf(val))
}

func (enc *Encoder) encode(val reflect.Value) string {
	var result string
	switch val.Kind() {
	case reflect.Int:
		result = enc.encodeInt(int(val.Int()))
	case reflect.String:
		result = enc.encodeString(val.String())
	case reflect.Slice, reflect.Array:
		result = enc.encodeList(val)
	case reflect.Map:
		result = enc.encodeDict(val)
	default:
		panic(fmt.Sprintf("Value %v (type %T) not recognized.", val, val))
	}
	return result
}

//////////////////////////////////////////////////////////////////////////////////////////
//
//  编码函数
//
//////////////////////////////////////////////////////////////////////////////////////////
// 编码整数
func (enc *Encoder) encodeInt(i int) string {
	val := fmt.Sprintf("i%de", i)
	return val
}

// 编码字符串
func (enc *Encoder) encodeString(s string) string {
	val := fmt.Sprintf("%d:%s", len(s), string(s))
	return val
}

// 编码一般数组
func (enc *Encoder) encodeList(l reflect.Value) string {
	val := ""
	for i := 0; i < l.Len(); i++ {
		v := l.Index(i)
		switch v.Kind() {
		case reflect.Interface:
			val += enc.encode(v.Elem())
		default:
			val += enc.encode(v)
		}
	}
	val = "l" + val + "e"
	return val
}

// 编码一般字典
func (enc *Encoder) encodeDict(d reflect.Value) string {
	val := ""
	for _, key := range SortKeys(d.MapKeys()) {
		v := d.MapIndex(reflect.ValueOf(key)).Elem()
		val += enc.encodeString(key) + enc.encode(v)
	}
	val = "d" + val + "e"
	return val
}

///

// 对字典的key排序
func SortKeys(l []reflect.Value) []string {
	var keys []string
	for _, v := range l {
		keys = append(keys, v.String())
	}
	sort.Strings(keys)
	return keys
}
