package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Completion scripts. urfave/cli drives completion by calling the program back
// with --generate-bash-completion, so the shell side is only the glue that
// forwards the current words and prints the result.
const (
	bashCompletion = `_changelog_complete() {
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
        opts=$( "${COMP_WORDS[@]:0:$COMP_CWORD}" "${cur}" --generate-bash-completion )
    else
        opts=$( "${COMP_WORDS[@]:0:$COMP_CWORD}" --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
    return 0
}
complete -o bashdefault -o default -o nospace -F _changelog_complete changelog
`

	zshCompletion = `#compdef changelog
_changelog_complete() {
    local -a opts
    local cur
    cur=${words[-1]}
    if [[ "$cur" == "-"* ]]; then
        opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
    else
        opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")
    fi
    if [[ "${opts[1]}" != "" ]]; then
        _describe 'values' opts
    else
        _files
    fi
}
compdef _changelog_complete changelog
`
)

func newCompletionCmd() *cli.Command {
	return &cli.Command{
		Name: "completion",
		UsageText: `Completion prints the shell completion script.

    changelog completion bash > /etc/bash_completion.d/changelog
    changelog completion zsh  > "${fpath[1]}/_changelog"`,
		ArgsUsage: "<bash|zsh>",
		Args:      true,
		Action:    completionAction,
	}
}

func completionAction(c *cli.Context) error {
	switch shell := c.Args().First(); shell {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "":
		return usageError("no shell specified: pass bash or zsh")
	default:
		return usageError("unsupported shell %q: pass bash or zsh", shell)
	}

	return nil
}

func newManCmd() *cli.Command {
	return &cli.Command{
		Name: "man",
		UsageText: `Man prints the man page in roff format.

    changelog man > /usr/share/man/man1/changelog.1`,
		Hidden: true,
		Action: manAction,
	}
}

func manAction(c *cli.Context) error {
	// section 1: a command a person runs, not a system administration tool
	page, err := c.App.ToManWithSection(1)
	if err != nil {
		return err
	}

	fmt.Print(page)

	return nil
}
