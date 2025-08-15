package devenv

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/s21platform/s21-cli/internal/devenv"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Запускает локальное окружение разработки",
		Long: `Команда start читает конфигурацию из devenv.toml файла и запускает
необходимые сервисы через docker-compose.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Проверяем наличие devenv.toml
			configPath := "devenv.toml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return fmt.Errorf("файл конфигурации %s не найден", configPath)
			}

			// Загружаем конфигурацию
			config, err := devenv.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("ошибка при чтении конфигурации: %v", err)
			}

			// Генерируем docker-compose.yml
			if err := devenv.GenerateCompose(config, "docker-compose.yml"); err != nil {
				return fmt.Errorf("ошибка при генерации docker-compose.yml: %v", err)
			}

			// Запускаем docker-compose
			dockerCmd := exec.Command("docker-compose", "up", "-d")
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				return fmt.Errorf("ошибка при запуске окружения: %v", err)
			}

			return nil
		},
	}
}
