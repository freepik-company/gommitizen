package docs

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

func generateCommandDocumentation(cmd *cobra.Command) CmdData {
	var d CmdData

	if !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand() {
		return d
	}

	// Command name
	d.Name = cmd.Name()
	// Command description
	d.LongDescription = cmd.Long
	// Command short description
	d.ShortDescription = cmd.Short
	// Command usage
	d.Usage = commandUsage(cmd)
	// Command example
	d.Example = cmd.Example
	// Command flags
	d.Flags = retrieveCommandFlags(cmd)
	// Command subcommands
	d.Commands = getCommands(cmd)

	return d
}

func commandUsage(cmd *cobra.Command) string {
	if cmd.Runnable() {
		line := ""
		line = fmt.Sprintf("%s", cmd.UseLine())
		if len(cmd.Commands()) > 0 {
			line = fmt.Sprintf("%s\n%s [command]", line, cmd.CommandPath())
		}

		return line
	}

	return ""
}

func retrieveCommandFlags(cmd *cobra.Command) []FlagData {
	var flags []FlagData

	visitFlag := func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		for _, flag := range flags {
			if flag.Name == f.Name {
				return
			}
		}
		flags = append(flags, FlagData{
			f.Name,
			f.Shorthand,
			f.Usage,
		})
	}
	cmd.Flags().VisitAll(visitFlag)
	cmd.PersistentFlags().VisitAll(visitFlag)

	return flags
}
func getCommands(cmd *cobra.Command) []CmdData {
	var c []CmdData

	if len(cmd.Commands()) > 0 {
		children := cmd.Commands()
		sort.Sort(byName(children))
		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			c = append(c, generateCommandDocumentation(child))
		}
	}

	return c
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
