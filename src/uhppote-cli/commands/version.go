package commands

import (
	"fmt"
	"uhppote"
)

type VersionCommand struct {
}

func (c *VersionCommand) Execute(ctx Context) error {
	fmt.Printf("%s\n", uhppote.VERSION)

	return nil
}

func (c *VersionCommand) CLI() string {
	return "version"
}

func (c *VersionCommand) Description() string {
	return "Displays the current version"
}

func (c *VersionCommand) Usage() string {
	return ""
}

func (c *VersionCommand) Help() {
	fmt.Println("Displays the uhppote-cli version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}
