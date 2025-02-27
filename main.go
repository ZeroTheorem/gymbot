package main

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

func main() {
	pref := tele.Settings{
		Token:  "TOKEN",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	mainMenu, createTraingBtn, createExerciseBtn := createMenuSelector()

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!", mainMenu)
	})
	b.Handle(createTraingBtn, func(c tele.Context) error {
		return c.Send("Traning")
	})
	b.Handle(createExerciseBtn, func(c tele.Context) error {
		return c.Send("Exercisej")
	})

	b.Start()
}
