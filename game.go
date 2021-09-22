package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Game struct {
	Cards         map[*Card]struct{}
	landlordCards [3]*Card
	player        [3]*Player

	currentPlayer *Player
	LastCards     *CardPartnerTester
}

func NewGame() *Game {
	cards := GenCards()
	g := new(Game)
	g.Cards = ChangeCardListToMap(cards)

	landlordCards := cards[:3]
	otherCards := cards[3:]
	for i := 0; i < 3; i++ {
		beginIdx := 17 * i
		endIdx := beginIdx + 17
		cards := otherCards[beginIdx:endIdx]
		g.player[i] = NewPlayer(false, i, cards)
	}

	g.landlordCards = [3]*Card{landlordCards[0], landlordCards[1], landlordCards[2]}
	return g
}

func (g *Game) GetPlayer(position int) *Player {
	return g.player[position]
}

func (g *Game) GetCurrentPlayer() *Player {
	return g.currentPlayer
}

func (g *Game) SetCurrentPlayer(player *Player) {
	g.currentPlayer = player
}

func (g *Game) TurnNextPlayer() *Player {
	position := g.GetCurrentPlayer().position
	position++
	if position == 3 {
		position = 0
	}
	player := g.GetPlayer(position)
	g.SetCurrentPlayer(player)
	return player
}

func (g *Game) GetLastCards() *CardPartnerTester {
	return g.LastCards
}

func (g *Game) GetLandlord() *Player {
	for i := 0; i < 3; i++ {
		p := g.GetPlayer(i)
		if p.IsLandlord() {
			return p
		}
	}
	return nil
}

func (g *Game) GetLandlordCards() [3]*Card {
	return g.landlordCards
}

func (g *Game) PlayerBecomeLandlord(position int) {
	fmt.Printf("player %d became landlord\n", position)
	player := g.GetPlayer(position)
	player.isLandlord = true
	landlordCards := CardList{g.landlordCards[0], g.landlordCards[1], g.landlordCards[2]}
	player.AppendCard(landlordCards)
	g.currentPlayer = player
}

func (g *Game) CurrentPlayerHasThereCard(playCards CardList) bool {
	currentPlayer := g.GetCurrentPlayer()

	for _, card := range playCards {
		if _, exist := currentPlayer.cardMap[card]; !exist {
			return false
		}
	}
	return true
}

func (g *Game) checkValidPlayCards(playCards CardList) bool {
	if len(playCards) == 0 {
		return true
	}
	cpt := NewCardPartnerTest(playCards)
	if !cpt.Valid() {
		return false
	}
	// 获取任意出牌权后的出牌
	if g.LastCards == nil {
		return true
	}
	currentPlayerCardPartnerType := cpt.GetCardPartnerType()
	lastPlayerCardPartnerType := g.LastCards.GetCardPartnerType()
	normalBomb, _ := PartnerTypeEnumConverter.Enum(PartnerTypeNormalBomb)
	rocket, _ := PartnerTypeEnumConverter.Enum(PartnerTypeRocket)

	// 炸弹做特殊判断
	// 王炸
	if currentPlayerCardPartnerType == rocket {
		return true
	}
	if lastPlayerCardPartnerType == rocket {
		return false
	}
	// 炸弹
	if currentPlayerCardPartnerType == normalBomb && lastPlayerCardPartnerType != normalBomb {
		return true
	}

	if currentPlayerCardPartnerType != lastPlayerCardPartnerType {
		return false
	}

	newMaxCard := cpt.GetMaxPointCard()
	oldMaxCard := g.LastCards.GetMaxPointCard()
	if newMaxCard != nil && oldMaxCard != nil && newMaxCard.Point > oldMaxCard.Point {
		return true
	}
	return false
}

func (g *Game) resetLastCards(playCards CardList) {
	if playCards == nil {
		g.LastCards = nil
	} else {
		cpt := NewCardPartnerTest(playCards)
		g.LastCards = cpt
	}
}

func (g *Game) removeCards(Cards CardList) error {
	for _, card := range Cards {
		if _, exist := g.Cards[card]; exist {
			delete(g.Cards, card)
		} else {
			return fmt.Errorf("no such card: %v", card)
		}
	}
	return nil
}

func (g *Game) PlayerPlayCards(position int, playCards CardList) error {
	if len(playCards) == 0 {
		return nil
	}
	success := g.checkValidPlayCards(playCards)
	if !success {
		return fmt.Errorf("player play error cards, %v", playCards)
	}

	play := g.GetPlayer(position)
	if err := play.PlayCards(playCards); err != nil {
		return err
	}
	if err := g.removeCards(playCards); err != nil {
		return err
	}
	g.resetLastCards(playCards)
	return nil
}

