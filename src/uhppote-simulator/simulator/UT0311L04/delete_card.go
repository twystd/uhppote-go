package UT0311L04

import (
	"fmt"
	"net"
	"uhppote/messages"
)

func (s *UT0311L04) deleteCard(addr *net.UDPAddr, request *messages.DeleteCardRequest) {
	if request.SerialNumber == s.SerialNumber {

		deleted := s.Cards.Delete(request.CardNumber)

		response := messages.DeleteCardResponse{
			SerialNumber: s.SerialNumber,
			Succeeded:    deleted,
		}

		s.send(addr, &response)

		if deleted {
			if err := s.Save(); err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}
}
