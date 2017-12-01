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
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "prompt",
	Long:  `prompt`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			panic("Not enough arguments to prompt")
		}

		pwd := args[0]
		components := []string{
			formatCurrentPath(pwd),
			formatGitInfo(pwd),
			formatKeymap(),
		}
		fmt.Println(joinNonEmpty(components, " "))
	},
}

func joinNonEmpty(xss []string, sep string) (ret string) {
	for _, xs := range xss {
		if xs != "" {
			ret += xs
			ret += sep
		}
	}
	return
}

func formatKeymap() string {
	return color.WhiteString("λ.")
}

func coalesce(xss ...string) string {
	for _, str := range xss {
		if len(str) > 0 {
			return str
		}
	}
	return ""
}

func formatGitInfo(path string) string {
	action := color.YellowString("")

	repo, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return ""
		}
		panic("Could not open git repository")
	}

	ref, err := repo.Head()
	if err != nil {
		return color.RedString("empty")
	}

	branch := ref.Name().Short()
	if branch != "" {
		return color.BlueString(branch) + action
	}

	cobj, err := repo.CommitObject(ref.Hash())
	if err == nil {
		commit := cobj.ID().String()
		return color.BlueString(commit) + action
	}

	position := ref.Target().String()
	return color.MagentaString(position) + action

}

func formatCurrentPath(path string) string {
	home := os.Getenv("HOME")
	return color.CyanString(trunc(strings.Replace(path, home, "~", 1)))
}

func trunc(path string) string {
	skip := false
	rarr := make([]rune, 16)
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
		rarr = append(rarr, c)
	}

	return string(rarr)
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
