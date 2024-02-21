package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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

		// Create file server and serve files
		fs := http.FileServer(http.Dir(path))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if single {
				// If required file is not found or is a directory, rewrite to index.html
				info, err := os.Stat(filepath.Join(path, r.URL.Path))
				if (err != nil && os.IsNotExist(err)) || (err == nil && info.IsDir()) {
					r.URL.Path = "/"
				}
			}

			fs.ServeHTTP(w, r)
		})

		fmt.Println("Listening on", listener.Addr().String())

		err = http.Serve(listener, nil)
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
