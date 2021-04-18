package main

import (
	"fmt"
)

const (
	CardPointThree      = "3"
	CardPointFore       = "4"
	CardPointFive       = "5"
	CardPointSix        = "6"
	CardPointSeven      = "7"
	CardPointEight      = "8"
	CardPointNine       = "9"
	CardPointTen        = "10"
	CardPointJack       = "J"
	CardPointQueen      = "Q"
	CardPointKing       = "K"
	CardPointAce        = "A"
	CardPointTwo        = "2"
	CardPointBlackJoker = "+"
	CardPointRedJoker   = "++"

	CardSuitSpade   = "S" // 黑桃
	CardSuitHeart   = "H" // 红桃
	CardSuitClub    = "C" // 梅花
	CardSuitDiamond = "D" // 方块
)

// 牌型
const (
	PartnerTypeEmpty                   = "Empty"
	PartnerTypeSingle                  = "Single"
	PartnerTypeDouble                  = "Double"
	PartnerTypeTriple                  = "Triple"
	PartnerTypeRocket                  = "Rocket"
	PartnerTypeNormalBomb              = "NormalBomb"
	PartnerTypeTripleWithSingle        = "TripleWithSingle"
	PartnerTypeTripleWithDouble        = "TripleWithDouble"
	PartnerTypeQuadrupleWithSingles    = "QuadrupleWithSingles"
	PartnerTypeQuadrupleWithDoublePair = "QuadrupleWithDoublePair"
	PartnerTypeStraight                = "Straight"
	PartnerTypeStraightPair            = "StraightPair"
	PartnerTypeAirPlane                = "AirPlane"
)

var (
	PartnerTypeEnumConverter   *TypeEnumConverter
	CardSuitEnumConverter      *TypeEnumConverter
	CardPointTypeEnumConverter *TypeEnumConverter
)

type TypeEnumConverter struct {
	byLabel map[string]int
	byEnum  map[int]string
}

func NewTypeEnumConverter() *TypeEnumConverter {
	c := new(TypeEnumConverter)
	c.byLabel = make(map[string]int)
	c.byEnum = make(map[int]string)
	return c
}

func (c *TypeEnumConverter) Set(label string, enum int) {
	c.byLabel[label] = enum
	c.byEnum[enum] = label
}

func (c *TypeEnumConverter) Enum(label string) (int, error) {
	if enum, ok := c.byLabel[label]; ok {
		return enum, nil
	} else {
		return 0, fmt.Errorf("invalid type")
	}
}

func (c *TypeEnumConverter) Label(enum int) (string, error) {
	if label, ok := c.byEnum[enum]; ok {
		return label, nil
	} else {
		return "", fmt.Errorf("invalid type")
	}
}

func NewPartnerTypeConverter() *TypeEnumConverter {
	c := NewTypeEnumConverter()
	c.Set(PartnerTypeEmpty, 1)
	c.Set(PartnerTypeSingle, 2)
	c.Set(PartnerTypeDouble, 3)
	c.Set(PartnerTypeTriple, 4)
	c.Set(PartnerTypeRocket, 5)
	c.Set(PartnerTypeNormalBomb, 6)
	c.Set(PartnerTypeTripleWithSingle, 7)
	c.Set(PartnerTypeTripleWithDouble, 8)
	c.Set(PartnerTypeQuadrupleWithSingles, 9)
	c.Set(PartnerTypeQuadrupleWithDoublePair, 10)
	c.Set(PartnerTypeStraight, 11)
	c.Set(PartnerTypeStraightPair, 12)
	c.Set(PartnerTypeAirPlane, 13)

	return c
}

func NewCardSuitEnumConverter() *TypeEnumConverter {
	c := NewTypeEnumConverter()
	c.Set(CardSuitSpade, 4)
	c.Set(CardSuitHeart, 3)
	c.Set(CardSuitClub, 2)
	c.Set(CardSuitDiamond, 1) // 0用于大小王
	return c
}

func NewCardPointTypeEnumConverter() *TypeEnumConverter {
	c := NewTypeEnumConverter()
	c.Set(CardPointThree, 3)
	c.Set(CardPointFore, 4)
	c.Set(CardPointFive, 5)
	c.Set(CardPointSix, 6)
	c.Set(CardPointSeven, 7)
	c.Set(CardPointEight, 8)
	c.Set(CardPointNine, 9)
	c.Set(CardPointTen, 10)
	c.Set(CardPointJack, 11)
	c.Set(CardPointQueen, 12)
	c.Set(CardPointKing, 13)
	c.Set(CardPointAce, 14)
	c.Set(CardPointTwo, 16) // 2 不能连对,不能顺子,所以不能是15
	c.Set(CardPointBlackJoker, 18)
	c.Set(CardPointRedJoker, 20)
	return c
}

func init() {
	CardPointTypeEnumConverter = NewCardPointTypeEnumConverter()
	CardSuitEnumConverter = NewCardSuitEnumConverter()
	PartnerTypeEnumConverter = NewPartnerTypeConverter()
}
