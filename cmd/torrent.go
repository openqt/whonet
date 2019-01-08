package cmd

import (
	"fmt"
	"github.com/openqt/whonet/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
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
	dat, err := ioutil.ReadFile(file)
	utils.Check(err)
	LOG.Debugf("Length: %d", len(dat))

	bc := utils.NewBencode()
	bc.Decode(string(dat))
	fmt.Println(bc.ToJson("  "))
}