func (g *Game) CurrentPlayerPlayCards(playCards CardList) error {
	current := g.GetCurrentPlayer()
	return g.PlayerPlayCards(current.position, playCards)
}

// 当前玩家获得随意出牌权
func (g *Game) CurrentPlayerHasPlayRight() bool {
	hasRight := true
	current := g.GetCurrentPlayer()
	for _, p := range g.player {
		if p != current {
			playerDoesntPlayInThisTurn := p.playCardsInThisTurn == nil
			hasRight = hasRight && playerDoesntPlayInThisTurn
		}
	}
	return hasRight
}

func (g *Game) GetCardsByName(names string) (CardList, error) {
	names = strings.TrimSpace(names)
	nameList := strings.Split(names, " ")
	var res CardList

	for _, name := range nameList {
		card, err := g.GetCardByName(name)
		if err != nil {
			return nil, err
		}
		res = append(res, card)
	}
	return res, nil
}

func (g *Game) GetCardByName(name string) (*Card, error) {
	if name == CardPointBlackJoker || name == CardPointRedJoker {
		point, err1 := CardPointTypeEnumConverter.Enum(name)
		if err1 != nil {
			return nil, fmt.Errorf("----- 没有这种牌")
		}
		suit := 0
		for card := range g.Cards {
			if card.Point == point && card.Suit == suit {
				return card, nil
			}
		}
	}

	nameRune := []rune(name)
	if len(nameRune) == 0 {
		return nil, fmt.Errorf("----- 没有这种牌")
	}

	point, err1 := CardPointTypeEnumConverter.Enum(string(nameRune[1:]))
	suit, err2 := CardSuitEnumConverter.Enum(string(nameRune[0]))
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("----- 没有这种牌")
	}

	for card := range g.Cards {
		if card.Point == point && card.Suit == suit {
			return card, nil
		}
	}
	return nil, nil
}

func (g *Game) GetCurrentPlayerPlayCardFromScan() (cards CardList, err error) {
	currentPlayer := g.GetCurrentPlayer()
	fmt.Println()
	fmt.Printf("player %d , pleace enter your cards(enter [pass] to pass):", currentPlayer.position)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if strings.ToLower(line) == "pass" {
		return nil, nil
	}
	cards, err = g.GetCardsByName(line)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (g *Game) ShowAllPlayerCards() {
	fmt.Println()
	for _, player := range g.player {
		cardList := player.GetCardList()
		var identity string
		if player.IsLandlord() {
			identity = "landlord"
		} else {
			identity = "farmer"
		}
		fmt.Printf("player %d (%s)\t: %s\n", player.position, identity, cardList)
	}
}

func (g *Game) Init() {
	// 没有设计抢地主功能，直接随机地主
	rand.Seed(time.Now().UnixNano())
	landlordPosition := rand.Intn(3)
	g.PlayerBecomeLandlord(landlordPosition)
}

func (g *Game) Play() error {
	player0 := g.GetPlayer(0)
	player1 := g.GetPlayer(1)
	player2 := g.GetPlayer(2)

	for !player0.RunOut() && !player1.RunOut() && !player2.RunOut() {
		g.ShowAllPlayerCards()

		currentPlayerHasPlayRight := g.CurrentPlayerHasPlayRight()
		currentPlayer := g.GetCurrentPlayer()
		if currentPlayerHasPlayRight {
			fmt.Printf("----- player %d 现在拥有任意出牌权\n", currentPlayer.position)
		}

		playCards, err := g.GetCurrentPlayerPlayCardFromScan()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 如果当前玩家拥有任意出牌权,重置LastCards为nil
		if currentPlayerHasPlayRight {
			g.resetLastCards(nil)
			// 拥有任意出牌权的玩家不能不出牌
			if len(playCards) == 0 {
				fmt.Println("----- 拥有任意出牌权的玩家不能不出牌")
				continue
			}
		}

		if len(playCards) == 0 {
			p := g.GetCurrentPlayer()
			p.Pass()
		} else {
			currentPlayerHasThereCard := g.CurrentPlayerHasThereCard(playCards)
			if !currentPlayerHasThereCard {
				fmt.Printf("----- player %d 没有这些牌\n", currentPlayer.position)
				continue
			}
			// 牌型正确,且大于上家出的牌
			rightCardPartner := g.checkValidPlayCards(playCards)
			if !rightCardPartner {
				fmt.Println("----- 牌型错误,或者小于上家出的牌")
				continue
			}
			fmt.Println("------------------------------------------------------------")
			fmt.Println("上家出牌:", playCards)
			err := g.CurrentPlayerPlayCards(playCards)
			if err != nil {
				return err
			}
		}
		g.TurnNextPlayer()
	}
	return nil
}
