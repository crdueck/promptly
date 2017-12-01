// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	git "gopkg.in/libgit2/git2go.v26"
)

// rpromptCmd represents the rprompt command
var rpromptCmd = &cobra.Command{
	Use:   "rprompt",
	Short: "rprompt",
	Long:  `rprompt`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd := args[0]
		components := []string{
			formatVim(),
			formatGitStatus(pwd),
		}
		fmt.Println(joinNonEmpty(components, " "))
	},
}

func formatVim() string {
	if os.Getenv("VIM") == "" {
		return ""
	}
	return color.CyanString("")
}

type StatusSummary struct {
	Untracked bool
	Unmerged  bool
	Modified  bool
	Staged    bool
}

func formatGitStatus(path string) string {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return ""
	}

	statuses, err := repo.StatusList(&git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptExcludeSubmodules | git.StatusOptIncludeUntracked | git.StatusOptDisablePathspecMatch,
	})

	if err != nil {
		return ""
	}

	return summarizeGitStatusList(statuses).String()
}

func (ws *StatusSummary) String() string {
	out := []string{}

	if ws.Staged {
		out = append(out, color.GreenString("●"))
	}

	if ws.Modified {
		out = append(out, color.YellowString("●"))
	} else if ws.Unmerged {
		out = append(out, color.YellowString("■"))
	}

	if ws.Untracked {
		out = append(out, color.RedString("●"))
	}

	return strings.Join(out, " ")
}

func summarizeGitStatusList(statuses *git.StatusList) *StatusSummary {
	summ := new(StatusSummary)

	count, err := statuses.EntryCount()
	if err != nil {
		return summ
	}

	for i := 0; i < count; i++ {
		entry, err := statuses.ByIndex(i)
		if err != nil {
			continue
		}

		switch entry.Status {
		case git.StatusIndexNew,
			git.StatusIndexModified,
			git.StatusIndexDeleted,
			git.StatusIndexRenamed:
			summ.Staged = true
		case git.StatusWtModified:
			summ.Modified = true
		case git.StatusWtNew:
			summ.Untracked = true
		}
	}

	return summ
}

func init() {
	rootCmd.AddCommand(rpromptCmd)
}
