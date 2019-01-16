package cmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/openqt/whonet/utils"
	"github.com/openqt/whonet/utils/bencode"
	"github.com/openqt/whonet/utils/torrent"
	"github.com/spf13/cobra"
	"io/ioutil"
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

	dec := bencode.NewDecoder()
	val := dec.Decode(bytes)
	torrent := torrent.NewTorrent(val)

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
		fmt.Printf("%02d: %x\n", i+1, hs)
		h.Reset()
	}
}
