/// Queries the Bitcoin DNS seeds
/// This program queries the seeds and retrieves node IP addresses
///
/// Examples:
/// rustc dns.rs && ./dns							Queries the default mainnet seeds
/// rustc dns.rs && ./dns seed.bitcoin.sipa.be		Queries the specifc seed
/// rustc dns.rs && ./dns x9.seed.bitcoin.sipa.be	Filters by service bits
use std::collections::HashSet;
use std::env;
use std::net::ToSocketAddrs;

fn main() {
    // Seeds from https://github.com/bitcoin/bitcoin/blob/v31.1/src/kernel/chainparams.cpp
    let default_hosts = vec![
        "seed.bitcoin.sipa.be",          // Pieter Wuille
        "dnsseed.bluematt.me",           // Matt Corallo
        "seed.bitcoin.jonasschnelli.ch", // Jonas Schnelli
        "seed.btc.petertodd.net",        // Peter Todd
        "seed.bitcoin.sprovoost.nl",     // Sjors Provoost
        "dnsseed.emzy.de",               // Stephan Oeste
        "seed.bitcoin.wiz.biz",          // Jason Maurice
        "seed.mainnet.achownodes.xyz",   // Ava Chow
    ];

    let args: Vec<String> = env::args().skip(1).collect();

    let hosts: Vec<&str>;
    if args.is_empty() {
        hosts = default_hosts;
    } else {
        hosts = args.iter().map(|s| s.as_str()).collect();
    }

    let mut unique_ip_set: HashSet<String> = HashSet::new();
    let mut ip_quantity = 0;

    for host in hosts {
        println!("{}", host);
        match (host, 0).to_socket_addrs() {
            Err(e) => {
                println!("{}", e);
                println!("");
                continue;
            }
            Ok(addrs) => {
                for addr in addrs {
                    let ip = addr.ip();
                    if ip.is_ipv4() {
                        println!("{} → IPv4", ip);
                    } else if ip.is_ipv6() {
                        println!("{} → IPv6", ip);
                    } else {
                        continue;
                    }
                    ip_quantity += 1;
                    unique_ip_set.insert(ip.to_string());
                }
                println!("");
            }
        }
    }
    println!("IPs quantity: {}", ip_quantity);
    println!("Unique IPs: {}", unique_ip_set.len());
}
