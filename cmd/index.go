package cmd

import (
	"github.com/spf13/cobra"
	"github.com/werkshy/likeness/index"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Build or update the index of your photo directory",
	Long:  `Walk through the directory, indexing files into the DB using md5sums.`,
	Run: func(cmd *cobra.Command, args []string) {
		index.StartIndex(mainDir, dbConnection)
	},
}

func init() {
	RootCmd.AddCommand(indexCmd)
}
