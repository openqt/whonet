package utils

type FileStruct struct {
	Length int64
	Path   []string
}

type InfoStruct struct {
	Files        []FileStruct
	Length       int64
	Name         string
	Piece_length int64
	Pieces       string
}

type TorrentStruct struct {
	Announce      string
	Announce_list []string
	Comment       string
	Created_by    int64
	Creation_date int64
	Info          InfoStruct
}
