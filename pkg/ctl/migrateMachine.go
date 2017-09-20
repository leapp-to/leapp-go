package ctl

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// migrateMachineCmd represents the migrateMachine command
var (
	sourceIP string
	targetIP string

	migrateMachineCmd = &cobra.Command{
		Use:   "migrate-machine",
		Short: "Executes a migration of a VM into a macrocontainer",
		Long: `This command migrates one or more application into containers by creating a macrocontainer.

This means that the entire system will be converted into a container, possibly bringing all the dirty with it.`,

		Args: func(cmd *cobra.Command, args []string) error {
			if sourceIP == "" {
				return errors.New("source-ip is mandatory")
			}
			if targetIP == "" {
				return errors.New("target-ip is mandatory")
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("target-ip %s\n", targetIP)
		},
	}
)

func init() {
	RootCmd.AddCommand(migrateMachineCmd)

	migrateMachineCmd.PersistentFlags().StringVarP(&sourceIP, "source-ip", "s", "", "Source machine IP address")
	migrateMachineCmd.PersistentFlags().StringVarP(&targetIP, "target-ip", "t", "", "Target machine IP address")
	migrateMachineCmd.MarkPersistentFlagRequired("target-ip")
}
