package devenv

import (
	"github.com/spf13/cobra"
)

// NewDevenvCmd создает новую команду devenv
func NewDevenvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devenv",
		Short: "Управление локальным окружением разработки",
		Long: `Команда devenv позволяет управлять локальным окружением разработки.
Она предоставляет возможность запускать и останавливать необходимые сервисы (PostgreSQL, Redis, Redpanda)
на основе конфигурации из devenv.toml файла.`,
	}

	// Добавляем подкоманды
	cmd.AddCommand(newStartCmd())
	cmd.AddCommand(newFinishCmd())
	cmd.AddCommand(newLoadCmd())

	return cmd
}
