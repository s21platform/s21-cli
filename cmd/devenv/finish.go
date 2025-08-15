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
			// Проверяем наличие docker-compose.yml
			if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
				return fmt.Errorf("файл docker-compose.yml не найден")
			}

			// Выполняем docker-compose down
			dockerCmd := exec.Command("docker-compose", "down", "-v")
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				return fmt.Errorf("ошибка при остановке окружения: %v", err)
			}

			// Удаляем сгенерированный docker-compose.yml
			if err := os.Remove("docker-compose.yml"); err != nil {
				return fmt.Errorf("ошибка при удалении docker-compose.yml: %v", err)
			}

			return nil
		},
	}
}
