package platform

func SupportedPlatforms() []Platform {
	return []Platform{PlatformCopilotCLI, PlatformClaudeCode, PlatformCodexCLI}
}

func IsSupportedTarget(target string) bool {
	switch Platform(target) {
	case PlatformCopilotCLI, PlatformClaudeCode, PlatformCodexCLI:
		return true
	default:
		return false
	}
}
