# i6r9c
basic terminal irc client in go, intended for use with tor.

e.g.

```
$ go install github.com/guest42069/i6r9c/cmd@main
$ ~/go/bin/cmd -h
Usage of /home/guest42069/go/bin/cmd:
  -nick string
    	IRC nickname to use.
  -proxy string
    	URL schema of proxy, [scheme]://[server]:[port]. (default "socks5://127.0.0.1:9050/")
  -sasl string
    	SASL cert and key prefix (I.E foo/bar for foo/bar.crt and foo/bar.key)
  -server string
    	URL schema of server, [scheme]://[server]:[port]. irc for non-TLS, ircs for TLS. (default "ircs://irc.oftc.net:6697/")
  -verify
    	Verify TLS certificates (I.E. an .onion with TLS but no valid cert.) (default true)
```

alternatively, as a container

```
$ podman run --rm -it ghcr.io/guest42069/i6r9c -h
Usage of /irc:
  -nick string
    	IRC nickname to use.
  -proxy string
    	URL schema of proxy, [scheme]://[server]:[port]. (default "socks5://127.0.0.1:9050/")
  -sasl string
    	SASL cert and key prefix (I.E foo/bar for foo/bar.crt and foo/bar.key)
  -server string
    	URL schema of server, [scheme]://[server]:[port]. irc for non-TLS, ircs for TLS. (default "ircs://irc.oftc.net:6697/")
  -verify
    	Verify TLS certificates (I.E. an .onion with TLS but no valid cert.) (default true)
```

or built locally

```
$ podman build -t localhost/i6r9c:latest -f Containerfile
$ podman run -it localhost/i6r9c:latest -h
Usage of /irc:
  -nick string
    	IRC nickname to use.
  -proxy string
    	URL schema of proxy, [scheme]://[server]:[port]. (default "socks5://127.0.0.1:9050/")
  -sasl string
    	SASL cert and key prefix (I.E foo/bar for foo/bar.crt and foo/bar.key)
  -server string
    	URL schema of server, [scheme]://[server]:[port]. irc for non-TLS, ircs for TLS. (default "ircs://irc.oftc.net:6697/")
  -verify
    	Verify TLS certificates (I.E. an .onion with TLS but no valid cert.) (default true)
```
only the container file is needed for the build, as it will handle fetching the source internally.

`podman` can be placed with `docker`.

Usage of the client (see: `/help`):

```
$ podman run -it --rm ghcr.io/guest42069/i6r9c:latest
[01:25:48] [kinetic.oftc.net@AUTH] *** Looking up your hostname... []
[01:25:48] [kinetic.oftc.net@AUTH] *** Checking Ident []
[01:25:48] [kinetic.oftc.net@AUTH] *** Couldn't look up your hostname []
[01:25:48] [kinetic.oftc.net@AUTH] *** No Ident response []
[01:25:48] [kinetic.oftc.net@quirky_edison] *** Connected securely via TLSv1.3 TLS_AES_128_GCM_SHA256-128 []
[01:25:48] [001] [kinetic.oftc.net@quirky_edison] Welcome to the OFTC Internet Relay Chat Network quirky_edison []
...
[01:25:48] [MODE] [quirky_edison@quirky_edison] +i []
[01:25:48] [kinetic.oftc.net@quirky_edison] Activating Cloak: 8VQAAFOZ0.tor-irc.dnsbl.oftc.net []
> /help
/msg <#channel/recipient> [message]
/join <#channel>
/part <#channel> [reason]
/nick <newnick>
/quit [reason]
> /join #oftc
[01:27:05] [quirky_edison!~quirky_ed@8VQAAFOZ0.tor-irc.dnsbl.oftc.net] has joined [#oftc]
[01:27:05] [332] [kinetic.oftc.net@quirky_edison] https://www.oftc.net | OFTC's public support channel. Our social channel is #moocows | Do NOT paste spam when reporting it | Want a cloak/vhost? See https://www.oftc.net/UserCloaks/ | https://www.oftc.net/Privacy_Policy/ [#oftc]
...
[01:27:05] [366] [kinetic.oftc.net@quirky_edison] End of /NAMES list. [#oftc]
#oftc>
```

this client is intentionally minimalist.
