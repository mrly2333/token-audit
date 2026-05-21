package migrations

import (
	"embed"
	"io/fs"
	"sort"
)

//go:embed *.sql
var Files embed.FS

func Ordered() ([]string, error) {
	entries, err := fs.ReadDir(Files, ".")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names, nil
}
