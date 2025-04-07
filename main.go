package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ZeroTheorem/gymbot/markups"
	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
		builder.WriteString(fmt.Sprintf(Msg2, expPerTraning))
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
		percent := getPercent(actualExp, xpForNextLvl)
		c.Send(fmt.Sprintf(Msg3, currentLvl, actualExp, xpForNextLvl, generateProgressBar(int(percent)), percent))

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
		percent := getPercent(currentLvl, xpForNextLvl)
		return c.Send(fmt.Sprintf(Msg3, currentLvl,
			currentXp, xpForNextLvl, generateProgressBar(int(percent)), percent))
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
			return c.Send(fmt.Sprintf(Msg4, actualExercise), menu.Selector)
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
			builder.WriteString(fmt.Sprintf(Msg1, wight, reps, wight*reps))
			return c.Send(builder.String(), menu.Selector)
		}
	})
	b.Start()
}
func writeData(lvl, exp int64) {
	f, err := os.Create("data.csv")
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.Write([]string{fmt.Sprint(lvl), fmt.Sprint(exp)})
	if err != nil {
		fmt.Printf("error write to file: %v", err)
		return
	}

}

func getData() [2]int64 {
	f, err := os.OpenFile("data.csv", os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return [2]int64{}
	}
	defer f.Close()
	reader := csv.NewReader(f)
	row, err := reader.Read()
	if err != nil {
		fmt.Println("err write file")
	}
	var data [2]int64
	for i, v := range row {
		numberConv, err := strconv.ParseInt(v, 10, 64)
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
func generateProgressBar(percent int) string {
	completed := percent * 20 / 100
	bar := strings.Repeat("█", completed) + strings.Repeat("░", 20-completed)
	return bar
}

func getPercent(num1, num2 int64) float64 {
	return (float64(num1) / float64(num2)) * 100
}
