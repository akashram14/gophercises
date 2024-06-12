package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Suit: Joker})
	// Output:
	// Ace of Hearts
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 52 {
		t.Error("Wrong no. of cards")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	exp := Card{Rank: Ace, Suit: Space}
	if cards[0] != exp {
		t.Error("Expected ace of spades in the first. Receiveid", cards[0])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != 3 {
		t.Error("Expected 3, found ", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(filter))
	flag := true
	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			flag = false
		}
	}
	if !flag {
		t.Error("Filter did not work!!")
	}
}

func TestShuffle(t *testing.T) {
	shuffleRand = rand.New(rand.NewSource(0))

	orig := New()
	first := orig[40]
	second := orig[35]
	_ = second
	cards := New(Shuffle)

	if cards[0] != first {
		t.Errorf("Expected the first card to be %s, got %s.", first, cards[0])
	}
	if cards[1] != second {
		t.Errorf("Expected the second card to be %s, got %s.", second, cards[1])
	}
}
