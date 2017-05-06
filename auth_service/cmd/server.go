package cmd

import (
	"github.com/lileio/auth_service/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Fatal(server.NewServer().
			ListenAndServe())
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
