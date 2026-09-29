package main

const MainnetPort = "8333"

// Seeds from https://github.com/bitcoin/bitcoin/blob/v31.1/src/kernel/chainparams.cpp
var MainnetSeeds = []string{
	"seed.bitcoin.sipa.be",          // Pieter Wuille
	"dnsseed.bluematt.me",           // Matt Corallo
	"seed.bitcoin.jonasschnelli.ch", // Jonas Schnelli
	"seed.btc.petertodd.net",        // Peter Todd
	"seed.bitcoin.sprovoost.nl",     // Sjors Provoost
	"dnsseed.emzy.de",               // Stephan Oeste
	"seed.bitcoin.wiz.biz",          // Jason Maurice
	"seed.mainnet.achownodes.xyz",   // Ava Chow
}
