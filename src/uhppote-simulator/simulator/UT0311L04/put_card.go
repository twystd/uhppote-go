package UT0311L04

import (
	"fmt"
	"net"
	"uhppote-simulator/entities"
	"uhppote/messages"
)

func (s *UT0311L04) putCard(addr *net.UDPAddr, request *messages.PutCardRequest) {
	if request.SerialNumber == s.SerialNumber {
		card := entities.Card{
			CardNumber: request.CardNumber,
			From:       request.From,
			To:         request.To,
			Doors: map[uint8]bool{1: request.Door1,
				2: request.Door2,
				3: request.Door3,
				4: request.Door4,
			},
		}

		s.Cards.Put(&card)

		response := messages.PutCardResponse{
			SerialNumber: s.SerialNumber,
			Succeeded:    true,
		}

		s.send(addr, &response)

		if err := s.Save(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
