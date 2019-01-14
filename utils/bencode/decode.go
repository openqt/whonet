package bencode

import (
	"fmt"
	"strconv"
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

type Decoder struct {
	buf []byte
	idx int // (内部)解码字段长度
}

//////////////////////////////////////////////////////////////////////////////////////////
//
//  编解码函数
//
//////////////////////////////////////////////////////////////////////////////////////////
// 创建一个新的编码对象
func NewDecoder() *Decoder {
	return &Decoder{}
}

// 从文本解码
func (dec *Decoder) Decode(buf []byte) interface{} {
	dec.buf = buf
	dec.idx = 0
	return dec.decode()
}

func (dec *Decoder) decode() interface{} {
	var val interface{}
	if dec.IsEnd() {
		return val
	}

	c := dec.Byte(0)
	switch {
	case c == 'i':
		val = dec.decodeInt()
	case c == 'l':
		val = dec.decodeList()
	case c == 'd':
		val = dec.decodeDict()
	case '0' <= c && c <= '9':
		val = dec.decodeString()
	default:
		panic(fmt.Sprintf("%s is invalid.", dec.buf))
	}

	return val
}

//////////////////////////////////////////////////////////////////////////////////////////
//
//  内部数据封装函数
//
//////////////////////////////////////////////////////////////////////////////////////////
// 当前字符
func (dec *Decoder) Byte(pos int) byte {
	if pos > 0 {
		return dec.buf[pos]
	}
	return dec.buf[dec.Pos()]
}

// 编码索引之后的字符串
func (dec *Decoder) String(length int) string {
	if length > dec.Pos() { // 如果长度大于当前索引
		pos := dec.Pos()
		dec.idx = length // 更新索引
		return string(dec.buf[pos:length])
	}
	if length > 0 {
		return string(dec.buf[length:dec.Pos()])
	}
	return string(dec.buf[dec.Pos():]) // 给出当前索引之后的字符串
}

// 编码索引当前位置
func (dec *Decoder) Pos() int {
	return dec.idx
}

// 编码索引增一
func (dec *Decoder) Next() {
	dec.idx += 1
}

// 编码的长度
func (dec *Decoder) Len() int {
	return len(dec.buf)
}

// 重置索引为0
func (dec *Decoder) Reset() {
	dec.idx = 0
}

// 编码结束标志
func (dec *Decoder) IsEnd() bool {
	return dec.Byte(0) == 'e'
}

//////////////////////////////////////////////////////////////////////////////////////////
//
//  编码函数
//
//////////////////////////////////////////////////////////////////////////////////////////

// 解码整数
func (dec *Decoder) decodeInt() int {
	if dec.Byte(0) != 'i' {
		panic(fmt.Sprintf("%s is not an int.", dec.String(0)))
	}
	dec.Next() // i<>e

	val, pos := 0, dec.Pos()
	for !dec.IsEnd() {
		dec.Next()
	}

	s := dec.String(pos)
	val, _ = strconv.Atoi(s)
	dec.Next() // 指向e下一个字符
	return val
}

// 解码字符串
func (dec *Decoder) decodeString() string {
	var length, n int
	for n = dec.Pos(); n < dec.Len(); n++ {
		c := dec.Byte(n)
		if c < '0' || c > '9' {
			length, _ = strconv.Atoi(dec.String(n))
			dec.Next()
			break
		}
	}

	if length < 0 {
		panic(fmt.Sprintf("%s cannot be decoded as string.", dec.String(0)))
	}

	val := dec.String(n + 1 + length)
	return val
}

// 解码数组
func (dec *Decoder) decodeList() []interface{} {
	if dec.Byte(0) != 'l' {
		panic(fmt.Sprintf("%s is not a list.", dec.String(0)))
	}
	dec.Next()

	var val []interface{}
	for !dec.IsEnd() {
		val = append(val, dec.decode())
	}
	dec.Next() // 指向e下一个字符
	return val
}

// 解码字典
func (dec *Decoder) decodeDict() map[string]interface{} {
	if dec.Byte(0) != 'd' {
		panic(fmt.Sprintf("%s is not a map.", dec.String(0)))
	}
	dec.Next()

	val := make(map[string]interface{})
	for !dec.IsEnd() {
		key := dec.decodeString()
		_val := dec.decode()
		// All strings must be UTF-8 encoded, except for pieces, which contains binary data.
		if key == "pieces" {
			val[key] = fmt.Sprintf("%x", _val)
		} else {
			val[key] = _val
		}
	}
	dec.Next() // 指向e下一个字符
	return val
}
