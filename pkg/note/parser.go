package note

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ParseNote(fullPath string) (*Note, error) {
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	note := &Note{
		ID: extractID(fullPath),
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}
	}

	frontmatter, content, err := splitFrontmatter(string(rawText))
	if err != nil {
		return nil, fmt.Errorf("Parsing %s: %w", fullPath, err)
	}
	note.Title = extractTitle(frontmatter)
	note.Created = extractCreated(frontmatter)
	note.Tags = extractTags(frontmatter)
	note.Links = extractLinks(string(content))

	return note, nil
}

func splitFrontmatter(content string) (frontmatter string, body string, err error) {
	parts := strings.SplitN(content, "---", 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed note: no --- seperator found")
	}
	frontmatter = strings.TrimSpace(parts[0])
	body = strings.TrimSpace(parts[1])

	return frontmatter, body, nil
}

func extractID(fullPath string) string {
	filePlusExt := filepath.Base(fullPath)
	fileName := strings.Split(filePlusExt, ".")
	return fileName[0]
}

func extractTitle(frontmatter string) (title string, err error) {
	start := strings.Index(frontmatter, "#")
	if start == -1 {
		return "", fmt.Errorf("Malformed note: no title header found")
	}
	end := strings.Index(frontmatter, "Created:")
	if end == -1 {
		return "", fmt.Errorf("Malformed note: frontmatter contains no 'Created:' field")
	}
	title = frontmatter[start+1 : end]
	return title, nil
}

func extractTags(frontmatter string) (tags []string, err error) {
	start := strings.Index(frontmatter, "Tags:")
	if start == -1 {
		return tags, nil
	}
	end := strings.Index(frontmatter, "---")
	if end == -1 {
		return tags, fmt.Errorf("Empty Note: No body found")
	}
	tagString := frontmatter[start+1 : end]
	tagSeq := strings.SplitSeq(tagString, ",")
	for tag := range tagSeq {
		tags = append(tags, tag)
	}
	return tags, nil
}

func extractLinks(content string) []string {

}

func extractCreated(frontmatter string) time.Time {

}
