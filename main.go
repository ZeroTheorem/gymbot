package main

import (
	"bufio"
	"fmt"
	"log"
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
	msg5 = `
<i>level: <b>%v</b></i>

<b>%v/%v</b> <i>exp</i>. | <b>%.2f%%</b>
`
)

var (
	builder        strings.Builder
	actualExercise string
	expPerTraning  int64
	ChooseExercise bool
)

func main() {
	pref := tele.Settings{
		Token:     "5881448051:AAGnJFe2NRnfochJ91PRw6NW73Fu4ufbXrk",
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
		data := getData()
		currentLvl := data[0]
		prevExp := data[1]
		actualExp := prevExp + expPerTraning
		xpForNextLvl := xpToNextLevel(currentLvl)
		for actualExp >= xpForNextLvl {
			currentLvl++
			actualExp = actualExp - xpForNextLvl
			xpForNextLvl = xpToNextLevel(currentLvl)
		}
		fmt.Println(actualExp)
		writeData(currentLvl, actualExp)
		c.Send(fmt.Sprintf(msg5, currentLvl, actualExp, xpForNextLvl, (float64(actualExp) / float64(xpForNextLvl) * 100)))

		// Reset to default settings
		expPerTraning = 0
		builder.Reset()
		menu.Selector.InlineKeyboard[0][0].Text = "Choose exercise"
		return nil
	})

	b.Handle("/cl", func(c tele.Context) error {
		data := getData()
		currentLvl := data[0]
		currentXp := data[1]
		xpForNextLvl := xpToNextLevel(currentLvl)
		return c.Send(fmt.Sprintf(msg5, currentLvl, currentXp, xpForNextLvl, (float64(currentXp) / float64(xpForNextLvl) * 100)))
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
func writeData(lvl, exp int64) {
	f, err := os.OpenFile("data.txt", os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%v\n%v\n", lvl, exp))
	if err != nil {
		fmt.Printf("error write to file: %v", err)
		return
	}

}

func getData() [2]int64 {
	f, err := os.OpenFile("data.txt", os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return [2]int64{}
	}
	defer f.Close()
	var data [2]int64
	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		number := scanner.Text()
		numberConv, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			return [2]int64{}
		}
		data[i] = numberConv
	}
	fmt.Println(data)
	return data

}

func xpToNextLevel(level int64) int64 {
	nextLevelXP := (level + 1) * (level + 1) * 400
	return nextLevelXP
}
