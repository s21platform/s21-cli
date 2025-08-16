package creds

import "strings"

// GetServiceEnvNames возвращает список имен переменных окружения для сервиса
func GetServiceEnvNames(serviceName string) []string {
	// Преобразуем имя сервиса в верхний регистр
	return []string{strings.ToUpper(serviceName)}
}
