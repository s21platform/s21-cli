package s21

import (
	"fmt"
	"os"

	"github.com/s21platform/s21-cli/cmd/codegen"
	"github.com/s21platform/s21-cli/internal/project"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "s21",
	Aliases: []string{"s21"},
	Short:   "S21 CLI - утилита командной строки для разработчиков S21",
	Long:    `S21 CLI предоставляет инструменты для кодогенерации, линтинга и других задач разработки.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Здесь можно добавить общие операции для всех команд
		// Например, инициализация логгера, проверка версии и т.д.

		// Пропускаем проверку проекта для определенных команд
		cmdName := cmd.Name()
		if cmdName == "version" || cmdName == "help" || cmdName == "completion" {
			return nil
		}

		// Проверяем, находится ли текущая директория в проекте S21
		isProject, projectRoot, err := project.IsS21Project()
		if err != nil {
			return fmt.Errorf("ошибка при определении проекта: %v", err)
		}

		// Если это не проект S21, выводим сообщение об ошибке и завершаем работу
		if !isProject {
			return fmt.Errorf("эта команда должна выполняться в директории проекта S21.\nПроверьте, что путь модуля в go.mod содержит github.com/s21platform")
		}

		// Устанавливаем каталог проекта как рабочую директорию
		if err := os.Chdir(projectRoot); err != nil {
			return fmt.Errorf("не удалось перейти в корень проекта: %v", err)
		}

		return nil
	},
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
