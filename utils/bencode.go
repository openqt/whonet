package utils

// 更新依赖库的环境参数
// HTTPS_PROXY=socks5://127.0.0.1:1080 go get -u -v github.com/sirupsen/logrus

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
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

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.SetReportCaller(true)
}
func NewBencode() *Bencode {
	p := new(Bencode)
	p.idx = 0
	return p
}

// 编码为文本
func (code *Bencode) Encode(val interface{}) string {
	code.Data = val
	code.Code = code.encode(val)
	code.idx = 0
	return code.Code
}

func (code *Bencode) encode(val interface{}) string {
	switch v := val.(type) {
	case int:
		return code.encodeInt(v)
	case string:
		return code.encodeString(v)
	case []string:
		return code.encodeStringList(v)
	case []int:
		return code.encodeIntList(v)
	case []interface{}:
		return code.encodeList(v)
	default:
		fmt.Printf("Value %v (type %T) not recognized.", v, v)
	}
	return ""
}

// 从文本解码
func (code *Bencode) Decode(s string) interface{} {
	code.Code = s
	code.idx = 0
	return code.decode()
}

func (code *Bencode) decode() interface{} {
	c := code.Code[code.idx]
	switch {
	case c == 'i':
		return code.decodeInt()
	case c == 'l':
		return code.DecodeList()
	case c == 'd':
		return code.DecodeDict()
	case '0' <= c && c <= '9':
		return code.decodeString()
	case c == 'e':
		code.idx += 1
	default:
		log.Errorf("%s is invalid.", code.Code)
	}
	return nil
}

// 编码整数
func (code *Bencode) encodeInt(i int) string {
	val := fmt.Sprintf("i%de", i)
	log.Debugf("%val => %val", i, val)
	return val
}

// 解码整数
func (code *Bencode) decodeInt() int {
	code.idx += 1 // i<>e

	var val int
	for i := code.idx; i < len(code.Code); i++ {
		c := code.Code[i]
		if c == 'e' {
			s := code.Code[code.idx:i]
			val, _ = strconv.Atoi(s)
			log.Debugf("%v => %v", s, val)
			code.idx += i // 指向e下一个字符
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
	for n = 0; n < len(code.Code); n++ {
		c := code.Code[n]
		if c < '0' || c > '9' {
			length, _ = strconv.Atoi(code.Code[:n])
			break
		}
	}

	if length < 0 {
		log.Errorf("%s cannot be decoded as string.", code.Code)
		return ""
	}

	code.idx += n + 1 + length // 指向下一个字符
	v := code.Code[n+1 : n+1+length]
	log.Debugf("%s", v)
	return v
}

// 编码字符串数组
func (code *Bencode) encodeStringList(ls []string) string {
	val := ""
	for _, i := range ls {
		val += code.encodeString(i)
	}
	return "l" + val + "e"
}

// 编码整数数组
func (code *Bencode) encodeIntList(ls []int) string {
	val := ""
	for _, i := range ls {
		val += code.encodeInt(i)
	}
	return "l" + val + "e"
}

// 编码一般数组
func (code *Bencode) encodeList(ls []interface{}) string {
	val := ""
	for _, i := range ls {
		val += code.encode(i)
	}
	return "l" + val + "e"
}

func (code *Bencode) DecodeList() []interface{} {
	code.idx += 1

	var val []interface{}
	for code.idx < len(code.Code) {
		v := code.Decode(code.Code[code.idx:])
		if v != nil {
			val = append(val, v)
		} else {
			break
		}
	}
	log.Debugf("%v", val)
	return val
}

// 编码一般字典
func (code *Bencode) EncodeDict(dt map[string]interface{}) string {
	val := ""

	var keys []string
	for k := range dt {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//for _, k := range keys {
	//    val += code.encodeString(k) + code.encodeValue(dt[k])
	//}

	return "d" + val + "e"
}

func (code *Bencode) DecodeDict() map[string]interface{} {
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
