package main

import (
	"fmt"
	"log"
)

type Card struct {
	Point int
	Suit  int
}

func NewCard(point, suit int) *Card {
	return &Card{
		Point: point,
		Suit:  suit,
	}
}

func (c *Card) String() string {
	point, err1 := CardPointTypeEnumConverter.Label(c.Point)
	if err1 != nil {
		log.Fatal("no such card")
	}
	if point == CardPointBlackJoker || point == CardPointRedJoker {
		return fmt.Sprintf("%s", point)
	}
	suit, err2 := CardSuitEnumConverter.Label(c.Suit)
	if err2 != nil {
		log.Fatal("no such card")
	}
	return fmt.Sprintf("%s%s", suit, point)
}

func (c *Card) isValid() bool {
	if _, err := CardPointTypeEnumConverter.Label(c.Point); err != nil {
		return false
	}
	if _, err := CardSuitEnumConverter.Label(c.Suit); err != nil {
		return false
	}
	return true
}

type CardList []*Card

func (cl CardList) Len() int           { return len(cl) }
func (cl CardList) Less(i, j int) bool { return cl[i].Point < cl[j].Point }
func (cl CardList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
