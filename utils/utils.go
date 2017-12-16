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

// RemoveSubDomain 删除域名的随机化子域名(如果存在该攻击或者模式)，仅保留域名本身
func RemoveSubDomain(domain string) string {
	if len(domain) > 1 && strings.HasSuffix(domain, ".") {
		domain = strings.TrimSuffix(domain, ".")
	}
	if domain != "." {
		return domain[strings.Index(domain, ".")+1:]
	}
	return domain
}

// GetTLDDomain 仅仅保留TLD域名
func GetTLDDomain(domain string) string {
	// remove the extra "."
	if len(domain) > 1 && strings.HasSuffix(domain, ".") {
		domain = strings.TrimSuffix(domain, ".")
	}
	if domain != "." {
		return domain[strings.LastIndex(domain, ".")+1:]
	}
	return domain
}
