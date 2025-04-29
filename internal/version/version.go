package version

// Версия CLI утилиты
var (
	Version   = "0.1.0"
	GitCommit = "development"
	BuildDate = "unknown"
)

// GetVersion возвращает текущую версию приложения
func GetVersion() string {
	return Version
}

// GetFullVersionInfo возвращает полную информацию о версии
func GetFullVersionInfo() string {
	return "Version: " + Version + "\n" +
		"Git commit: " + GitCommit + "\n" +
		"Build date: " + BuildDate
}
