# i6r9c
basic terminal irc client in go, intended for use with tor.

e.g.

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

alternatively, as a container

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

or built locally

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

only the container file is needed for the build, as it will handle fetching the source internally.

`podman` can be placed with `docker`.
