package devenv

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/s21platform/s21-cli/internal/devenv"
	"github.com/s21platform/s21-cli/internal/devenv/creds"
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

			// Генерируем .env.s21-cli файл
			envFile := ".env.s21-cli"
			if err := devenv.GenerateEnvFile(config, envFile); err != nil {
				return fmt.Errorf("ошибка при генерации %s файла: %v", envFile, err)
			}

			// Если включен сервис учетных данных, получаем переменные
			if config.Creds.Enabled {
				credsClient, err := creds.NewClient(config.Creds.Endpoint)
				if err != nil {
					return fmt.Errorf("ошибка при подключении к сервису учетных данных: %v", err)
				}
				defer credsClient.Close()

				// Собираем список необходимых переменных
				var envNames []string

				// Добавляем только переменные из секции [creds]
				for _, service := range config.Creds.Services {
					envNames = append(envNames, creds.GetServiceEnvNames(service)...)
				}

				// Получаем значения переменных
				envVars, err := credsClient.GetEnvVars(context.Background(), envNames)
				if err != nil {
					return fmt.Errorf("ошибка при получении учетных данных: %v", err)
				}

				// Добавляем полученные переменные в конфигурацию
				for k, v := range envVars {
					config.Env[k] = v
				}
			}

			// Генерируем .env.s21-cli файл с учетом полученных кредов
			if err := devenv.GenerateEnvFile(config, envFile); err != nil {
				return fmt.Errorf("ошибка при генерации %s файла: %v", envFile, err)
			}

			// Загружаем переменные окружения
			if err := devenv.LoadEnvFile(envFile); err != nil {
				return fmt.Errorf("ошибка при загрузке переменных окружения из %s: %v", envFile, err)
			}

			// Генерируем docker-compose.s21-cli.yml
			composeFile := "docker-compose.s21-cli.yml"
			if err := devenv.GenerateCompose(config, composeFile); err != nil {
				return fmt.Errorf("ошибка при генерации %s: %v", composeFile, err)
			}

			// Запускаем docker-compose с указанием файла конфигурации
			dockerCmd := exec.Command("docker-compose", "-f", composeFile, "up")
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				if strings.Contains(err.Error(), "yaml:") {
					return fmt.Errorf("ошибка в конфигурации %s. Проверьте, что все необходимые переменные окружения указаны в секции [env] вашего devenv.toml файла", composeFile)
				}
				return fmt.Errorf("ошибка при запуске окружения: %v", err)
			}

			fmt.Println("Окружение успешно запущено!")
			return nil
		},
	}
}
