package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/paradoxgery/batch-rename/cmd/copy"
	"github.com/paradoxgery/batch-rename/cmd/rename"
)

var rootCmd = &cobra.Command{
	Use:   "batch-rename",
	Short: "offers batch utility for files",
}

// Execute is the entry to this command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error running command: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print more error messages")
	rootCmd.PersistentFlags().StringP("seperator", "s", ";", "csv seperator")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("seperator", rootCmd.PersistentFlags().Lookup("seperator"))

	rootCmd.AddCommand(rename.RenameCmd)
	rootCmd.AddCommand(copy.CopyCmd)
}
