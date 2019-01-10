package utils

type FileStruct struct {
	Length int64
	Path   []string
}

type InfoStruct struct {
	Files       []FileStruct
	Length      int64
	Name        string
	PieceLength int64 `json:"piece length"`
	//Pieces       string  // TODO: hash list structure
}

type TorrentStruct struct {
	Announce     string
	AnnounceList [][]string `json:"announce-list"`
	Comment      string
	CreatedBy    string `json:"created by"`
	CreationDate int64  `json:"creation date"`
	Info         InfoStruct
}
