package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/trugamr/diro/server"
)

func init() {
	rootCmd.Flags().StringP("listen", "l", ":8080", "Address to listen for incoming requests")
	rootCmd.Flags().BoolP("single", "s", false, "Rewrite not-found requests to index.html")
}

var rootCmd = &cobra.Command{
	Use:   "diro [path]",
	Short: "diro is a lightweight file server",
	Long:  "diro is a file server that allows you to serve files or single page applications with ease",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// Collect flags and arguments
		listen := cmd.Flag("listen").Value.String()
		single, _ := cmd.Flags().GetBool("single")

		// If no path is provided, use the current directory
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		// Create new listener and start server
		listener, err := net.Listen("tcp", listen)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer listener.Close()

		// Create new server and start serving
		server := server.New(http.Dir(path), single)

		fmt.Println("Listening on", listener.Addr().String())

		err = http.Serve(listener, server)
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
