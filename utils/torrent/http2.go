package torrent

type GetStruct struct {
	InfoHash   string `json:"info_hash"`
	PeerId     string `json:"peer_id"`
	Port       int    `json:"port,omitempty"`
	Uploaded   int    `json:"uploaded,omitempty"`
	Downloaded int    `json:"downloaded,omitempty"`
	Left       int    `json:"left,omitempty"`
	Compact    int    `json:"compact,omitempty"`
	NoPeerId   int    `json:"no_peer_id,omitempty"`
	Event      string `json:"event,omitempty"`
	IP         string `json:"ip,omitempty"`
	NumWant    int    `json:"num_want,omitempty"`
	Key        string `json:"key,omitempty"`
	TrackerId  string `json:"trackerid,omitempty"`
}

func init() {
	//
}
