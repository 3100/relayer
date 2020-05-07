package relayer

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	"github.com/datachainlab/cross/x/ibc/cross/types"
	"github.com/tendermint/tendermint/libs/log"
)

const defaultPacketRecvQuery = "recv_packet.packet_src_port=%s&recv_packet.packet_src_channel=%s&recv_packet.packet_sequence=%d"

type PacketStatus = types.UnacknowledgedPacket

func (c *Chain) QueryUnacknowledgedPackets() ([]PacketStatus, error) {
	req := types.QueryUnacknowledgedPacketsRequest{}
	bz, err := c.Amino.MarshalJSON(req)
	if err != nil {
		return nil, err
	}
	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryUnacknowledgedPackets)
	res, _, err := c.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}
	var response types.QueryUnacknowledgedPacketsResponse
	c.Amino.MustUnmarshalJSON(res, &response)
	return response.Packets, nil
}

func (c *Chain) queryRecvPackets(sourcePort, sourceChannel string, sequence uint64) ([]relayPacket, error) {
	eveSend, err := ParseEvents(fmt.Sprintf(defaultPacketRecvQuery, sourcePort, sourceChannel, sequence))
	if err != nil {
		return nil, err
	}

	res, err := c.QueryTxs(0, 1, 10, eveSend)
	if err != nil {
		return nil, err
	}

	if l := len(res.Txs); l == 0 { // there are no unack packets
		return nil, nil
	} else if l > 1 {
		return nil, fmt.Errorf("duplicated sequence on same channel: length=%v sourceChannel=%v sequence=%v", l, sourceChannel, sequence)
	}
	tx := res.Txs[0]
	for i, msg := range tx.Tx.GetMsgs() {
		m, ok := msg.(channel.MsgPacket)
		if !ok {
			continue
		}
		if m.SourceChannel == sourceChannel && m.Sequence == sequence {
			events := validateAndStringifyEvents(tx.Logs[i].Events, c.logger)
			packets, err := parsePackets(events, sourcePort, sourceChannel)
			if err != nil {
				return nil, err
			}
			return packets, nil
		}
	}
	return nil, nil
}

func validateAndStringifyEvents(events sdk.StringEvents, logger log.Logger) map[string][]string {
	result := make(map[string][]string)
	for _, event := range events {
		if len(event.Type) == 0 {
			logger.Debug("Got an event with an empty type (skipping)", "event", event)
			continue
		}

		for _, attr := range event.Attributes {
			if len(attr.Key) == 0 {
				logger.Debug("Got an event attribute with an empty key(skipping)", "event", event)
				continue
			}

			compositeTag := fmt.Sprintf("%s.%s", event.Type, string(attr.Key))
			result[compositeTag] = append(result[compositeTag], string(attr.Value))
		}
	}

	return result
}

func parsePackets(events map[string][]string, sourcePort, sourceChannel string) (rlyPkts []relayPacket, err error) {

	// then, check for packet acks
	if pdval, ok := events["recv_packet.packet_data"]; ok {
		for i, pd := range pdval {
			// Ensure that we only relay over the channel and port specified
			// OPTIONAL FEATURE: add additional filtering options
			srcChan, srcPort := events["recv_packet.packet_src_channel"], events["recv_packet.packet_src_port"]
			// dstChan, dstPort := events["recv_packet.packet_dst_channel"], events["recv_packet.packet_dst_port"]

			if sourcePort == srcPort[i] && sourceChannel == srcChan[i] {
				rp := &relayMsgPacketAck{packetData: []byte(pd)}

				// first get the ack
				if ack, ok := events["recv_packet.packet_ack"]; ok {
					rp.ack = []byte(ack[i])
				}
				// next, get and parse the sequence
				if sval, ok := events["recv_packet.packet_sequence"]; ok {
					seq, err := strconv.ParseUint(sval[i], 10, 64)
					if err != nil {
						return nil, err
					}
					rp.seq = seq
				}

				// finally, get and parse the timeout
				if sval, ok := events["recv_packet.packet_timeout_height"]; ok {
					timeout, err := strconv.ParseUint(sval[i], 10, 64)
					if err != nil {
						return nil, err
					}
					rp.timeout = timeout
				}

				// finally, get and parse the timeout
				if sval, ok := events["recv_packet.packet_timeout_timestamp"]; ok {
					timeout, err := strconv.ParseUint(sval[i], 10, 64)
					if err != nil {
						return nil, err
					}
					rp.timeoutStamp = timeout
				}

				// queue the packet for return
				rlyPkts = append(rlyPkts, rp)
			}
		}
	}

	return rlyPkts, nil
}
