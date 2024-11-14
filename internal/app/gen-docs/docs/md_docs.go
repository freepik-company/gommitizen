package docs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Flag struct {
	Name        string
	ShortHand   string
	Description string
}

const (
	titleStartingLevel = 2
	templateFile       = "README.md.template"
)

// GenMarkdown creates markdown output.
func GenMarkdown(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	// App description
	descriptionBuf := new(bytes.Buffer)
	descriptionBuf.WriteString(fmt.Sprintf("%s\n\n", cmd.Long))

	// Commands
	commandBuf := new(bytes.Buffer)
	generateCommandDocumentation(cmd, commandBuf, 0)

	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New("readme").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		Description string
		Usage       string
	}{
		Description: descriptionBuf.String(),
		Usage:       commandBuf.String(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func generateCommandDocumentation(cmd *cobra.Command, buf *bytes.Buffer, level int) {
	if !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand() {
		return
	}

	title(cmd, buf, level)
	usage(cmd, buf)
	flags(cmd, buf)
	commands(cmd, buf, level)
	example(cmd, buf)
}

func title(cmd *cobra.Command, buf *bytes.Buffer, level int) {
	if level > 0 {
		title := strings.Repeat("#", level+titleStartingLevel)
		buf.WriteString(fmt.Sprintf("%s %s\n\n", title, cmd.Name()))
		buf.WriteString(fmt.Sprintf("%s\n\n", cmd.Long))
	}
}

func usage(cmd *cobra.Command, buf *bytes.Buffer) {
	if cmd.Runnable() {
		line := ""
		line = fmt.Sprintf("%s", cmd.UseLine())
		if len(cmd.Commands()) > 0 {
			line = fmt.Sprintf("%s\n%s [command]", line, cmd.CommandPath())
		}
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", line))
	}
}

func flags(cmd *cobra.Command, buf *bytes.Buffer) {
	var flags []Flag
	visitFlag := func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		for _, flag := range flags {
			if flag.Name == f.Name {
				return
			}
		}
		flags = append(flags, Flag{
			f.Name,
			f.Shorthand,
			f.Usage,
		})
	}
	cmd.Flags().VisitAll(visitFlag)
	cmd.PersistentFlags().VisitAll(visitFlag)

	if len(flags) > 0 {
		var b bytes.Buffer
		var maxNameLen, maxShorthandLen int
		for _, flag := range flags {
			if len(flag.Name) > maxNameLen {
				maxNameLen = len(flag.Name)
			}
			if len(flag.ShortHand) > maxShorthandLen {
				maxShorthandLen = len(flag.ShortHand)
			}
		}
		for _, flag := range flags {
			if flag.ShortHand != "" {
				b.WriteString(fmt.Sprintf("  -%s, --%s  %s%s\n", flag.ShortHand, flag.Name, strings.Repeat(" ", maxNameLen-len(flag.Name)), flag.Description))
			} else {
				b.WriteString(fmt.Sprintf("      --%s  %s%s\n", flag.Name, strings.Repeat(" ", maxNameLen-len(flag.Name)), flag.Description))
			}
		}
		buf.WriteString(fmt.Sprintf("**Flags**:\n\n```\n%s```\n\n", b.String()))
	}
}

func example(cmd *cobra.Command, buf *bytes.Buffer) {
	if len(cmd.Example) > 0 {
		buf.WriteString(fmt.Sprintf("**Examples**:\n\n%s\n\n", cmd.Example))
	}
}

func commands(cmd *cobra.Command, buf *bytes.Buffer, level int) {
	if len(cmd.Commands()) > 0 {
		children := cmd.Commands()
		sort.Sort(byName(children))
		buf.WriteString("**Available commands**:\n\n")
		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			buf.WriteString(fmt.Sprintf("- **%s**: %s\n", child.Name(), child.Short))
		}
		buf.WriteString(fmt.Sprintf("\n"))

		// subcommand information
		for _, child := range children {
			generateCommandDocumentation(child, buf, level+1)
		}
	}
}
