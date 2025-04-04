package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ZeroTheorem/gymbot/markups"
	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

const (
	msg  = "<b>%v</b> <i>kg</i> - <b>%v</b> <i>rep.</i> | <b>%v</b> <i>exp</i>.\n"
	msg2 = "<i>Your current level:</i> <b>%v</b>\n<i>The next level</i> <b>%v</b>. <i>Remain</i> <b>%v</b> <i>exp.</i>"
	msg4 = "<i>Your current level:</i> <b>%v</b>\n<i>You've reached the maximum level! Your exp:</i> <b>%v</b>\n<i>Congratulations!</i>🎉"
	msg3 = "\n<i>Total</i>: <b>%v</b> <i>exp.</i>"
	msg5 = "<i>lvl:</i> <b>%v</b>\n<i>To the next level:</i> <b>%v</b> <i>exp</i>."
)

var (
	builder        strings.Builder
	actualExercise string
	expPerTraning  int64
	ChooseExercise bool
)

func main() {
	pref := tele.Settings{
		Token:     "8137726417:AAEcQP9p_ejkUM9KyRvofUzQl0iNJvrT9Fw",
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	menu := markups.CreateMenuSelector()

	b.Handle("/start", func(c tele.Context) error {
		ChooseExercise = false
		return c.Send("<i>Start working out right now!</i>", menu.Selector)
	})

	b.Handle("/reset", func(c tele.Context) error {
		expPerTraning = 0
		builder.Reset()
		menu.Selector.InlineKeyboard[0][0].Text = "Choose exercise"
		return c.Send("<i>Reset compleated!</i>")
	})

	b.Handle("/cmpl", func(c tele.Context) error {
		builder.WriteString(fmt.Sprintf(msg3, expPerTraning))
		c.Send(builder.String())

		// Update user exp and check current lvl
		exp := getExp()
		actualExp := exp + expPerTraning
		writeExp(actualExp)
		c.Send(fmt.Sprintf(msg5, calculateLevel(actualExp), xpToNextLevel(actualExp)))
		// Reset to default settings
		expPerTraning = 0
		builder.Reset()
		menu.Selector.InlineKeyboard[0][0].Text = "Choose exercise"
		return nil
	})

	b.Handle("/cl", func(c tele.Context) error {
		exp := getExp()
		actualLvl := calculateLevel(exp)
		return c.Send(fmt.Sprintf(msg5, actualLvl, xpToNextLevel(exp)))

	})
	b.Handle(menu.ChooseExerciseBtn, func(c tele.Context) error {
		ChooseExercise = true
		return c.Send("<i>Enter the name of the exercise</i>")
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		switch {
		case ChooseExercise:
			actualExercise = c.Message().Text
			ChooseExercise = false
			menu.Selector.InlineKeyboard[0][0].Text = "Change exercise"
			builder.WriteString(fmt.Sprintf("\n<i>%v:</i>\n", actualExercise))
			return c.Send(fmt.Sprintf("<i>Exercise</i>: <b>%v</b>\n<i>Good luck with your approach, bro.</i>", actualExercise), menu.Selector)
		default:
			data := strings.Split(c.Message().Text, " ")
			wight, err := strconv.ParseInt(data[0], 10, 64)
			if err != nil {
				return c.Send("<i>Enter a number</i>")
			}
			reps, err := strconv.ParseInt(data[1], 10, 64)
			if err != nil {
				return c.Send("<i>Enter a number</i>")
			}
			expPerTraning += wight * reps
			builder.WriteString(fmt.Sprintf(msg, wight, reps, wight*reps))
			return c.Send(builder.String(), menu.Selector)
		}
	})
	b.Start()
}
func writeExp(exp int64) {
	f, err := os.OpenFile("data.txt", os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%v", exp))
	if err != nil {
		fmt.Printf("error write to file: %v", err)
		return
	}
}

func getExp() int64 {
	f, err := os.OpenFile("data.txt", os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return 0
	}
	defer f.Close()
	var exp string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		exp = scanner.Text()
	}
	expInt, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		return 0
	}
	return expInt

}

func calculateLevel(xp int64) int64 {
	return int64(math.Floor(math.Sqrt(float64(xp) / 400)))
}

func xpToNextLevel(xp int64) int64 {
	level := calculateLevel(xp)
	nextLevelXP := (level + 1) * (level + 1) * 400
	return nextLevelXP - xp
}
