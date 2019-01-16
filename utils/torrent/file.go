package torrent

//
// 关于Go的json处理说明
// 参考 https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
//

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openqt/whonet/utils"
	"github.com/openqt/whonet/utils/bencode"
	"time"
)

type FileStruct struct {
	Length int64    `json:"length"`
	Path   []string `json:"path"`
	//Md5sum string `json:",omitempty"`
}

func (j FileStruct) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["length"] = j.Length
	result["path"] = j.Path
	return result
}

type InfoStruct struct {
	Pieces      Pieces `json:"pieces"`
	PieceLength int64  `json:"piece length"`
	Name        string `json:"name"`

	Private *int `json:"private,omitempty"`

	// Single file mode
	Length *int64 `json:"length,omitempty"`
	//Md5sum string `json:"md5sum,omitempty"`

	// Multiple file mode
	Files    []FileStruct `json:"files,omitempty"`
	RootHash *string      `json:"root hash,omitempty"`
}

func (j InfoStruct) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["pieces"] = j.Pieces.O
	result["piece length"] = j.PieceLength
	result["name"] = j.Name
	if j.Private != nil {
		result["private"] = *j.Private
	}
	if j.Length != nil {
		result["length"] = *j.Length
	}
	if j.Files != nil {
		var v []interface{}
		for _, i := range j.Files {
			v = append(v, i.ToMap())
		}
		result["files"] = v
	}
	if j.RootHash != nil {
		result["root hash"] = j.RootHash
	}

	return result
}

type TorrentStruct struct {
	Info         InfoStruct `json:"info"`
	Announce     string     `json:"announce"`
	AnnounceList [][]string `json:"announce-list,omitempty"`
	CreationDate *Timestamp `json:"creation date,omitempty"`
	CreatedBy    *string    `json:"created by,omitempty"`
	Comment      *string    `json:"comment,omitempty"`
	Encoding     *string    `json:"encoding,omitempty"`
}

func (j TorrentStruct) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["info"] = j.Info.ToMap()
	result["announce"] = j.Announce
	if j.AnnounceList != nil {
		result["announce-list"] = j.AnnounceList
	}
	if j.CreationDate != nil {
		result["creation date"] = j.CreationDate.Unix()
	}
	if j.CreatedBy != nil {
		result["created by"] = *j.CreatedBy
	}
	if j.Comment != nil {
		result["comment"] = *j.Comment
	}
	if j.Encoding != nil {
		result["encoding"] = *j.Encoding
	}

	return result
}

type Pieces struct {
	S string
	O string `json:"-"` // 二进制原始内容，避免转义
}

// 显示pieces的数量及前三个SHA1
func (j Pieces) MarshalJSON() ([]byte, error) {
	var s string
	const LEN = 40
	s = fmt.Sprintf("[%d]", len(j.S)/LEN)
	for i := 0; i < 3; i++ {
		n := i * LEN
		s += fmt.Sprintf(" %s", j.S[n:n+LEN])
	}
	s += " ..."
	return json.Marshal(s)
}

func (j *Pieces) UnmarshalJSON(data []byte) error {
	json.Unmarshal(data, &j.S)
	return nil
}

type Timestamp struct {
	time.Time
}

//func (j Timestamp) MarshalJSON() ([]byte, error) {
//	return json.Marshal(j.Time.Unix())
//}

// 转换成时间格式显示
func (j *Timestamp) UnmarshalJSON(data []byte) error {
	var t int64
	json.Unmarshal(data, &t)
	j.Time = time.Unix(t, 0)
	return nil
}

// 数据结构转Torrent结构
func NewTorrent(data []byte) *TorrentStruct {
	dec := bencode.NewDecoder()
	val := dec.Decode(data)

	torrent := new(TorrentStruct)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // 不做字符转换

	err := enc.Encode(val)
	utils.CheckError(err)

	json.Unmarshal(buf.Bytes(), torrent)
	torrent.Info.Pieces.O = dec.Pieces

	return torrent
}
