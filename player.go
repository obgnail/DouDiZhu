package main

import (
	"fmt"
	"sort"
)

type Player struct {
	isLandlord          bool
	position            int
	cardMap             map[*Card]struct{}
	playCardsInThisTurn CardList
}

func NewPlayer(isLandlord bool, position int, cardList CardList) *Player {
	return &Player{
		isLandlord: isLandlord,
		position:   position,
		cardMap:    ChangeCardListToMap(cardList),
	}
}

func (p *Player) IsLandlord() bool {
	return p.isLandlord
}

func (p *Player) GetPosition() int {
	return p.position
}

func (p *Player) GetCardList() CardList {
	var cardList CardList
	for card := range p.cardMap {
		cardList = append(cardList, card)
	}
	sort.Sort(cardList)
	return cardList
}

func (p *Player) AppendCard(cards CardList) CardList {
	for _, card := range cards {
		p.cardMap[card] = struct{}{}
	}
	return p.GetCardList()
}

func (p *Player) RunOut() bool {
	return len(p.cardMap) == 0
}

func (p *Player) PlayCards(playCards CardList) error {
	for _, card := range playCards {
		if _, exist := p.cardMap[card]; exist {
			delete(p.cardMap, card)
		} else {
			return fmt.Errorf("no such card: %v", card)
		}
	}
	p.playCardsInThisTurn = playCards
	return nil
}

func (p *Player) Pass() {
	p.playCardsInThisTurn = nil
}
