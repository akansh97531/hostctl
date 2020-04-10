package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a profile from your hosts file.",
	Long: `
Completely remove a profile content from your hosts file.
It cannot be undone unless you have a backup and restore it.

If you want to remove a profile but would like to use it later,
use 'hosts disable' instead.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		all, _ := cmd.Flags().GetBool("all")

		if !all && profile == "" {
			return host.MissingProfileError
		}

		if profile == "default" {
			return host.DefaultProfileError
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		dst, _ := cmd.Flags().GetString("host-file")
		quiet, _ := cmd.Flags().GetBool("quiet")

		err := host.RemoveProfile(dst, profile)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Printf("Profile '%s' removed.\n\n", profile)
		}

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return postActionCmd(cmd, args, nil)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().Bool("all", false, "Remove all profiles")

	removeCmd.AddCommand(removeDomainsCmd)
}
