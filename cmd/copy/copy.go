package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type FileMove struct {
	Data     string `csv:"Ur-Pfad+Ur-Datei"`
	CopyPath string `csv:"Zielpfad"`
}

var CopyCmd = &cobra.Command{
	Use:     "copy",
	Short:   "copies batch of file to a new location",
	Long:    "given a csv with header 'Ur-Pfad+Ur-Datei' and 'Zielpfad' this tool will copy all named files",
	Args:    cobra.ExactArgs(1),
	Example: "copy some.csv",
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("could not find csv: %s\n", path)
			os.Exit(1)
		}
		r := transform.NewReader(f, charmap.ISO8859_1.NewDecoder())

		var copies []FileMove
		err = csvtag.LoadFromReader(r, &copies, csvtag.CsvOptions{Separator: rune(viper.GetString("seperator")[0])})
		if err != nil {
			fmt.Printf("error reading csv file %q: %s\n", path, err)
		}

		var missing []FileMove
		var existing []FileMove
		if viper.GetBool("verbose") {
			fmt.Println("looking for these files:")
		}
		for _, r := range copies {
			if viper.GetBool("verbose") {
				fmt.Println(r.Data)
			}
			_, err := os.Stat(r.Data)
			if err != nil {
				missing = append(missing, r)
				continue
			}

			existing = append(existing, r)
		}

		fmt.Printf("found [%d/%d] files\n", len(existing), len(copies))

		if viper.GetBool("verbose") {
			if len(missing) > 0 {
				fmt.Println("missing files:")
				for _, m := range missing {
					fmt.Printf("\t%s\n", m.Data)
				}
			}
		}

		fmt.Printf("start copying ? (y|N): ")
		var y string
		_, err = fmt.Scanln(&y)
		if err != nil || y != "y" {
			fmt.Println("aborting")
			os.Exit(0)
		}

		errors := make(map[string]error)
		for _, e := range existing {
			err := os.MkdirAll(e.CopyPath, os.ModePerm)
			if err != nil {
				errors[e.Data] = err
				continue
			}

			if viper.GetBool("no-copy") {
				continue
			}

			f, err := os.Open(e.Data)
			if err != nil {
				errors[e.Data] = err
				continue
			}

			dst, err := os.Create(e.CopyPath + "/" + filepath.Base(e.Data))
			if err != nil {
				errors[e.Data] = err
				continue
			}

			_, err = io.Copy(dst, f)
			if err != nil {
				errors[e.Data] = err
				continue
			}
		}

		if len(errors) > 0 {
			fmt.Println("copy errors:")
			for f, err := range errors {
				fmt.Println(">>")
				fmt.Printf("File:\t%s\n", f)
				fmt.Printf("Error:\t%s\n", err)
			}
		}
	},
}

func init() {
	CopyCmd.PersistentFlags().BoolP("no-copy", "n", false, "just create directories but DO NOT copy")
	_ = viper.BindPFlag("no-copy", CopyCmd.PersistentFlags().Lookup("no-copy"))
}
