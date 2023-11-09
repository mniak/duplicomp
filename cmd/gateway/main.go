package main

import (
	"os"
	"syscall"
	"time"

	"github.com/mniak/duplicomp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).
		Level(zerolog.TraceLevel).With().
		Timestamp().
		Caller().
		Logger()

	config := LoadConfig()

	var listenAddress string
	var primaryTarget duplicomp.Target
	var shadowTarget duplicomp.Target

	cmd := cobra.Command{
		Use: "gateway",
		Run: func(cmd *cobra.Command, args []string) {
			cmp := LogComparator{
				Logger:           logger,
				AliasesPerMethod: config.AliasesPerMethod(),
				HintsPerMethod:   config.HintsPerMethod(),
			}

			stopGw, err := duplicomp.StartNewGateway(
				listenAddress,
				primaryTarget,
				shadowTarget,
				cmp,
			)
			cobra.CheckErr(err)

			wait(syscall.SIGTERM, syscall.SIGINT)
			stopGw.GracefulStop()
		},
	}

	cmd.Flags().StringVar(&listenAddress, "listen", ":9091", "TCP address to listen on")
	cmd.Flags().StringVar(&primaryTarget.Address, "target", "", "Connection target")
	cmd.Flags().BoolVar(&primaryTarget.UseTLS, "target-tls", false, "Use TLS in connection target")
	cmd.Flags().StringVar(&shadowTarget.Address, "shadow-target", "", "Shadow connection target")
	cmd.Flags().BoolVar(&shadowTarget.UseTLS, "shadow-target-tls", false, "Use TLS in shadow connection target")
	cmd.MarkFlagRequired("target")
	cmd.MarkFlagRequired("shadow-target")

	cmd.Execute()
}
