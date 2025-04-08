package main

import (
	tele "gopkg.in/telebot.v4"
)

type menu struct {
	Selector          *tele.ReplyMarkup
	ChooseExerciseBtn *tele.Btn
}

func createMenuSelector() *menu {
	selector := &tele.ReplyMarkup{}
	chooseExerciseBtn := selector.Data("Choose exercise", "ChooseExercise")
	selector.Inline(
		selector.Row(chooseExerciseBtn),
	)
	m := menu{
		Selector:          selector,
		ChooseExerciseBtn: &chooseExerciseBtn,
	}
	return &m
}
