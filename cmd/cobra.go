package cmd

import (
	"errors"
	"fmt"
	"mango-admin/pkg/sdk/pkg"
	"mango-admin/cmd/app"
	"mango-admin/common/global"
	"os"

	"github.com/spf13/cobra"

	"mango-admin/cmd/api"
	"mango-admin/cmd/config"
	"mango-admin/cmd/migrate"
	"mango-admin/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:          "mango-admin",
	Short:        "mango-admin",
	SilenceUsage: true,
	Long:         `mango-admin`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(pkg.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`mango-admin `+global.Version) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	usageStr1 := `也可以参考 https://doc.mango-admin.dev/guide/ksks.html 里边的【启动】章节`
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1)
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(app.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
