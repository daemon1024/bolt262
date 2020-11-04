package cmd

import (
	"errors"
	"log"
	"os"

	"bolt262/internals/runtests"

	"github.com/spf13/cobra"
)

// SourcePath is the source directory for test262
var SourcePath string

// IncludePath is the directory where helper files are found for test262
var IncludePath string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "the actual test runner / better description for latter :P",
	Long:  `Usage:  bolt262 run [options] <test-file/directory>`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a test file/directory as argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fi, err := os.Stat(args[0])
		if err != nil {
			log.Print(err)
		}
		if fi.Mode().IsDir() {
			runtests.Dir(args[0], IncludePath)
		} else {
			runtests.File(args[0], IncludePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().StringVarP(&SourcePath, "test262Dir", "s", "", "Source directory for test262")
	runCmd.Flags().StringVarP(&IncludePath, "includesDir", "i", "", "directory where helper files are found ")
}
