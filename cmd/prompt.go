// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the License for the specific language governing permissions and limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	git "gopkg.in/libgit2/git2go.v26"
)

var keymapFlag string

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "prompt",
	Long:  `prompt`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd := args[0]
		components := []string{
			formatCurrentPath(pwd),
			formatGitInfo(pwd),
			formatKeymap(),
		}
		fmt.Println(joinNonEmpty(components, " "))
	},
}

func joinNonEmpty(xss []string, sep string) string {
	return strings.Join(filterEmpty(xss), sep)
}

func filterEmpty(in []string) (out []string) {
	for _, s := range in {
		if s != "" {
			out = append(out, s)
		}
	}
	return
}

func formatKeymap() string {
	sym := "λ."

	if keymapFlag == "vicmd" {
		return color.HiWhiteString(sym)
	}

	return color.WhiteString(sym)
}

func formatGitInfo(path string) string {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return ""
	}

	ref, err := repo.Head()
	if err != nil {
		return ""
	}

	gi := &GitInfo{
		Action: showRepositoryState(repo.State()),
		Branch: ref.Shorthand(),
		Commit: ref.Target().String(),
	}

	return gi.String()
}

func showRepositoryState(state git.RepositoryState) string {
	switch state {
	case git.RepositoryStateMerge:
		return "merge"
	case git.RepositoryStateRevert:
		return "revert"
	case git.RepositoryStateCherrypick:
		return "cherry-pick"
	case git.RepositoryStateBisect:
		return "bisect"
	case git.RepositoryStateRebase:
		return "rebase"
	case git.RepositoryStateRebaseInteractive:
		return "rebase-interactive"
	case git.RepositoryStateRebaseMerge:
		return "rebase-merge"
	default:
		return ""
	}
}

type GitInfo struct {
	Action string
	Branch string
	Commit string
}

func (gi *GitInfo) String() (out string) {
	var info string

	if gi.Branch != "" && gi.Branch != "HEAD" {
		info += color.BlueString(gi.Branch)
	} else if len(gi.Commit) > 1000 {
		info += color.BlueString(gi.Commit[:7])
	} else {
		res, err := exec.Command("git", "describe", "--all", "--contains").Output()
		if err == nil {
			desc := string(res[:len(res)-1]) // strip trailing '\n'
			info += color.MagentaString(desc)
		}
	}

	if gi.Action != "" {
		info += ":" + color.YellowString(gi.Action)
	}

	return info
}

func formatCurrentPath(path string) string {
	home := os.Getenv("HOME")
	return color.CyanString(trunc(strings.Replace(path, home, "~", 1)))
}

func trunc(path string) string {
	skip := false
	buff := make([]rune, 16)
	last := strings.LastIndex(path, "/")

	for i, c := range path {
		switch c {
		case '~':
			if skip {
				continue
			}
		case '.':
			if i == 0 {
				skip = false
			} else {
				continue
			}
		case '/':
			skip = false
		default:
			if skip && i < last {
				continue
			} else {
				skip = true
			}
		}
		buff = append(buff, c)
	}

	return string(buff)
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().StringP("keymap", "k", keymapFlag, "Provide the $KEYMAP environment variable")
}
