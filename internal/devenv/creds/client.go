package creds

import (
	"context"
	"fmt"
	"time"

	"github.com/s21platform/creds/pkg/creds"
	"github.com/s21platform/s21-cli/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client представляет клиент для получения учетных данных
type Client struct {
	client creds.CredentialsServiceClient
	conn   *grpc.ClientConn
}

// NewClient создает новый клиент
func NewClient(target string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к сервису учетных данных: %v", err)
	}

	return &Client{
		client: creds.NewCredentialsServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

// GetEnvVars получает значения переменных окружения
func (c *Client) GetEnvVars(ctx context.Context, names []string) (map[string]string, error) {
	// Загружаем токен из конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке конфигурации: %v", err)
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("требуется авторизация. Выполните команду: s21 login")
	}

	// Создаем контекст с токеном
	ctx = metadata.AppendToOutgoingContext(ctx, "token", cfg.Token)

	fmt.Printf("Запрашиваем переменные: %v\n", names)
	resp, err := c.client.GetCreds(ctx, &creds.GetCredsRequest{
		Names: names,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении учетных данных: %v", err)
	}

	// Преобразуем ответ в map
	envVars := make(map[string]string)
	for _, cred := range resp.Credentials {
		envVars[cred.Name] = cred.Value
	}

	fmt.Printf("Получены переменные: %v\n", envVars)
	return envVars, nil
}
