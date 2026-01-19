package service

type Services struct {
	Auth     *AuthService
	Users    *UserService
	Scripts  *ScriptService
	Tasks    *TaskService
	Reports  *ReportService
	Runner   *TaskRunner
	Settings *SettingsService
	Stats    *StatsService
}
