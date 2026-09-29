package main

type NodeConfig struct {
	ListenPort            string
	ExplicitOutboundPeers []PeerAddress
	MaxOutboundPeers      int
}
