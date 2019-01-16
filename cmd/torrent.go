package cmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/openqt/whonet/utils"
	"github.com/openqt/whonet/utils/bencode"
	"github.com/openqt/whonet/utils/torrent"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"os"
)

var (
	Filename string // Torrent文件路径
	LOG      = utils.GetLogger()
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show torrent file information",
	Long:  `Parsing torrent file and list data in JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			fmt.Println(">>>", file)
			ShowTorrent(file)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func ShowTorrent(file string) {
	bytes, err := ioutil.ReadFile(file)
	utils.CheckError(err)
	LOG.Debugf("Length: %d", len(bytes))

	torrent := torrent.NewTorrent(bytes)

	b, err := json.MarshalIndent(torrent, "", "  ")
	fmt.Println(string(b))

	f, err := os.Open("/data/Go/src/github.com/openqt/whonet/tests/ubuntu-18.10-desktop-amd64.iso")
	utils.CheckError(err)
	bytes = make([]byte, torrent.Info.PieceLength)

	h := sha1.New()
	for i := 0; i < 10; i++ {
		_, err := f.Read(bytes)
		utils.CheckError(err)

		h.Write(bytes)
		hs := h.Sum(nil)
		fmt.Printf("%02d: %X\n", i+1, hs)
		h.Reset()
	}

	s := ""
	enc := bencode.NewEncoder()

	s = enc.Encode(torrent.Info.ToMap())
	h.Write([]byte(s))
	s = string(h.Sum(nil))
	fmt.Printf("Info SHA1: %X, %s\n", s, url.QueryEscape(s))
	h.Reset()

	h.Write(uuid.NewV4().Bytes())
	fmt.Printf("Peer ID: %s\n", url.QueryEscape(string(h.Sum(nil))))
	h.Reset()

	s = enc.Encode(torrent.ToMap())
	err = ioutil.WriteFile("t.to", []byte(s), 0644)
	utils.CheckError(err)
}
