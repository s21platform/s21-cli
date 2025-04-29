package s21

import (
	"fmt"

	"github.com/s21platform/s21-cli/internal/version"
	"github.com/spf13/cobra"
)

// newVersionCmd создает команду для отображения версии
func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Показать версию",
		Long:  "Показать информацию о версии S21 CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.GetFullVersionInfo())
		},
	}

	return cmd
}
