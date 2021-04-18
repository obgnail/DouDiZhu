package main

import "log"

type CardPartnerTester struct {
	*CardPartner
	nameLabel string
}

func NewCardPartnerTest(cardList CardList) *CardPartnerTester {
	cp := NewCardPartner(cardList)
	return &CardPartnerTester{CardPartner: cp}
}

// 单条
func (cpt *CardPartnerTester) IsSingle() bool {
	return cpt.GetLength() == 1
}

// 一对
func (cpt *CardPartnerTester) IsDouble() bool {
	return cpt.GetLength() == 2 && cpt.GetKindCount() == 1
}

// 三条
func (cpt *CardPartnerTester) IsTriple() bool {
	return cpt.GetLength() == 3 && cpt.GetKindCount() == 1
}

// 普通炸弹
func (cpt *CardPartnerTester) IsNormalBomb() bool {
	return cpt.GetLength() == 4 && cpt.GetKindCount() == 1
}

// 王炸
func (cpt *CardPartnerTester) IsRocket() bool {
	if cpt.GetLength() != 2 || cpt.GetKindCount() != 2 {
		return false
	}
	CardPointBlackJoker, err1 := CardPointTypeEnumConverter.Enum(CardPointBlackJoker)
	CardPointRedJoker, err2 := CardPointTypeEnumConverter.Enum(CardPointRedJoker)
	if err1 != nil || err2 != nil {
		log.Fatal("no such card")
	}
	Jokers := map[int]int{
		CardPointBlackJoker: 1,
		CardPointRedJoker:   1,
	}
	return Equal(cpt.mapPointToCount, Jokers)
}

// 炸弹
func (cpt *CardPartnerTester) IsBomb() bool {
	return cpt.IsRocket() || cpt.IsNormalBomb()
}

// 三带一
func (cpt *CardPartnerTester) IsTripleWithSingle() bool {
	if cpt.GetKindCount() != 2 || cpt.GetLength() != 4 {
		return false
	}
	var hasTripleCard bool
	var hasSingleCard bool
	for _, cardCount := range cpt.mapPointToCount {
		if cardCount == 3 {
			hasTripleCard = true
		} else if cardCount == 1 {
			hasSingleCard = true
		} else {
			return false
		}
	}
	return hasTripleCard && hasSingleCard
}

// 三带二
func (cpt *CardPartnerTester) IsTripleWithDouble() bool {
	if cpt.GetKindCount() != 2 || cpt.GetLength() != 5 {
		return false
	}
	var hasTripleCard bool
	var hasDoubleCard bool
	for _, cardCount := range cpt.mapPointToCount {
		if cardCount == 3 {
			hasTripleCard = true
		} else if cardCount == 2 {
			hasDoubleCard = true
		} else {
			return false
		}
	}
	return hasTripleCard && hasDoubleCard
}

// 四带二
func (cpt *CardPartnerTester) IsQuadrupleWithSingles() bool {
	if cpt.GetKindCount() != 3 && cpt.GetKindCount() != 2 {
		return false
	}
	if cpt.GetLength() != 6 {
		return false
	}
	var hasQuadrupleCard bool
	for _, cardCount := range cpt.mapPointToCount {
		if cardCount == 4 {
			hasQuadrupleCard = true
		}
	}
	return hasQuadrupleCard
}

// 四带两对
func (cpt *CardPartnerTester) IsQuadrupleWithDoublePair() bool {
	if cpt.GetKindCount() != 3 || cpt.GetLength() != 8 {
		return false
	}
	var hasQuadrupleCard bool
	var doubleCardCount int
	for _, cardCount := range cpt.mapPointToCount {
		if cardCount == 4 {
			hasQuadrupleCard = true
		} else if cardCount == 2 {
			doubleCardCount++
		} else {
			return false
		}
	}
	return hasQuadrupleCard && doubleCardCount == 2
}

// 顺子
func (cpt *CardPartnerTester) IsStraight() bool {
	cardListLength := cpt.GetLength()
	cardKindCount := cpt.GetKindCount()
	minCard := cpt.GetMinPointCard()

	// 牌的种类少于5
	if cardKindCount < 5 {
		return false
	}
	// 存在重复的牌
	if cardListLength != cardKindCount {
		return false
	}

	startPoint := minCard.Point
	for i := 0; i < cardListLength; i++ {
		if _, exist := cpt.mapPointToCount[startPoint]; !exist {
			return false
		}
		startPoint++
	}
	return true
}

