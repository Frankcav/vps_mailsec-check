package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/foxcpp/mailsec-check/dns"
	"github.com/mitchellh/colorstring"
)

var (
	active = flag.Bool("active", false, "Do some tests that require making connections to the SMTP servers")
)

func printStatus(level Level, name, desc string) {
	var color, mark string
	switch level {
	case LevelUnknown:
		color = "[dark_gray]"
		mark = " "
		desc = "not evaluated;"
	case LevelSecure:
		color = "[green]"
		mark = "+"
	case LevelInsecure:
		color = "[yellow]"
		mark = " "
	case LevelMissing:
		color = "[red]"
		mark = " "
	case LevelInvalid:
		color = "[red]"
		mark = "!"
	}

	colorstring.Println(fmt.Sprintf("[%s%s[reset]] %s[bold]%s:[reset] \t %s", color, mark, color, name, desc))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Println("Usage:", os.Args[0], "<domain>")
		os.Exit(2)
	}
	domain := flag.Args()[0]

	var err error
	extR, err = dns.NewExtResolver()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	res, err := evaluateAll(domain)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	colorstring.Println("[bold]-- Source forgery protection[reset]")
	printStatus(res.dkim, "DKIM", res.dkimDesc)
	printStatus(res.spf, "SPF", res.spfDesc)
	printStatus(res.dmarc, "DMARC", res.dmarcDesc)
	fmt.Println()

	colorstring.Println("[bold]-- TLS enforcement[reset]")
	printStatus(res.mtasts, "MTA-STS", res.mtastsDesc)
	printStatus(res.dane, "DANE", res.daneDesc)
	fmt.Println()

	colorstring.Println("[bold]-- DNS consistency[reset]")
	printStatus(res.fcrdns, "FCrDNS", res.fcrdnsDesc)
	printStatus(res.dnssecMX, "DNSSEC", res.dnssecMXDesc)
}
