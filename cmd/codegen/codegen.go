package codegen

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCodegenCmd создает команду для кодогенерации
func NewCodegenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codegen",
		Short: "Команды для кодогенерации",
		Long:  `Различные инструменты для генерации кода.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Команды кодогенерации будут добавлены позже")
		},
	}

	// Здесь будем добавлять подкоманды для разных типов кодогенерации
	// cmd.AddCommand(newModelCmd())
	// cmd.AddCommand(newControllerCmd())

	return cmd
}
