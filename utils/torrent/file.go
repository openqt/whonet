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
	"time"
)

type FileStruct struct {
	Length int64
	Path   []string
	Md5sum string `json:",omitempty"`
}

type InfoStruct struct {
	Pieces      Pieces `json:"pieces"`
	PieceLength int64  `json:"piece length"`
	Private     int    `json:",omitempty"`
	Name        string `json:"name"`

	// Single file mode
	Length int64  `json:"length,omitempty"`
	Md5sum string `json:",omitempty"`

	// Multiple file mode
	Files    []FileStruct `json:"files,omitempty"`
	RootHash string       `json:"root hash,omitempty"`
}

type TorrentStruct struct {
	Info         InfoStruct `json:"info"`
	Announce     string     `json:"announce"`
	AnnounceList [][]string `json:"announce-list,omitempty"`
	CreationDate Timestamp  `json:"creation date,omitempty"`
	CreatedBy    string     `json:"created by,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	Encoding     string     `json:",omitempty"`
}

type Pieces struct {
	string
}

// 显示pieces的数量及前三个SHA1
func (j Pieces) MarshalJSON() ([]byte, error) {
	const LEN = 40
	s := fmt.Sprintf("[%d]", len(j.string)/LEN)
	for i := 0; i < 3; i++ {
		n := i * LEN
		s += fmt.Sprintf(" %s", j.string[n:n+LEN])
	}
	s += " ..."
	return json.Marshal(s)
}

func (j *Pieces) UnmarshalJSON(data []byte) error {
	json.Unmarshal(data, &j.string)
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
func NewTorrent(data interface{}) *TorrentStruct {
	torrent := new(TorrentStruct)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // 不做字符转换

	err := enc.Encode(data)
	utils.CheckError(err)

	json.Unmarshal(buf.Bytes(), torrent)

	return torrent
}
