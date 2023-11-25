package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type Changelog struct {
	BcChanges   []string
	FeatChanges []string
	FixChanges  []string
	Version     string
	Date        string
	Footer      string
	dirPath     string
}

var stringTemplate = `## {{.Version}} ({{.Date}})
{{- if gt (len .BcChanges) 0 -}}
{{"\n"}}
### BREAKING CHANGES
{{ range .BcChanges }}
- {{.}}
{{- end}}
{{- end}}
{{- if gt (len .FeatChanges) 0 -}}
{{"\n"}}
### Feat
{{range .FeatChanges}}
- {{.}}
{{- end}}
{{- end}}
{{- if gt (len .FixChanges) 0 -}}
{{"\n"}}
### Fix
{{ range .FixChanges}}
- {{.}}
{{- end}}
{{- end}}
{{if .Footer}}
{{.Footer}}
{{end}}
`

const CHANGELOG_FILE = "CHANGELOG.md"

// Public methods

// New creates a new Changelog struct
// version is the version of the project
// dirPath is the path to the CHANGELOG.md file where the file is going to be read from and written to
func New(version string, dirPath string) *Changelog {
	return &Changelog{
		Version: version,
		Date:    time.Now().Format("2006-01-02"),
		dirPath: dirPath,
	}
}

func (c *Changelog) AddBreakingChange(change string) {
	change = removeScope(change)
	c.BcChanges = append(c.BcChanges, change)
}

func (c *Changelog) AddFeatureChange(change string) {
	change = removeScope(change)
	c.FeatChanges = append(c.FeatChanges, change)
}

func (c *Changelog) AddFixChange(change string) {
	change = removeScope(change)
	c.FixChanges = append(c.FixChanges, change)
}

func (c *Changelog) Parse() (*template.Template, error) {
	footer, err := c.footer()
	if err != nil {
		return nil, err
	}

	c.Footer = footer

	return template.Must(template.New("changelog").Parse(stringTemplate)), nil
}

func (c *Changelog) Write() error {
	t, err := c.Parse()
	if err != nil {
		return err
	}

	filePath := filepath.Join(c.dirPath, CHANGELOG_FILE)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating CHANGELOG.md: %s", err)
	}

	err = t.Execute(file, c)
	if err != nil {
		return fmt.Errorf("Error writing CHANGELOG.md: %s", err)
	}

	return nil
}

// Private methods

// footer reads the footer from CHANGELOG.md
func (c *Changelog) footer() (string, error) {
	filePath := filepath.Join(c.dirPath, CHANGELOG_FILE)

	// Check if the file exists and return an empty string if it doesn't
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", nil
	}

	footer, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening CHANGELOG.md: %s", err)
	}

	// retrieve the size of the file
	fileInfo, err := footer.Stat()
	if err != nil {
		return "", fmt.Errorf("Error reading CHANGELOG.md: %s", err)
	}

	// read the whole footer file from CHANGELOG.md
	footerBytes := make([]byte, fileInfo.Size())
	_, err = footer.Read(footerBytes)
	if err != nil {
		return "", fmt.Errorf("Error reading CHANGELOG.md: %s", err)
	}

	return string(footerBytes), nil
}

// Auxiliary functions

// removeScope removes the scope of a change
func removeScope(change string) string {
	// remove the first part of the change delineated by the colon
	// this is the scope of the change
	// e.g. "feat(foo): add bar" -> "add bar"
	// e.g. "feat: add bar" -> "add bar"

	// Find the index of the first colon
	colonIndex := -1
	for i, char := range change {
		if char == ':' {
			colonIndex = i
			break
		}
	}

	// If there is a colon, remove the first part of the string
	if colonIndex != -1 {
		change = change[colonIndex+1:]
	}

	return change
}
