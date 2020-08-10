package rename

import (
	"fmt"
	"os"
	"path/filepath"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/paradoxgery/batch-rename/utils"
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

var RenameCmd = &cobra.Command{
	Use:     "rename",
	Short:   "renames files in batch",
	Long:    "this command will read from a given *.csv with ; as seperator with 'Ur-Pfad+Ur-Datei;Zielpfad+Zieldatei' as header and rename all files listed in the csv\nif there are errors it will list all files that had errors",
	Args:    cobra.ExactArgs(1),
	Example: "rename some.csv",
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("could not find csv: %s\n", path)
			os.Exit(1)
		}
		r := transform.NewReader(f, charmap.ISO8859_1.NewDecoder())

		var renames []FileRename
		err = csvtag.LoadFromReader(r, &renames, csvtag.CsvOptions{Separator: rune(viper.GetString("seperator")[0])})
		if err != nil {
			fmt.Printf("error reading csv file %q: %s\n", path, err)
		}

		if viper.GetBool("verbose") {
			fmt.Println("looking for these files:")
			for _, r := range renames {
				fmt.Println(r.From)
			}
		}

		var rens []FileRename
		fmt.Print(CountTxt)
		count := 0
		maxLen := 0
		fmt.Printf("%d", count)
		var finderrs []RenameErr
		for _, fr := range renames {
			fmt.Print(utils.Deleters(fmt.Sprintf("%d", count)))
			_, err := os.Stat(fr.From)
			if err != nil {
				finderrs = append(finderrs, RenameErr{Fr: fr, err: err})

				fmt.Printf("%d", count)
				continue
			}
			count++
			if len(fr.From) > maxLen {
				maxLen = len(fr.From)
			}
			rens = append(rens, fr)

			fmt.Printf("%d", count)
		}
		fmt.Println("")

		if viper.GetBool("verbose") {
			fmt.Println("Errors:")
			for _, e := range finderrs {
				fmt.Println(e.err)
			}
		}

		fmt.Printf("found %d/%d files\n", count, len(renames))
		if len(rens) == 0 {
			fmt.Println("nothing to do")
			os.Exit(0)
		}
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
		fmt.Printf(RenFmt, count, len(rens), len(errs), utils.Pad(last, maxLen))
		for _, ren := range rens {
			fmt.Print(utils.Deleters(fmt.Sprintf(RenFmt, count, len(rens), len(errs), utils.Pad(last, maxLen))))
			count++
			last = ren.From

			err := os.MkdirAll(filepath.Dir(ren.To), os.ModePerm)
			if err != nil {
				errs = append(errs, RenameErr{Fr: ren, err: err})
				fmt.Printf(RenFmt, count, len(rens), len(errs), utils.Pad(last, maxLen))
				continue
			}

			err = os.Rename(ren.From, ren.To)
			if err != nil {
				errs = append(errs, RenameErr{Fr: ren, err: err})
			}

			fmt.Printf(RenFmt, count, len(rens), len(errs), utils.Pad(last, maxLen))
		}
		fmt.Println(utils.Deleters(utils.Pad(last, maxLen+2)))

		for _, rerr := range errs {
			fmt.Printf("%s\n", rerr.err)
		}

		fmt.Println("done")
	},
}
