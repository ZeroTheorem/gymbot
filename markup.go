package main

import (
	tele "gopkg.in/telebot.v4"
)

func createMenuSelector() (*tele.ReplyMarkup, *tele.Btn, *tele.Btn) {
	selector := &tele.ReplyMarkup{}
	createTraningBtn := selector.Data("Создать тренировку", "create")
	createExerciseBtn := selector.Data("Создать тренировку", "create")
	selector.Inline(
		selector.Row(createTraningBtn),
		selector.Row(createExerciseBtn),
	)
	return selector, &createTraningBtn, &createExerciseBtn

}
