package cmd

import (
	"fmt"
	"os"

	lacon "github.com/laconlab/lacon-go-tiny-scrapy/pkg"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "scrapy",
		Short: "downloads websites according to provided yml config",
		Long:  "downloads websites according to provided yml config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config-path")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return lacon.RunScrapy(cfgPath)
		},
	}

	cmd.Flags().String(
		"config-path",
		"resources/application.yml",
		"path to config file",
	)

	root.AddCommand(cmd)
}
