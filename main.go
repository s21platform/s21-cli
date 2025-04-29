package main

import "github.com/s21platform/s21-cli/cmd/s21"

func main() {
	// Добавляем все команды к корневой команде
	s21.AddCommands()

	// Запускаем CLI
	s21.Execute()
}