// 连对
func (cpt *CardPartnerTester) IsStraightPair() bool {
	cardListLength := cpt.GetLength()
	cardKindCount := cpt.GetKindCount()
	minCard := cpt.GetMinPointCard()

	// 牌的种类少于3
	if cardKindCount <= 3 {
		return false
	}
	// 长度不是种类的一半
	if cardListLength != cardKindCount*2 {
		return false
	}

	startPoint := minCard.Point
	for i := 0; i < cardKindCount; i++ {
		if count, exist := cpt.mapPointToCount[startPoint]; !exist || count != 2 {
			return false
		}
		startPoint++
	}
	return true
}

// 不带翅膀的飞机
func (cpt *CardPartnerTester) IsAirPlaneWithoutWing() bool {
	cardListLength := cpt.GetLength()
	cardKindCount := cpt.GetKindCount()

	if cardKindCount*3 != cardListLength {
		return false
	}

	var doubleTripleCount int
	minTripleCardPoint := 100
	for point, cardCount := range cpt.mapPointToCount {
		if cardCount == 3 {
			doubleTripleCount++
			minTripleCardPoint = Min(minTripleCardPoint, point)
		} else {
			return false
		}
	}
	if doubleTripleCount != cardKindCount {
		return false
	}

	// 校验顺子
	startPoint := minTripleCardPoint
	for i := 0; i < cardKindCount; i++ {
		if _, exist := cpt.mapPointToCount[startPoint]; !exist {
			return false
		}
		startPoint++
	}
	return true
}

// 带单翅的飞机
func (cpt *CardPartnerTester) IsAirPlaneWithSingleWing() bool {
	cardListLength := cpt.GetLength()
	var TripleCardList CardList

	hasRightLen := (cardListLength-8)%4 == 0
	if !hasRightLen {
		return false
	}

	for point, cardCount := range cpt.mapPointToCount {
		cards := cpt.GetCardsByCardPoint(point)
		if cardCount == 3 {
			TripleCardList = append(TripleCardList, cards...)
		} else if cardCount == 4 {
			TripleCardList = append(TripleCardList, cards[1:]...)
		}
	}

	airPlaneWithoutWingTester := NewCardPartnerTest(TripleCardList)
	isAirPlaneWithoutWing := airPlaneWithoutWingTester.IsAirPlaneWithoutWing()
	isRightWingLength := (len(TripleCardList)/3)*4 == cardListLength
	return isAirPlaneWithoutWing && isRightWingLength
}

// 带双翅的飞机
func (cpt *CardPartnerTester) IsAirPlaneWithDoubleWing() bool {
	cardListLength := cpt.GetLength()
	var TripleCardList CardList

	hasRightLen := (cardListLength-10)%5 == 0
	var doubleLen int
	if !hasRightLen {
		return false
	}
	for point, cardCount := range cpt.mapPointToCount {
		if cardCount == 3 {
			cards := cpt.GetCardsByCardPoint(point)
			TripleCardList = append(TripleCardList, cards...)
		} else if cardCount == 2 {
			doubleLen++
		} else {
			return false
		}
	}
	airPlaneWithoutWingTester := NewCardPartnerTest(TripleCardList)
	isAirPlaneWithoutWing := airPlaneWithoutWingTester.IsAirPlaneWithDoubleWing()
	isRightWingLength := len(TripleCardList)%3 == doubleLen
	return isAirPlaneWithoutWing && isRightWingLength
}

func (cpt *CardPartnerTester) IsAirPlane() bool {
	return cpt.IsAirPlaneWithoutWing() || cpt.IsAirPlaneWithSingleWing() || cpt.IsAirPlaneWithDoubleWing()
}

func (cpt *CardPartnerTester) Valid() bool {
	cardListLength := cpt.GetLength()
	switch cardListLength {
	case 0:
		return true
	case 1:
		return cpt.IsSingle()
	case 2:
		return cpt.IsDouble() || cpt.IsRocket()
	case 3:
		return cpt.IsTriple()
	case 4:
		return cpt.IsTripleWithSingle() || cpt.IsNormalBomb()
	case 5:
		return cpt.IsTripleWithDouble() || cpt.IsStraight()
	default:
		return cpt.IsStraight() || cpt.IsStraightPair() || cpt.IsAirPlane() || cpt.IsQuadrupleWithSingles() || cpt.IsQuadrupleWithDoublePair()
	}
}

