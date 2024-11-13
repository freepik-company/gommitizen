package docs

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/spf13/cobra"
)

// GenMarkdown creates markdown output.
func GenMarkdown(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)

	// ?? Hacer todo una función recursiva o dejar el raiz así???
	//genCommands(cmd, buf)

	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Description\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("\n### Usage:\n\n```\n%s\n%s [command]\n```\n\n", cmd.UseLine(), cmd.CommandPath()))
	}

	buf.WriteString(fmt.Sprintf("Flags:\n\n```\n%s\n```\n\n", cmd.Flags().FlagUsages()))

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```bash\n%s\n```\n\n", cmd.Example))
	}

	children := cmd.Commands()
	sort.Sort(byName(children))
	buf.WriteString("### Available commands\n\n")
	for _, child := range children {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}
		child.Name()
		buf.WriteString(fmt.Sprintf("- **%s**: %s\n", child.Name(), child.Short))
	}
	buf.WriteString(fmt.Sprintf("\n"))

	for _, child := range children {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}
		buf.WriteString(fmt.Sprintf("#### Command: %s\n\n", child.Name()))
		genCommands(child, buf)
	}

	_, err := buf.WriteTo(w)
	return err
}

func genCommands(cmd *cobra.Command, buf *bytes.Buffer) {
	if !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand() {
		return
	}

	buf.WriteString(fmt.Sprintf("%s\n\n", cmd.Long))

	if cmd.Runnable() {
		line := ""
		line = fmt.Sprintf("%s", cmd.UseLine())
		if len(cmd.Commands()) > 0 {
			line = fmt.Sprintf("%s\n%s [command]", line, cmd.CommandPath())
		}
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", line))
	}

	if len(cmd.Flags().FlagUsages()) > 0 {
		buf.WriteString(fmt.Sprintf("Flags:\n\n```\n%s\n```\n\n", cmd.Flags().FlagUsages()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("Example: \n\n")
		buf.WriteString(fmt.Sprintf("```bash\n%s\n```\n\n", cmd.Example))
	}

	if len(cmd.Commands()) > 0 {
		children := cmd.Commands()
		sort.Sort(byName(children))
		buf.WriteString("Available commands:\n\n")
		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			child.Name()
			buf.WriteString(fmt.Sprintf("- **%s**: ", child.Name()))
			genCommands(child, buf)
		}
		buf.WriteString(fmt.Sprintf("\n"))
	}
}
