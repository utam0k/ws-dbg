package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wsdbg",
	Short: "CLI for wsdbg",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.PersistentFlags().StringP("namespace", "n", "k8s.io", "namespace of containerd")
	rootCmd.PersistentFlags().StringP("address", "a", "/run/containerd/containerd.sock", "grpc address back to main containerd")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getNamespace() (string, error) {
	ns, err := rootCmd.PersistentFlags().GetString("namespace")
	if err != nil {
		return "", err
	}
	return ns, nil
}

func getAddress() (string, error) {
	addr, err := rootCmd.PersistentFlags().GetString("address")
	if err != nil {
		return "", err
	}
	return addr, nil
}
