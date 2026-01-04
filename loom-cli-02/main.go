// Package main is the entry point for loom-cli.
//
// Implements: l2/package-structure.md PKG-001
// See: l2/interface-contracts.md, l2/tech-specs.md TS-ARCH-001
package main

import (
	"os"

	"loom-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
