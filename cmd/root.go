package cmd

import (
	"github.com/ariden/auto-add-golang-comments/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "autocomment",
	Short: "goautocommenter",
	Long: `
╔═╗╦ ╦╔╦╗╔═╗  ╔═╗╔═╗╔╦╗╔╦╗╔═╗╔╗╔╔╦╗╔═╗╦═╗
╠═╣║ ║ ║ ║ ║  ║  ║ ║║║║║║║║╣ ║║║ ║ ║╣ ╠╦╝
╩ ╩╚═╝ ╩ ╚═╝  ╚═╝╚═╝╩ ╩╩ ╩╚═╝╝╚╝ ╩ ╚═╝╩╚═
Go auto commenter. Add comments to exported fields automatically
`,
	Run: func(cmd *cobra.Command, args []string) {
		dirs := make([]string, 0)
		files := make([]string, 0)
		if ifDotExist(args) || len(args) == 0 {
			log.Println("Auto commenting the current directory and ignoring other arguments")
			pkg.AutoCommentDir(".")

		} else {
			for index := range args {
				if isDir(args[index]) {
					dirs = append(dirs, args[index])
				} else if isFileExist(args[index]) {
					files = append(files, args[index])
				}
			}
		}

		for index := range dirs {
			pkg.AutoCommentDir(dirs[index])
		}

		pkg.AutoCommentFiles(files...)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
