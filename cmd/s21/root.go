package s21

import (
	"fmt"
	"os"

	"github.com/s21platform/s21/cmd/codegen"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "s21",
	Aliases: []string{"s21"},
	Short:   "S21 CLI - утилита командной строки для разработчиков S21",
	Long:    `S21 CLI предоставляет инструменты для кодогенерации, линтинга и других задач разработки.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Если не указана подкоманда, показываем помощь
		cmd.Help()
	},
}

// Execute выполняет корневую команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// AddCommands добавляет все подкоманды к корневой команде
func AddCommands() {
	// Регистрируем команду кодогенерации
	rootCmd.AddCommand(codegen.NewCodegenCmd())

	// Регистрируем команду версии
	rootCmd.AddCommand(newVersionCmd())
}
