package main

import (
	"fmt"
	"sort"
)

type Player struct {
	isDiZhu             bool
	position            int
	cardMap             map[*Card]byte
	playCardsInThisTurn CardList
}

func NewPlayer(isDiZhu bool, position int, cardList CardList) *Player {
	p := new(Player)
	p.isDiZhu = isDiZhu
	p.position = position
	p.cardMap = ChangeCardListToMap(cardList)
	return p
}

func (p *Player) IsDiZhu() bool {
	return p.isDiZhu
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
		p.cardMap[card] = 1
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
