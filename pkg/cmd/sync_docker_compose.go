package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
)

type composeInfo struct {
	ProjectName string
	File        string
}

func newSyncDockerComposeCmd(removeCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "docker-compose [profile] [flags]",
		Short: "Sync your docker-compose containers IPs with a profile.",
		Long: `
Reads from a docker-compose.yml file  the list of containers and add names and IPs to a profile in your hosts file.
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			profile, _ := cmd.Flags().GetString("profile")

			if profile == "default" {
				return host.ErrDefaultProfileError
			}
			return nil
		},
		Args: commonCheckArgs,
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			domain, _ := cmd.Flags().GetString("domain")
			network, _ := cmd.Flags().GetString("network")
			prefix, _ := cmd.Flags().GetBool("prefix")

			compose, err := getComposeInfo(cmd)
			if err != nil {
				return err
			}

			profile := profiles[0]

			if profile == "" && compose.ProjectName == "" {
				return host.ErrMissingProfile
			}

			if profile == "" {
				profile = compose.ProjectName
				profiles = append(profiles, profile)
				cmd.SetArgs(profiles)
			}

			if domain == "" {
				domain = fmt.Sprintf("%s.loc", profile)
			}

			ctx := context.Background()

			p, err := host.NewProfileFromDocker(ctx, &host.DockerOptions{
				Domain:      domain,
				Network:     network,
				ComposeFile: compose.File,
				ProjectName: compose.ProjectName,
				KeepPrefix:  prefix,
				Cli:         nil,
			})
			if err != nil {
				return err
			}

			h, err := host.NewFile(src)
			if err != nil {
				return err
			}

			p.Name = profile
			p.Status = host.Enabled

			err = h.AddProfile(*p)
			if err != nil {
				return err
			}

			return h.Flush()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return postActionCmd(cmd, args, removeCmd, false)
		},
	}
}

func getComposeInfo(cmd *cobra.Command) (*composeInfo, error) {
	name, _ := cmd.Flags().GetString("project-name")
	f, _ := cmd.Flags().GetString("compose-file")

	if f == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		f = path.Join(cwd, "docker-compose.yml")
	}

	if name == "" {
		reg := regexp.MustCompile("[^a-z0-9-]+")
		name = path.Base(path.Dir(f))
		name = strings.ToLower(name)
		name = reg.ReplaceAllString(name, "")
	}

	return &composeInfo{
		ProjectName: name,
		File:        f,
	}, nil
}
