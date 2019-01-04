package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
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
	Data interface{} // 解码结构
	Code string      // 编码文本
	idx  int         // (内部)解码字段长度
}

// 全局设置
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.SetReportCaller(true)
}

// 创建一个新的编码对象
func NewBencode() *Bencode {
	p := new(Bencode)
	p.reset()
	return p
}

// 编码为文本
func (code *Bencode) Encode(val interface{}) string {
	code.Data = val
	code.Code = code.encode(reflect.ValueOf(val))
	code.reset()
	return code.Code
}

func (code *Bencode) encode(val reflect.Value) string {
	log.Debugf("Encode %v (%v)", val, val.Kind())
	switch val.Kind() {
	case reflect.Int:
		return code.encodeInt(int(val.Int()))
	case reflect.String:
		return code.encodeString(val.String())
	case reflect.Slice, reflect.Array:
		return code.encodeList(val)
	case reflect.Map:
		return code.encodeDict(val)
	}
	panic(fmt.Sprintf("Value %v (type %T) not recognized.", val, val))
}

// 从文本解码
func (code *Bencode) Decode(s string) interface{} {
	code.Code = s
	code.reset()
	return code.decode()
}

func (code *Bencode) decode() interface{} {
	c := code.currentByte(0)
	switch {
	case c == 'i':
		return code.decodeInt()
	case c == 'l':
		return code.decodeList()
	case c == 'd':
		return code.decodeDict()
	case '0' <= c && c <= '9':
		return code.decodeString()
	case c == 'e': // 目前仅忽略，不做一致性校验
		code.next()
		return nil
	default:
		panic(fmt.Sprintf("%s is invalid.", code.Code))
	}
}

func (code *Bencode) currentByte(pos int) byte {
	if pos > 0 {
		return code.Code[pos]
	}
	return code.Code[code.current()]
}

// 编码索引之后的字符串
func (code *Bencode) currentString(length int) string {
	if length > code.current() { // 如果长度大于当前索引
		pos := code.current()
		code.idx = length // 更新索引
		return code.Code[pos:length]
	}
	return code.Code[code.current():] // 给出当前索引之后的字符串
}

// 编码索引当前位置
func (code *Bencode) current() int {
	return code.idx
}

// 编码索引增一
func (code *Bencode) next() {
	code.idx += 1
}

// 编码的长度
func (code *Bencode) length() int {
	return len(code.Code)
}

// 重置索引为0
func (code *Bencode) reset() {
	code.idx = 0
}

// 编码整数
func (code *Bencode) encodeInt(i int) string {
	val := fmt.Sprintf("i%de", i)
	log.Debugf("%v => %v", i, val)
	return val
}

// 解码整数
func (code *Bencode) decodeInt() int {
	if code.currentByte(0) != 'i' {
		panic(fmt.Sprintf("%s is not an int.", code.currentString(0)))
	}
	code.next() // i<>e

	var val int
	for n := code.current(); n < code.length(); n++ {
		c := code.currentByte(n)
		if c == 'e' {
			s := code.currentString(n)
			val, _ = strconv.Atoi(s)
			log.Debugf("%v => %v", s, val)
			code.next() // 指向e下一个字符
			break
		}
	}

	return val
}

// 编码字符串
func (code *Bencode) encodeString(s string) string {
	return fmt.Sprintf("%d:%s", len(s), string(s))
}

// 解码字符串
func (code *Bencode) decodeString() string {
	var length, n int
	for n = code.current(); n < code.length(); n++ {
		c := code.currentByte(n)
		if c < '0' || c > '9' {
			length, _ = strconv.Atoi(code.currentString(n))
			code.next()
			break
		}
	}

	if length < 0 {
		panic(fmt.Sprintf("%s cannot be decoded as string.", code.currentString(0)))
	}

	v := code.currentString(n + 1 + length)
	log.Debugf("%s", v)
	return v
}

// 编码一般数组
func (code *Bencode) encodeList(l reflect.Value) string {
	val := ""
	for i := 0; i < l.Len(); i++ {
		v := l.Index(i)
		switch v.Kind() {
		case reflect.Interface:
			val += code.encode(v.Elem())
		default:
			val += code.encode(v)
		}
	}
	return "l" + val + "e"
}

// 解码数组
func (code *Bencode) decodeList() []interface{} {
	if code.currentByte(0) != 'l' {
		panic(fmt.Sprintf("%s is not a list.", code.currentString(0)))
	}
	code.next()

	var val []interface{}
	for code.current() < code.length() {
		if v := code.decode(); v != nil {
			val = append(val, v)
		} else {
			break
		}
	}
	log.Debugf("%v", val)
	return val
}

// 编码一般字典
func (code *Bencode) encodeDict(d reflect.Value) string {
	val := ""
	for _, key := range SortKeys(d.MapKeys()) {
		v := d.MapIndex(reflect.ValueOf(key)).Elem()
		val += code.encodeString(key) + code.encode(v)
	}
	return "d" + val + "e"
}

// 解码字典
func (code *Bencode) decodeDict() map[string]interface{} {
	//for n, c := range s {
	//    switch c {
	//    case 'i':
	//        ns := code.decodeInt()
	//        fmt.Println(ns)
	//    case 'l':
	//        //code.decodeStringList(s)
	//    case 'd':
	//        //
	//    default:
	//        code.decodeString()
	//    }
	//}
	return nil
}
