package torrent

import (
	"github.com/openqt/whonet/utils"
	"testing"
	"io/ioutil"
	"github.com/openqt/whonet/utils/bencode"
)

func TestFile(t *testing.T) {
	data := []string {
		"../../tests/puppy.torrent",
		"../../tests/ubuntu-18.10-desktop-amd64.iso.torrent",
		"../../tests/CentOS-7-x86_64-Minimal-1810.torrent",
	}

	enc := bencode.NewEncoder()
	for _, name := range data {
		bs, err := ioutil.ReadFile(name)
		utils.CheckError(err)
		torrent := NewTorrent(bs)

		s := enc.Encode(torrent.ToMap())
		if s != string(bs) {
			t.Errorf("File %s is after encode then decode.\n", name)
		}
	}
}
