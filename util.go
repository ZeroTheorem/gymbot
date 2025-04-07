package main

import "strings"

func xpToNextLevel(level int64) int64 {
	nextLevelXP := (level + 1) * (level + 1) * 100
	return nextLevelXP
}
func generateProgressBar(percent int) string {
	completed := percent * 23 / 100
	bar := strings.Repeat("█", completed) + strings.Repeat("░", 23-completed)
	return bar
}

func getPercent(num1, num2 int64) float64 {
	return (float64(num1) / float64(num2)) * 100
}
