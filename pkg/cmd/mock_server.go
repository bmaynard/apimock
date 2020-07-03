package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/bmaynard/apimock/pkg/server"
)

type ServerOptions struct {
	pemPath         string
	keyPath         string
	gracefulTimeout int64
	addr            string
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		gracefulTimeout: 15,
		addr:            "127.0.0.1:8000",
		pemPath:         "",
		keyPath:         "",
	}
}

func NewCmdMockServer() *cobra.Command {
	o := NewServerOptions()

	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Run the mock server to respond to API requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.pemPath, "pemPath", "p", o.pemPath, "The path to the pem file")
	cmd.Flags().StringVarP(&o.keyPath, "keyPath", "k", o.pemPath, "The path to the key file")
	cmd.Flags().Int64VarP(&o.gracefulTimeout, "graceful-timeout", "g", o.gracefulTimeout, "The time for which the server will gracefully wait for existing connections to finish")
	cmd.Flags().StringVarP(&o.addr, "addr", "a", o.addr, "The listen addr eg: 127.0.0.1:8000")

	return cmd
}

func (o *ServerOptions) Run() error {
	server.Serve(time.Duration(o.gracefulTimeout)*time.Second, o.addr, o.pemPath, o.keyPath)
	return nil
}
