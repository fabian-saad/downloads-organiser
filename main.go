package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Skip already existing file or rename and move
var skip = false

var categories = map[string][]string{
	"Bilder, Videos, Audios": {".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico", ".mp3", ".wav", ".mp4", ".mkv", ".mov", ".tif", ".bmp", ".wmf", ".emf", ".wmv"},
	"Textdatei":              {".txt", ".csv", ".md", ".rmd", ".log", ".epub", ".mobi", ".rtf", ".odt", ".odp"},
	"Archive":                {".zip", ".tar", ".gz", ".7z", ".rar", ".xz", ".bz2"},
	"Installer":              {".exe", ".dmg", ".msi", ".pkg", ".deb", ".iso", ".jar"},
	"PDF":                    {".pdf"},
	"Word":                   {".docx", ".doc", ".docm", ".dotx", ".dotm", ".dot"},
	"PowerPoint":             {".pptx", ".ppt", ".pptm", ".potx", ".potm", ".pot", ".ppsx", ".ppsm", ".pps", ".ppam"},
	"Excel":                  {".xlsx", ".xls", ".xlsm", ".xlsb", ".xltx", ".xltm", ".xlt", ".xlam", ".xla"},
	"Code":                   {".py", ".go", ".js", ".jsx", ".ts", ".tsx", ".html", ".css", ".java", ".cpp", ".c", ".h", ".php", ".rs", ".json", ".xml", ".sql", ".sh", ".bat", ".ipynb", ".r", ".rdata", ".tex", ".yaml", ".yml", ".toml", ".sqlite", ".db"},
}

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error retrieving home directory: %v\n", err)
		return
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")
	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		fmt.Printf("Error retrieving downloads directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		item := entry.Name()

		if entry.IsDir() {
			_, isCategory := categories[item]
			if isCategory || item == "Ordner" || item == "Sonstiges" {
				continue
			}
			moveItem(downloadsDir, item, "Ordner")
			continue
		}

		ext := strings.ToLower(filepath.Ext(item))

		moved := false
		for category, extensions := range categories {

			containsExt := slices.Contains(extensions, ext)
			if containsExt {
				moveItem(downloadsDir, item, category)
				moved = true
				break
			}
		}

		if !moved {
			moveItem(downloadsDir, item, "Sonstiges")
		}
	}
}

func moveItem(downloadsDir, item, category string) {

	targetFolder := filepath.Join(downloadsDir, category)

	err := os.MkdirAll(targetFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("Could not create directory: %v\n", err)
		return
	}

	oldPath := filepath.Join(downloadsDir, item)
	newPath := filepath.Join(targetFolder, item)

	if skip {
		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("Skipped (already exists): %s\n", item)
			return
		}
	}

	ext := filepath.Ext(item)
	for i := 1; ; i++ {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		newPath = filepath.Join(targetFolder, fmt.Sprintf("%s (%d)%s", strings.TrimSuffix(item, ext), i, ext))
	}

	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Printf("Error moving %s: %v\n", item, err)
	} else {
		fmt.Printf("Moved: %s -> %s/%s\n", item, category, filepath.Base(newPath))
	}
}
