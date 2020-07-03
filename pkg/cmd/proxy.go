package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bmaynard/apimock/pkg/proxy"
)

type ProxyOptions struct {
	pemPath string
	keyPath string
	addr    string
	isK8s   bool
}

func NewProxyOptions() *ProxyOptions {
	return &ProxyOptions{
		pemPath: "",
		keyPath: "",
		addr:    "127.0.0.1:8888",
		isK8s:   false,
	}
}

func NewCmdProxyServer() *cobra.Command {
	o := NewProxyOptions()

	var cmd = &cobra.Command{
		Use:   "proxy",
		Short: "Run a proxy server to capture requests and save as mock files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.pemPath, "pemPath", "p", o.pemPath, "The path to the pem file")
	cmd.Flags().StringVarP(&o.keyPath, "keyPath", "k", o.pemPath, "The path to the key file")
	cmd.Flags().StringVarP(&o.addr, "addr", "a", o.addr, "The listen address (default: 127.0.0.1:8888)")
	cmd.Flags().BoolVarP(&o.isK8s, "k8s", "K", o.isK8s, "Is a proxy for a Kubernetes service (will rewrite all requests to localhost)")

	return cmd
}

func (o *ProxyOptions) Run() error {
	proxy.Listen(o.pemPath, o.keyPath, o.addr, o.isK8s)
	return nil
}
