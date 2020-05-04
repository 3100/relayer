package cmd

import (
	"fmt"

	"github.com/iqlusioninc/relayer/relayer"
	"github.com/spf13/cobra"
)

func crossCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cross",
		Aliases: []string{"cr"},
		Short:   "cross utility commands",
	}

	cmd.AddCommand(
		unacknowledgedPacketsCmd(),
		relayAckPacketCmd(),
	)

	return cmd
}

func unacknowledgedPacketsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unacknowledged-packets [chain-id]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}
			packets, err := chain.QueryUnacknowledgedPackets()
			if err != nil {
				return err
			}
			bz, err := chain.Amino.MarshalJSON(packets)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
		},
	}
	return cmd
}

func relayAckPacketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "relay-ack [path-name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}
			// Fetch latest headers for each chain and store them in sync headers
			sh, err := relayer.NewSyncHeaders(c[src], c[dst])
			if err != nil {
				return err
			}
			return relayer.RelayUnacknowledgedPackets(c[src], c[dst], sh)
		},
	}
	return cmd
}
