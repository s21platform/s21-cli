package s21

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/s21platform/s21-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Авторизация в CLI",
		Long:  `Команда login сохраняет токен авторизации для доступа к сервисам S21.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Запрашиваем токен у пользователя
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Введите ваш токен: ")
			token, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("ошибка при чтении токена: %v", err)
			}

			// Очищаем токен от пробелов и переносов строк
			token = strings.TrimSpace(token)
			if token == "" {
				return fmt.Errorf("токен не может быть пустым")
			}

			// Сохраняем токен в конфигурацию
			cfg := &config.Config{
				Token: token,
			}
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("ошибка при сохранении токена: %v", err)
			}

			fmt.Println("Токен успешно сохранен!")
			return nil
		},
	}
}
