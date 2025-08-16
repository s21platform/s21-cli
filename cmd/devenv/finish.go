package devenv

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newFinishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "finish",
		Short: "Останавливает локальное окружение разработки",
		Long: `Команда finish останавливает все сервисы, запущенные через docker-compose,
и удаляет созданные контейнеры.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeFile := "docker-compose.s21-cli.yml"

			// Проверяем наличие docker-compose.yml
			if _, err := os.Stat(composeFile); os.IsNotExist(err) {
				return fmt.Errorf("файл %s не найден", composeFile)
			}

			// Выполняем docker-compose down
			dockerCmd := exec.Command("docker-compose", "-f", composeFile, "down", "-v")
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				return fmt.Errorf("ошибка при остановке окружения: %v", err)
			}

			// Удаляем сгенерированные файлы
			files := []string{composeFile, ".env.s21-cli"}
			for _, file := range files {
				if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("ошибка при удалении файла %s: %v", file, err)
				}
			}

			return nil
		},
	}
}
