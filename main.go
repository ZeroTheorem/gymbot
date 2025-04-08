package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

var (
	builder         strings.Builder
	currentExercise string
	sessionExp      int64
	isChoosing      bool
	isTraining      bool
)

func main() {
	pref := tele.Settings{
		Token:     "8137726417:AAEcQP9p_ejkUM9KyRvofUzQl0iNJvrT9Fw",
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	menu := createMenuSelector()

	bot.Handle("/start", func(c tele.Context) error {
		resetState(menu)
		return c.Send("<i>Start working out right now!</i>", menu.Selector)
	})

	bot.Handle("/reset", func(c tele.Context) error {
		resetState(menu)
		return c.Send("<i>Reset completed!</i>")
	})

	bot.Handle("/cmpl", func(c tele.Context) error {
		builder.WriteString(fmt.Sprintf(Msg2, sessionExp))
		c.Send(builder.String())

		data := getData()
		lvl, xp := data[0], data[1]
		newLvl, newXp, messages := updateLevel(lvl, xp, sessionExp)
		for _, msg := range messages {
			c.Send(msg)
		}

		writeData(newLvl, newXp)

		percent := getPercent(newXp, xpToNextLevel(newLvl))
		c.Send(fmt.Sprintf(Msg3, newLvl, defineRank(newLvl), newXp, xpToNextLevel(newLvl),
			generateProgressBar(int(percent)), percent))

		resetState(menu)
		return nil
	})

	bot.Handle("/cl", func(c tele.Context) error {
		data := getData()
		lvl, xp := data[0], data[1]
		percent := getPercent(xp, xpToNextLevel(lvl))
		return c.Send(fmt.Sprintf(Msg3, lvl, defineRank(lvl), xp,
			xpToNextLevel(lvl), generateProgressBar(int(percent)), percent))
	})

	bot.Handle(menu.ChooseExerciseBtn, func(c tele.Context) error {
		isChoosing = true
		isTraining = false
		return c.Send("<i>Enter the name of the exercise</i>")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		switch {
		case isChoosing:
			currentExercise = c.Message().Text
			isChoosing = false
			isTraining = true
			menu.Selector.InlineKeyboard[0][0].Text = "Change exercise"
			builder.WriteString(fmt.Sprintf("\n<i>%v:</i>\n", currentExercise))
			return c.Send(fmt.Sprintf(Msg4, currentExercise), menu.Selector)

		case isTraining:
			input := strings.Split(c.Message().Text, " ")
			if len(input) != 2 {
				return c.Send("<i>Enter two numbers: weight reps</i>")
			}
			weight, err1 := strconv.ParseInt(input[0], 10, 64)
			reps, err2 := strconv.ParseInt(input[1], 10, 64)
			if err1 != nil || err2 != nil {
				return c.Send("<i>Invalid input. Use numbers like: 50 10</i>")
			}
			exp := weight * reps
			sessionExp += exp
			builder.WriteString(fmt.Sprintf(Msg1, weight, reps, exp))
			return c.Send(builder.String(), menu.Selector)

		default:
			return c.Send("<i>Please choose exercise first!</i>", menu.Selector)
		}
	})

	bot.Start()
}
