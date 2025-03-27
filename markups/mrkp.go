package markups

import (
	tele "gopkg.in/telebot.v4"
)

type SubMenu struct {
	Selector          *tele.ReplyMarkup
	ChooseExerciseBtn *tele.Btn
}

func CreateMenuSelector() *SubMenu {
	selector := &tele.ReplyMarkup{}
	chooseExerciseBtn := selector.Data("Choose exercise", "ChooseExercise")
	selector.Inline(
		selector.Row(chooseExerciseBtn),
	)
	m := SubMenu{
		Selector:          selector,
		ChooseExerciseBtn: &chooseExerciseBtn,
	}
	return &m
}
