package markups

import (
	tele "gopkg.in/telebot.v4"
)

type MainMenu struct {
	Selector         *tele.ReplyMarkup
	CreateTraningBtn *tele.Btn
}

func CreateMainMenuSelector() *MainMenu {
	selector := &tele.ReplyMarkup{}
	createTraningBtn := selector.Data("Начать тренировку", "createTraning")
	selector.Inline(
		selector.Row(createTraningBtn),
	)
	m := MainMenu{
		Selector:         selector,
		CreateTraningBtn: &createTraningBtn,
	}
	return &m
}
