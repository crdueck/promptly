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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
)

// rpromptCmd represents the rprompt command
var rpromptCmd = &cobra.Command{
	Use:   "rprompt",
	Short: "rprompt",
	Long:  `rprompt`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Not enough arguments to prompt")
			os.Exit(1)
		}

		pwd := args[0]
		components := []string{
			formatVim(),
			formatGitStatus(pwd),
		}
		fmt.Println(joinNonEmpty(components, " "))
	},
}

func formatVim() string {
	if os.Getenv("VIM") != "" {
		return color.CyanString("")
	}
	return ""
}

func formatGitStatus(path string) string {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return ""
		}
		fmt.Println("Could not open git repository")
		os.Exit(1)
	}

	tree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Unable to open working tree")
		os.Exit(1)
	}

	status, err := tree.Status()
	if err != nil {
		fmt.Println("Could not get working status")
		os.Exit(1)
	}

	var untracked, unmerged, modified, staged bool
	for _, fs := range status {
		switch fs.Staging {
		case git.Unmodified, git.Untracked:
			// no fallthrough
		default:
			staged = true
			continue
		}

		switch fs.Worktree {
		case git.Untracked:
			untracked = true
		case git.Modified:
			modified = true
		case git.UpdatedButUnmerged:
			unmerged = true
		}
	}

	statusString := ""

	if staged {
		statusString += color.GreenString("● ")
	}

	if modified {
		statusString += color.YellowString("● ")
	} else if unmerged {
		statusString += color.YellowString("■ ")
	}

	if untracked {
		statusString += color.RedString("●")
	}

	return statusString
}

func init() {
	rootCmd.AddCommand(rpromptCmd)
}
