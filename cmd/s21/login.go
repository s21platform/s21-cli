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
		Long:  `Команда login сохраняет токен авторизации и nickname для доступа к сервисам S21.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)

			// Запрашиваем токен
			fmt.Print("Введите ваш токен: ")
			token, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("ошибка при чтении токена: %v", err)
			}
			token = strings.TrimSpace(token)
			if token == "" {
				return fmt.Errorf("токен не может быть пустым")
			}

			// Запрашиваем nickname
			fmt.Print("Введите ваш nickname: ")
			nickname, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("ошибка при чтении nickname: %v", err)
			}
			nickname = strings.TrimSpace(nickname)
			if nickname == "" {
				return fmt.Errorf("nickname не может быть пустым")
			}

			// Сохраняем конфигурацию
			cfg := &config.Config{
				Token:    token,
				Nickname: nickname,
			}
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("ошибка при сохранении конфигурации: %v", err)
			}

			fmt.Println("Токен и nickname успешно сохранены!")
			return nil
		},
	}
}
