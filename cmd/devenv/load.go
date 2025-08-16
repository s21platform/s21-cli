package devenv

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

func newLoadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "load",
		Short: "Загружает переменные окружения из .env.s21-cli",
		Long: `Команда load загружает переменные окружения из .env.s21-cli файла и запускает новый shell.
После выхода из shell переменные окружения вернутся в исходное состояние.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile := ".env.s21-cli"

			// Проверяем наличие файла
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				return fmt.Errorf("файл %s не найден. Сначала выполните команду 's21 devenv start'", envFile)
			}

			// Читаем файл
			file, err := os.Open(envFile)
			if err != nil {
				return fmt.Errorf("ошибка при чтении файла %s: %v", envFile, err)
			}
			defer file.Close()

			// Получаем текущие переменные окружения
			env := os.Environ()

			// Читаем и добавляем переменные из файла
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				env = append(env, line)
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("ошибка при чтении файла %s: %v", envFile, err)
			}

			// Определяем shell пользователя
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "/bin/sh"
			}

			// Запускаем shell с новыми переменными
			fmt.Printf("Загружены переменные окружения из %s\n", envFile)
			fmt.Println("Для выхода используйте 'exit' или Ctrl+D")

			shellCmd := exec.Command(shell)
			shellCmd.Env = env
			shellCmd.Stdin = os.Stdin
			shellCmd.Stdout = os.Stdout
			shellCmd.Stderr = os.Stderr

			// Используем syscall.Exec для замены текущего процесса
			return syscall.Exec(shell, []string{shell}, env)
		},
	}
}
