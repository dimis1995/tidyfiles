package cmd

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"text/tabwriter"
	"tidyfiles/config"
	"time"

	"github.com/spf13/cobra"
)

type fileEntry struct {
	Path    string
	Name    string
	Ext     string
	Size    int64
	ModTime time.Time
}

var scanCmd = &cobra.Command{
	Use:   "scan [directory]",
	Args:  cobra.ExactArgs(1),
	Short: "displays the content of the given directory",
	Run:   printDirectoryContent,
}

func printDirectoryContent(cmd *cobra.Command, args []string) {
	dirName := args[0]
	m := make(map[string]int)
	var files []fileEntry
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
			files = append(files, fileEntry{
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

	f, err := os.Create(config.AppConfig.Output.Filename + ".txt")
	if err != nil {
		slog.Error("could not create report file", "err", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "Directory report for: %s\n", dirName)
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Fprintf(f, "Found a total of %d files\n", total)
	fmt.Fprintf(f, "Total size: %d bytes\n", totalSize)
	fmt.Fprintf(f, "Largest file size: %d bytes\n", largest)
	fmt.Fprintln(f, "Oldest file modified:  "+oldest.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(f, "Most recent modified:  "+newest.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(f)

	fmt.Fprintln(f, "File extension breakdown:")
	w := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Extension\tCount")
	fmt.Fprintln(w, "---------\t-----")

	exts := make([]string, 0, len(m))
	for ext := range m {
		exts = append(exts, ext)
	}
	sort.Strings(exts)

	for _, ext := range exts {
		fmt.Fprintf(w, "%s\t%d\n", ext, m[ext])
	}
	w.Flush()
	fmt.Fprintln(f)

	fmt.Fprintf(f, "Files (sorted by %s):\n", config.AppConfig.Output.SortBy)
	fw := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
	fmt.Fprintln(fw, "Name\tSize (bytes)\tModified\tPath")
	fmt.Fprintln(fw, "----\t------------\t--------\t----")
	for _, file := range files {
		fmt.Fprintf(fw, "%s\t%d\t%s\t%s\n",
			file.Name,
			file.Size,
			file.ModTime.Format("2006-01-02 15:04:05"),
			file.Path,
		)
	}
	fw.Flush()
	fmt.Println("Report written to " + config.AppConfig.Output.Filename + ".txt")

}

func init() {
	rootCmd.AddCommand(scanCmd)
}
