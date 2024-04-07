package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hackeroneapiclient/pkg/csvfiles"
	"hackeroneapiclient/pkg/targetretrievalservice"
	"io"
	"os"
)

var output string

var scanCmd = &cobra.Command{
	Use:   "webapptargets",
	Short: "Searches the HackerOne API for web application targets",
	Long:  `Searches the HackerOne API for web application targets and outputs them to a file`,
	Run: func(cmd *cobra.Command, args []string) {
		o := csvfiles.OutputFile{Path: output}
		u := viper.GetString("hackeroneapiclient_username")
		t := viper.GetString("hackeroneapiclient_token")
		we := io.Writer(os.Stderr)
		wo := io.Writer(os.Stdout)
		targetretrievalservice.SearchForWebApps(&o, u, t, we, wo)
	},
}

func init() {
	scanCmd.Flags().StringVarP(&output, "output", "o", "", "Path to write an output .csv file")
	_ = scanCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(scanCmd)
}
