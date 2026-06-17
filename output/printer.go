package output

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"text/tabwriter"
	"tidyfiles/config"
	"time"
)

type FileEntry struct {
	Path    string
	Name    string
	Ext     string
	Size    int64
	ModTime time.Time
}

func WriteReport(dirName string, total int, totalSize int, largest int64, oldest time.Time, newest time.Time, files []FileEntry, m map[string]int) {
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
