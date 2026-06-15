package cmd

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"tidyfiles/config"
	"time"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [directory]",
	Args:  cobra.ExactArgs(1),
	Short: "displays the content of the given directory",
	Run:   printDirectoryContent,
}

func printDirectoryContent(cmd *cobra.Command, args []string) {
	dirName := args[0]
	m := make(map[string]int)
	var total = 0
	var totalSize = 0
	var largest int64
	var newest time.Time
	var oldest = time.Now()
	filepath.WalkDir(dirName, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if slices.Contains(config.AppConfig.AllowedExtensions, filepath.Ext(d.Name())) {
			m[filepath.Ext(d.Name())] += 1
			total += 1
			fileInfo, err := d.Info()
			if err != nil {
				slog.Debug("Error when getting information on file", "err", err)
				return err
			}
			if fileInfo.Size() > largest {
				largest = fileInfo.Size()
			}
			totalSize += int(fileInfo.Size())
			if fileInfo.ModTime().Before(oldest) {
				oldest = fileInfo.ModTime()
			}
			if fileInfo.ModTime().After(newest) {
				newest = fileInfo.ModTime()
			}
			fmt.Println(d.Name())
		}
		return nil
	})
	fmt.Printf("Found a total of %d files\n", total)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Largest file size is %d\n", largest)
	fmt.Println("Oldest file modified: " + oldest.Format("2006-01-02 15:04:05"))
	fmt.Println("Most recent modified: " + newest.Format("2006-01-02 15:04:05"))
	for ext, count := range m {
		fmt.Printf("%s: %d\n", ext, count)
	}

}

func init() {
	rootCmd.AddCommand(scanCmd)
}
