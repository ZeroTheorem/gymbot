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
)

const (
	L1  int64 = 0
	L2  int64 = 100_000
	L3  int64 = 300_000
	L4  int64 = 600_000
	L5  int64 = 100_000_000
	L6  int64 = 100_500_000
	L7  int64 = 200_100_000
	L8  int64 = 200_800_000
	L9  int64 = 300_600_000
	L10 int64 = 400_500_000
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
		return c.Send("<i>Start working out right now!</i>", menu.Selector)
	})

	b.Handle("/reset", func(c tele.Context) error {
		ChooseExercise = true
		expPerTraning = 0
		builder.Reset()
		menu.Selector.InlineKeyboard[0][0].Text = "Choose exercise"
		c.Send("<i>Reset compleated!</i>")
		return c.Send("<i>Start working out right now!</i>", menu.Selector)
	})

	b.Handle("/cmpl", func(c tele.Context) error {
		builder.WriteString(fmt.Sprintf(msg3, expPerTraning))
		c.Send(builder.String())

		// Update user exp and check current lvl
		exp := getExp()
		actualExp := exp + expPerTraning
		writeExp(actualExp)
		checkLevel(c, actualExp)

		// Reset to default settings
		ChooseExercise = true
		expPerTraning = 0
		builder.Reset()
		menu.Selector.InlineKeyboard[0][0].Text = "Choose exercise"
		return c.Send("<i>Start working out right now!</i>", menu.Selector)

	})

	b.Handle("/cl", func(c tele.Context) error {
		exp := getExp()
		checkLevel(c, exp)
		return c.Send(fmt.Sprintf("<i>Your exp:</i> <b>%v</b>", exp))

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

func checkLevel(c tele.Context, exp int64) {
	switch {
	case exp >= L10:
		c.Send(fmt.Sprintf(msg4, "Toji Fushiguro", exp))
	case exp >= L9:
		c.Send(fmt.Sprintf(msg2, "Saitama (One Punch Man)", "Zeno-sama (Dragon Ball)", L10-exp))
	case exp >= L8:
		c.Send(fmt.Sprintf(msg2, "Gojo Satoru (Jujutsu Kaisen)", "Saitama (One Punch Man)", L9-exp))
	case exp >= L7:
		c.Send(fmt.Sprintf(msg2, "Madara Uchiha (Naruto)", "Gojo Satoru (Jujutsu Kaisen)", L8-exp))
	case exp >= L6:
		c.Send(fmt.Sprintf(msg2, "Whitebeard (One Piece)", "Madara Uchiha (Naruto)", L7-exp))
	case exp >= L5:
		c.Send(fmt.Sprintf(msg2, "Escanor (Seven Deadly Sins)", "Whitebeard (One Piece)", L6-exp))
	case exp >= L4:
		c.Send(fmt.Sprintf(msg2, "All Might (My Hero Academia)", "Escanor (Seven Deadly Sins)", L5-exp))
	case exp >= L3:
		c.Send(fmt.Sprintf(msg2, "Levi Ackerman (Attack on Titan)", "All Might (My Hero Academia)", L4-exp))
	case exp >= L2:
		c.Send(fmt.Sprintf(msg2, "Trafalgar Law (One Piece)", "Levi Ackerman (Attack on Titan)", L3-exp))
	default:
		c.Send(fmt.Sprintf(msg2, "Killua Zoldyck (Hunter x Hunter)", "Trafalgar Law (One Piece)", L2-exp))
	}
}
