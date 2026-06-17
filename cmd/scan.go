package cmd

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"sort"
	"tidyfiles/config"
	"tidyfiles/output"
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
	var files []output.FileEntry
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
			files = append(files, output.FileEntry{
				Path:    s,
				Name:    d.Name(),
				Ext:     filepath.Ext(d.Name()),
				Size:    fileInfo.Size(),
				ModTime: fileInfo.ModTime(),
			})
		}
		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		var less bool
		switch config.AppConfig.Output.SortBy {
		case "date":
			less = files[i].ModTime.Before(files[j].ModTime)
		default:
			less = files[i].Size < files[j].Size
		}
		if config.AppConfig.Output.Ascending {
			return less
		}
		return !less
	})

	output.WriteReport(dirName, total, totalSize, largest, oldest, newest, files, m)
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
