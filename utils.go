package main

import (
	"math/rand"
	"time"
)

func Equal(aMap, bMap map[int]int) bool {
	if aMap == nil || bMap == nil {
		return true
	}

	for key, aValue := range aMap {
		if bValue, exist := bMap[key]; !exist || aValue != bValue {
			return false
		}
	}

	for key, bValue := range bMap {
		if aValue, exist := aMap[key]; !exist || aValue != bValue {
			return false
		}
	}
	return true
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func In(x int, slc []int) bool {
	for _, v := range slc {
		if v == x {
			return true
		}
	}
	return false
}

func GenCards() CardList {
	blackJoker, _ := CardPointTypeEnumConverter.Enum(CardPointBlackJoker)
	redJoker, _ := CardPointTypeEnumConverter.Enum(CardPointRedJoker)

	var cards CardList
	for point := range CardPointTypeEnumConverter.byEnum {
		for suit := range CardSuitEnumConverter.byEnum {
			if point == blackJoker || point == redJoker {
				continue
			}
			card := NewCard(point, suit)
			cards = append(cards, card)
		}
	}
	blackJokerCard := NewCard(blackJoker, 0)
	redJokerCard := NewCard(redJoker, 0)
	cards = append(cards, blackJokerCard)
	cards = append(cards, redJokerCard)
	Shuffle(cards)
	return cards
}

func Shuffle(slice CardList) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func ChangeCardListToMap(cards CardList) map[*Card]struct{} {
	res := make(map[*Card]struct{}, len(cards))
	for _, card := range cards {
		res[card] = struct{}{}
	}
	return res
}
