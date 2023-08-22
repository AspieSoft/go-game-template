package game

import (
	"game/gamehandler"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type randSeedHandler struct {
	rand rand.Source
	ind uint64
}

var GameRandSeed *randSeedHandler = &randSeedHandler{rand: rand.NewSource(rand.Int63())}
var InconsistentRand bool = false

func Init(game *gamehandler.Game){
	InconsistentRand = game.InconsistentRand

	//todo: may need to wait for level load menu

	GameRandSeed.rand.Seed(6405275983374102578)

	//todo: init game data and get level info

	go func(){
		for {
			time.Sleep(1 * time.Second)
	
			//todo: add level objects each second
	
		}
	}()
}

func (randSeed *randSeedHandler) Get(min, max int) int {
	if min == max {
		return min
	}else if min > max {
		m := max
		max = min
		min = m
	}

	num := []byte(strconv.Itoa(int(randSeed.rand.Int63())))
	
	minL := len(strings.TrimPrefix(strconv.Itoa(min), "-"))
	maxL := len(strings.TrimPrefix(strconv.Itoa(max), "-"))

	needL := maxL + (maxL-minL) + 2
	for needL > len(num) {
		num = append(num, []byte(strconv.Itoa(int(randSeed.rand.Int63())))...)
	}

	// InconsistentRand will randomly replace indexes with a new random int sometimes
	//
	// this can produce a random number from a consistant seed,
	// giving a uniquely different result each time,
	// while still having some kind of pattern
	//
	// having some kind of subtle pattern gives players the ability to learn that pattern, making the game more fun
	//
	// having the slight inconsistency makes learning that pattern more challenging, preveting the game from being too easy
	//
	// this could also create the illusion of a learning AI, without the need for the complexity of machine learning
	//
	// I also reccomend choosing a different random seed per level
	//
	// for an easy difficulty mode, you could simply disable the 'InconsistentRand' option
	//
	// you can also play around with the math in this method to make things more or less consistant
	if InconsistentRand {
		for i := range num {
			r := strconv.Itoa(int(rand.Int31()))
			if r[1] % 2 == 0 && ((r[2] % 2 == 0 || r[3] % 2 == 0) && (r[4] % 2 == 0 || r[5] % 2 == 0) || (r[2] % 2 == 0 || r[4] % 2 == 0)) {
				num[i] = r[0]
			}
		}
	}

	res := num[:maxL]
	num = num[maxL:]

	if minL < maxL {
		r := res[minL:]
		res = res[:minL]
		rm := num[:(maxL-minL)]
		num = num[(maxL-minL):]
		for i, n := range rm {
			if n % 2 == 0 {
				res = append(res, r[i])
			}
		}
	}

	pos := true
	if min < 0 && max >= 0 {
		pos = num[0] % 2 == 0
		num = num[1:]
	}else if min >= 0 && max < 0 {
		pos = num[0] % 2 != 0
		num = num[1:]
	}else if min < 0 && max < 0 {
		pos = false
	}

	if n, err := strconv.Atoi(string(res)); err == nil {
		dif := int(num[0])
		num = num[1:]
		if dif <= 0 {
			minA := int(math.Abs(float64(min)))
			dif = int(math.Abs(float64(max)))-minA
			if dif > minA && minA >= 0 {
				dif = minA
			}
		}
		if dif <= 0 {
			dif = 1
		}
		for n > max {
			n -= dif
		}
		for n < min {
			n += dif
		}

		if n > max {
			n = max
		}else if n < min {
			n = min
		}

		if !pos {
			n *= -1
		}

		return n
	}

	return min
}
