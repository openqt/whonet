package utils

//
// 关于Go的json处理说明
// 参考 https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
//

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FileStruct struct {
	Length int64
	Path   []string
}

type InfoStruct struct {
	Files       []FileStruct `json:"files"`
	Length      int64        `json:"length"`
	Name        string       `json:"name"`
	PieceLength int64        `json:"piece length"`
	Pieces      Pieces       `json:"pieces"`
	RootHash    string       `json:"root hash"`
	Private     int          `json:"private"`
}

type TorrentStruct struct {
	Announce     string     `json:"announce"`
	AnnounceList [][]string `json:"announce-list"`
	Comment      string     `json:"comment"`
	CreatedBy    string     `json:"created by"`
	CreationDate int64      `json:"creation date"`
	Info         InfoStruct `json:"info"`
}

type Pieces struct {
	S string
}

// 显示pieces的数量及前三个SHA1
func (j Pieces) MarshalJSON() ([]byte, error) {
	const LEN = 40
	s := fmt.Sprintf("[%d]", len(j.S)/LEN)
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

// 数据结构转Torrent结构
func NewTorrent(data interface{}) *TorrentStruct {
	torrent := new(TorrentStruct)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // 不做字符转换

	err := enc.Encode(data)
	CheckError(err)

	json.Unmarshal(buf.Bytes(), torrent)

	return torrent
}
