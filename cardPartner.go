package main

// 牌型
type CardPartner struct {
	CardList           CardList
	mapPointToCount    map[int]int
	mapPointToCardList map[int]CardList
}

func NewCardPartner(cardList CardList) *CardPartner {
	cp := &CardPartner{
		CardList:           cardList,
		mapPointToCount:    make(map[int]int),
		mapPointToCardList: make(map[int]CardList),
	}
	for _, c := range cardList {
		cp.mapPointToCount[c.Point]++
		cp.mapPointToCardList[c.Point] = append(cp.mapPointToCardList[c.Point], c)
	}
	return cp
}

// 几种牌
func (cp *CardPartner) GetKindCount() int {
	return len(cp.mapPointToCount)
}

// 几张牌
func (cp *CardPartner) GetLength() int {
	return len(cp.CardList)
}

// 某张牌有几张
func (cp *CardPartner) GetPointCount(card *Card) int {
	return cp.mapPointToCount[card.Point]
}

// 根据 牌号 获取全部的相同牌
func (cp *CardPartner) GetCardsByCardPoint(cardPoint int) CardList {
	return cp.mapPointToCardList[cardPoint]
}

func (cp *CardPartner) GetCardsByCount(count int) CardList {
	var res CardList
	for point, cardCount := range cp.mapPointToCount {
		if cardCount == count {
			cards := cp.GetCardsByCardPoint(point)
			res = append(res, cards...)
		}
	}
	return res
}

// 第一张牌
func (cp *CardPartner) GetFirstCard() *Card {
	return cp.CardList[0]
}

// 最小值的牌
func (cp *CardPartner) GetMinPointCard() *Card {
	minCardPoint := 100
	for cardPoint := range cp.mapPointToCount {
		minCardPoint = Min(minCardPoint, cardPoint)
	}
	cards := cp.GetCardsByCardPoint(minCardPoint)
	if len(cards) > 0 {
		return cards[0]
	}
	return nil
}

func (cp *CardPartner) GetMaxPointCard() *Card {
	maxCardPoint := 0
	for cardPoint := range cp.mapPointToCount {
		maxCardPoint = Max(maxCardPoint, cardPoint)
	}
	cards := cp.GetCardsByCardPoint(maxCardPoint)
	if len(cards) > 0 {
		return cards[0]
	}
	return nil
}
