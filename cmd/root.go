package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dyne/tgcom/utils/modfile"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

/* In these variables we store the arguments passed to flags -f, -l, -d and -a */
var FileToRead string
var inputFlag modfile.Config

/*
	rootCmd is the command tgcom. "Use" is the name of the command, "Short" is a brief description of the command, "Long

is a longer description of the command, Run is the action that must be executed when command tgcom is called"
*/
var rootCmd = &cobra.Command{
	Use:   "tgcom",
	Short: "tgcom is tool that allows users to comment or uncomment pieces of code",
	Long: `tgcom is a CLI library written in Go that allows users to
    comment or uncomment pieces of code. It supports many different
    languages including Go, C, Java, Python, Bash, and many others...`,

	Run: func(cmd *cobra.Command, args []string) {
		if noFlagsGiven(cmd) {
			customUsageFunc(cmd)
			os.Exit(1)
		}
		ReadFlags(cmd)
	},
}

/* the one command used to run the main function (set by default by cobra-cli) */
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

/*
	pass in this function the flags of the command tgcom. Flags can be Persistend (so that if tgcom has a sub-command, e.g.

subtgcom, the flag defined for tgcom can be used as flags of subtgcom) or local (so flags are usable only for tgcom command)
*/
func init() {

	rootCmd.SetHelpFunc(customHelpFunc)
	rootCmd.SetUsageFunc(customUsageFunc)

	rootCmd.PersistentFlags().StringVarP(&FileToRead, "file", "f", "", "pass argument to the flag and will modify the file content")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.LineNum, "line", "l", "", "pass argument to line flag and will modify the line in the specified range")
	rootCmd.PersistentFlags().BoolVarP(&inputFlag.DryRun, "dry-run", "d", false, "pass argument to dry-run flag and will print the result")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.Action, "action", "a", "toggle", "pass argument to action to comment/uncomment/toggle some lines")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.StartLabel, "start-label", "s", "", "pass argument to start-label to modify lines after start-label")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.EndLabel, "end-label", "e", "", "pass argument to end-label to modify lines up to end-label")
	rootCmd.PersistentFlags().StringVarP(&inputFlag.Lang, "language", "L", "", "pass argument to language to specify the language of the input code")
	rootCmd.MarkFlagsRequiredTogether("start-label", "end-label")
	rootCmd.MarkFlagsMutuallyExclusive("line", "start-label")
	rootCmd.MarkFlagsMutuallyExclusive("line", "end-label")
	rootCmd.MarkFlagsOneRequired("file", "language")
	rootCmd.MarkFlagsMutuallyExclusive("file", "language")

}

/* function to see if no flag is given */
func noFlagsGiven(cmd *cobra.Command) bool {
	hasFlags := false
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			hasFlags = true
		}
	})
	return !hasFlags
}

// ReadFlags parses command line flags and applies them to modify files or display information accordingly.
func ReadFlags(cmd *cobra.Command) {
	if strings.Contains(FileToRead, ",") {
		if cmd.Flags().Changed("line") {
			fmt.Println("Warning: when passed multiple file to flag -f don't use -l flag")
		}
		if cmd.Flags().Changed("start-label") && cmd.Flags().Changed("end-label") {
			fileInfo := strings.Split(FileToRead, ",")
			for i := 0; i < len(fileInfo); i++ {
				inputFlag.Filename = fileInfo[i]
				if err := modfile.ChangeFile(inputFlag); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fileInfo := strings.Split(FileToRead, ",")
			for i := 0; i < len(fileInfo); i++ {
				if strings.Contains(fileInfo[i], ":") {
					parts := strings.Split(fileInfo[i], ":")
					if len(parts) != 2 {
						log.Fatalf("invalid syntax. Use 'File:lines'")
					}
					inputFlag.Filename = parts[0]
					inputFlag.LineNum = parts[1]
					if err := modfile.ChangeFile(inputFlag); err != nil {
						log.Fatal(err)
					}

				} else {
					log.Fatalf("invalid syntax. Use 'File:lines'")
				}
			}
		}
	} else {
		if cmd.Flags().Changed("line") || cmd.Flags().Changed("start-label") && cmd.Flags().Changed("end-label") {
			inputFlag.Filename = FileToRead
			if err := modfile.ChangeFile(inputFlag); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("Not specified what you want to modify: add -l flag or -s and -e flags")
		}

	}
}

func customHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Tgcom CLI Application")
	fmt.Println()
	fmt.Println(cmd.Long)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tgcom [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "action" {
			fmt.Printf("  -%s, --%s: %s (default: %s)\n", flag.Shorthand, flag.Name, flag.Usage, flag.DefValue)
		} else {
			fmt.Printf("  -%s, --%s: %s\n", flag.Shorthand, flag.Name, flag.Usage)
		}
	})
	fmt.Println()
	fmt.Println("Supported Languages:")
	for lang := range modfile.CommentChars {
		fmt.Println(lang)
	}
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Toggle comments on lines 1-5 in example.go")
	fmt.Println("  tgcom -f example.go -l 1-5 -a toggle")
	fmt.Println()
	fmt.Println("  # Dry run: show the changes without modifying the file")
	fmt.Println("  tgcom -f example.go -s START -e END -a toggle -d")
}

func customUsageFunc(cmd *cobra.Command) error {
	fmt.Println(cmd.Short)
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", cmd.UseLine())
	fmt.Println()
	fmt.Println("Flags:")
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "action" {
			fmt.Printf("  -%s, --%s: %s (default: %s)\n", flag.Shorthand, flag.Name, flag.Usage, flag.DefValue)
		} else {
			fmt.Printf("  -%s, --%s: %s\n", flag.Shorthand, flag.Name, flag.Usage)
		}
	})
	return nil
}
