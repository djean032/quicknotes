package note

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func ParseNote(filepath string) (*Note, error) {
	rawText, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	note := &Note{
		ID: extractID(filepath),
	}
	frontmatter, content, err := splitFrontmatter(string(rawText))
	if err != nil {
		return nil, fmt.Errorf("Parsing %s: %w", filepath, err)
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
	content = strings.TrimSpace(parts[1])

	return frontmatter, content, nil
}

func extractID(filepath string) string {

}

func extractTitle(frontmatter string) string {

}

func extractTags(frontmatter string) []string {

}

func extractLinks(content string) []string {

}

func extractCreated(frontmatter string) time.Time {

}
