package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/spf13/cobra"
)

type FileRename struct {
	From string `csv:"Ur-Pfad+Ur-Datei"`
	To   string `csv:"Zielpfad+Zieldatei"`
}

type RenameErr struct {
	Fr  FileRename
	err error
}

const CountTxt = `counting files to rename: `
const RenFmt = "renaming %d/%d (errors: %d): %s"

var rootCmd = &cobra.Command{
	Use:     "batch-rename",
	Short:   "renames files in batch",
	Long:    "this command will read from a given *.csv with ; as seperator with 'Ur-Pfad+Ur-Datei;Zielpfad+Zieldatei' as header and rename all files listed in the csv\nif there are errors it will list all files that had errors",
	Args:    cobra.ExactArgs(1),
	Example: "batch-rename some.csv",
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		var renames []FileRename
		err := csvtag.LoadFromPath(path, &renames, csvtag.CsvOptions{Separator: ';'})
		if err != nil {
			fmt.Printf("error reading csv file %q: %s\n", path, err)
		}

		var rens []FileRename
		fmt.Print(CountTxt)
		count := 0
		maxLen := 0
		fmt.Printf("%d", count)
		for _, fr := range renames {
			fmt.Print(deleters(fmt.Sprintf("%d", count)))
			_, err := os.Stat(fr.From)
			if err == nil {
				count++
				if len(fr.From) > maxLen {
					maxLen = len(fr.From)
				}
				rens = append(rens, fr)
			}

			fmt.Printf("%d", count)
		}
		fmt.Println("")

		fmt.Printf("found %d/%d files\n", count, len(renames))
		fmt.Printf("start renaming ? (y|N): ")
		var y string
		_, err = fmt.Scanln(&y)
		if err != nil || y != "y" {
			fmt.Println("aborting")
			os.Exit(0)
		}

		count = 0
		var errs []RenameErr
		last := ""
		fmt.Printf(RenFmt, count, len(rens), len(errs), pad(last, maxLen))
		for _, ren := range rens {
			fmt.Print(deleters(fmt.Sprintf(RenFmt, count, len(rens), len(errs), pad(last, maxLen))))
			count++
			last = ren.From

			err := os.MkdirAll(filepath.Dir(ren.To), os.ModePerm)
			if err != nil {
				errs = append(errs, RenameErr{Fr: ren, err: err})
				fmt.Printf(RenFmt, count, len(rens), len(errs), pad(last, maxLen))
				continue
			}

			err = os.Rename(ren.From, ren.To)
			if err != nil {
				errs = append(errs, RenameErr{Fr: ren, err: err})
			}

			fmt.Printf(RenFmt, count, len(rens), len(errs), pad(last, maxLen))
		}
		fmt.Println()

		for _, rerr := range errs {
			fmt.Printf("%s\n", rerr.err)
		}

		fmt.Println("done")
	},
}

func deleters(text string) string {
	res := ""

	for range text {
		res += "\b"
	}

	return res
}

func pad(text string, l int) string {
	res := text
	for i := len(text); i < l; i++ {
		res += " "
	}

	return res
}

// Execute is the entry to this command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error running command: %s\n", err)
		os.Exit(1)
	}
}
