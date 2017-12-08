package utils

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func UsageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, "Error: %s", msg)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	os.Exit(1)
}

// RemoveSubDomain remove the first part of domain
func RemoveSubDomain(domain string) string {
	if strings.Contains(domain, ".") {
		if domain != "." {
			return domain[strings.Index(domain, ".")+1:]
		}
	}
	return domain
}
