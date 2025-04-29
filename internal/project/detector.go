package project

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	// GoModFile - имя файла go.mod
	GoModFile = "go.mod"
	// S21ModulePrefix - префикс модуля S21
	S21ModulePrefix = "github.com/s21platform"
)

// IsS21Project проверяет, является ли текущая директория проектом S21
func IsS21Project() (bool, string, error) {
	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return false, "", err
	}

	// Проверка по имени модуля в go.mod
	projectRoot, found := findS21ModuleInParents(currentDir)
	if found {
		return true, projectRoot, nil
	}

	return false, "", nil
}

// findS21ModuleInParents рекурсивно ищет go.mod с префиксом S21 в родительских директориях
func findS21ModuleInParents(dir string) (string, bool) {
	// Проверяем наличие go.mod
	goModPath := filepath.Join(dir, GoModFile)
	if isS21GoModule(goModPath) {
		return dir, true
	}

	// Получаем родительскую директорию
	parentDir := filepath.Dir(dir)

	// Если мы достигли корня файловой системы, go.mod не найден
	if parentDir == dir {
		return "", false
	}

	// Проверяем родительскую директорию
	return findS21ModuleInParents(parentDir)
}

// isS21GoModule проверяет, содержит ли go.mod модуль S21
func isS21GoModule(goModPath string) bool {
	file, err := os.Open(goModPath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Ищем строку с объявлением модуля
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			// Проверяем, содержит ли имя модуля префикс S21
			moduleName := strings.TrimPrefix(line, "module ")
			return strings.Contains(moduleName, S21ModulePrefix)
		}
	}

	return false
}

// IsGitRepo проверяет, является ли директория Git-репозиторием
func IsGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return true
	}
	return false
}

// GetRemoteURL возвращает URL удаленного репозитория Git
// Можно использовать как дополнительную проверку
func GetRemoteURL(dir string) (string, error) {
	// Здесь будет код для получения URL удаленного репозитория
	// Для полной реализации нужно выполнить git команду
	// Это упрощенная заглушка
	return "", nil
}
