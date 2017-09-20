package ctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "leappctl",
	Short: "Control interface for leapp-daemon",
	Long: `leapctl is one of the front ends to the LeApp application.
	
	It is designed to help the administrato control the functioning of leapp-daemon along with the execution of commands.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.leapp.yaml)")
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
