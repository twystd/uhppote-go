package UT0311L04

import (
	"net"
	"uhppote/messages"
)

func (s *UT0311L04) getCardByID(addr *net.UDPAddr, request *messages.GetCardByIDRequest) {
	if request.SerialNumber == s.SerialNumber {
		response := messages.GetCardByIDResponse{
			SerialNumber: s.SerialNumber,
		}

		for _, card := range s.Cards {
			if request.CardNumber == card.CardNumber {
				response.CardNumber = card.CardNumber
				response.From = &card.From
				response.To = &card.To
				response.Door1 = card.Doors[1]
				response.Door2 = card.Doors[2]
				response.Door3 = card.Doors[3]
				response.Door4 = card.Doors[4]
				break
			}
		}

		s.send(addr, &response)
	}
}

func (s *UT0311L04) getCardByIndex(addr *net.UDPAddr, request *messages.GetCardByIndexRequest) {
	if request.SerialNumber == s.SerialNumber {
		response := messages.GetCardByIndexResponse{
			SerialNumber: s.SerialNumber,
		}

		if request.Index > 0 && request.Index <= uint32(len(s.Cards)) {
			card := s.Cards[request.Index-1]
			response.CardNumber = card.CardNumber
			response.From = &card.From
			response.To = &card.To
			response.Door1 = card.Doors[1]
			response.Door2 = card.Doors[2]
			response.Door3 = card.Doors[3]
			response.Door4 = card.Doors[4]
		}

		s.send(addr, &response)
	}
}
