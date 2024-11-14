package docs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

type DocData struct {
	templateFile string

	Title       string
	Description string
	RootCommand CmdData
}

type CmdData struct {
	Name             string
	LongDescription  string
	ShortDescription string
	Usage            string
	Example          string
	Flags            []FlagData
	Commands         []CmdData
}

type FlagData struct {
	Name        string
	ShortHand   string
	Description string
}

// GenMarkdown creates markdown output.
func GenMarkdown(cmd *cobra.Command, w io.Writer, tmplFile string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	// Root command information
	rootCmdDoc := generateCommandDocumentation(cmd)

	// Root command
	data := DocData{
		templateFile: tmplFile,

		Title:       cmd.Name(),
		Description: cmd.Long,
		RootCommand: rootCmdDoc,
	}

	// Debug
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to YAML: %w", err)
	}
	slog.Info(string(yamlData))
	// Debug

	// Generate the documentation
	templateContent, err := os.ReadFile(data.templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New("readme").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
