package commands

import (
	"fmt"
)

type GetCardsCommand struct {
}

func (c *GetCardsCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	N, err := ctx.uhppote.GetCards(serialNumber)

	if err != nil {
		return err
	}

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.uhppote.GetCardByIndex(serialNumber, index+1)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", record)
	}

	return nil
}

func (c *GetCardsCommand) CLI() string {
	return "get-cards"
}

func (c *GetCardsCommand) Description() string {
	return "Returns the list of cards stored on the controller"
}

func (c *GetCardsCommand) Usage() string {
	return "<serial number>"
}

func (c *GetCardsCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-cards <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the number of cards in the controller card list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-cards 12345678")
	fmt.Println()
}
