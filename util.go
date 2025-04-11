package main

import (
	"fmt"
	"strings"
)

func xpToNextLevel(level int64) int64 {
	nextLevelXP := (level + 1) * (level + 1) * 10
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

func resetState(menu *menu) {
	sessionExp = 0
	currentExercise = ""
	isChoosing = false
	isTraining = false
	builder.Reset()
	menu.Selector.InlineKeyboard[0][0].Text = "Choose exercise"
}

func updateLevel(lvl, xp, gainedXP int64) (int64, int64, []string) {
	xp += gainedXP
	messages := []string{}
	for xp >= xpToNextLevel(lvl) {
		xp -= xpToNextLevel(lvl)
		lvl++
		messages = append(messages,
			fmt.Sprintf("🎉 <i>Congratulations, you've reached a new level</i>: <b>%v</b>", lvl))
	}
	return lvl, xp, messages
}

func defineRank(lvl int64) string {
	switch {
	case lvl >= 300:
		return "<i>Rank:</i> <b>S+</b>"
	case lvl >= 250:
		return "<i>Rank:</i> <b>S</b> <i>(next rank</i> <b>S+</b> <i>on level</i> <b>300</b><i>)</i>"
	case lvl >= 200:
		return "<i>Rank:</i> <b>A</b> <i>(next rank</i> <b>S</b> <i>on level</i> <b>250</b><i>)</i>"
	case lvl >= 150:
		return "<i>Rank:</i> <b>B</b> <i>(next rank</i> <b>A</b> <i>on level</i> <b>200</b><i>)</i>"
	case lvl >= 100:
		return "<i>Rank:</i> <b>C</b> <i>(next rank</i> <b>B</b> <i>on level</i> <b>150</b><i>)</i>"
	case lvl >= 50:
		return "<i>Rank:</i> <b>D</b> <i>(next rank</i> <b>C</b> <i>on level</i> <b>100</b><i>)</i>"
	default:
		return "<i>Rank:</i> <b>E</b> <i>(next rank</i> <b>D</b> <i>on level</i> <b>50</b><i>)</i>"
	}
}