// 判断牌型,0为非法
func (cpt *CardPartnerTester) GetCardPartnerType() int {
	var cardPartnerType string
	cardListLength := cpt.GetLength()
	switch cardListLength {
	case 0:
		cardPartnerType = PartnerTypeEmpty
	case 1:
		if cpt.IsSingle() {
			cardPartnerType = PartnerTypeSingle
		}
	case 2:
		if cpt.IsDouble() {
			cardPartnerType = PartnerTypeDouble
		} else if cpt.IsRocket() {
			cardPartnerType = PartnerTypeRocket
		}
	case 3:
		if cpt.IsTriple() {
			cardPartnerType = PartnerTypeSingle
		}
	case 4:
		if cpt.IsTripleWithSingle() {
			cardPartnerType = PartnerTypeTripleWithSingle
		} else if cpt.IsNormalBomb() {
			cardPartnerType = PartnerTypeNormalBomb
		}
	case 5:
		if cpt.IsTripleWithDouble() {
			cardPartnerType = PartnerTypeTripleWithDouble
		} else if cpt.IsStraight() {
			cardPartnerType = PartnerTypeStraight
		}
	default:
		switch {
		case cpt.IsStraight():
			cardPartnerType = PartnerTypeStraight
		case cpt.IsStraightPair():
			cardPartnerType = PartnerTypeStraightPair
		case cpt.IsAirPlane():
			cardPartnerType = PartnerTypeAirPlane
		case cpt.IsQuadrupleWithSingles():
			cardPartnerType = PartnerTypeQuadrupleWithSingles
		case cpt.IsQuadrupleWithDoublePair():
			cardPartnerType = PartnerTypeQuadrupleWithDoublePair
		}
	}
	res, err := PartnerTypeEnumConverter.Enum(cardPartnerType)
	if err != nil {
		log.Fatal("no such card")
	}
	cpt.nameLabel = cardPartnerType
	return res
}

func (cpt *CardPartnerTester) GetSingleMaxPoint() int {
	return cpt.CardList[0].Point
}

func (cpt *CardPartnerTester) GetDoubleMaxPoint() int {
	return cpt.CardList[0].Point
}

func (cpt *CardPartnerTester) GetTripleMaxPoint() int {
	return cpt.CardList[0].Point
}

func (cpt *CardPartnerTester) GetNormalBombMaxPoint() int {
	return cpt.CardList[0].Point
}

func (cpt *CardPartnerTester) GetRocketMaxPoint() int {
	res, err := PartnerTypeEnumConverter.Enum(CardPointRedJoker)
	if err != nil {
		log.Fatal("no such card")
	}
	return res
}

func (cpt *CardPartnerTester) GetBombMaxPoint() int {
	if cpt.IsNormalBomb() {
		return cpt.GetNormalBombMaxPoint()
	} else if cpt.IsRocket() {
		return cpt.GetRocketMaxPoint()
	}
	return 0
}

func (cpt *CardPartnerTester) GetTripleWithSingleMaxPoint() int {
	cards := cpt.GetCardsByCount(3)
	if len(cards) > 0 {
		return cards[0].Point
	}
	return 0
}

func (cpt *CardPartnerTester) GetTripleWithDoubleMaxPoint() int {
	cards := cpt.GetCardsByCount(3)
	if len(cards) > 0 {
		return cards[0].Point
	}
	return 0
}

func (cpt *CardPartnerTester) GetQuadrupleWithSinglesMaxPoint() int {
	cards := cpt.GetCardsByCount(4)
	if len(cards) > 0 {
		return cards[0].Point
	}
	return 0
}

func (cpt *CardPartnerTester) GetQuadrupleWithDoublePairMaxPoint() int {
	cards := cpt.GetCardsByCount(4)
	if len(cards) > 0 {
		return cards[0].Point
	}
	return 0
}

func (cpt *CardPartnerTester) GetStraightMaxPoint() int {
	maxCard := cpt.GetMaxPointCard()
	if maxCard != nil {
		return maxCard.Point
	}
	return 0
}

func (cpt *CardPartnerTester) GetStraightPairMaxPoint() int {
	maxCard := cpt.GetMaxPointCard()
	if maxCard != nil {
		return maxCard.Point
	}
	return 0
}

func (cpt *CardPartnerTester) GetAirPlaneMaxPoint() int {
	cards1 := cpt.GetCardsByCount(3)
	cards2 := cpt.GetCardsByCount(4)

	var maxPoint1 int
	var maxPoint2 int
	maxCard1 := NewCardPartner(cards1).GetMaxPointCard()
	maxCard2 := NewCardPartner(cards2).GetMaxPointCard()
	if maxCard1 != nil {
		maxPoint1 = maxCard1.Point
	}
	if maxCard2 != nil {
		maxPoint2 = maxCard2.Point
	}
	return Max(maxPoint1, maxPoint2)
}
