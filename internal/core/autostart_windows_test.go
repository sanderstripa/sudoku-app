//go:build windows

package core

import (
	"strings"
	"testing"
)

func TestAutostartUsesNormalInteractiveUserAndQuotesExecutable(t *testing.T) {
	script := windowsAutostartScript(`C:\Users\O'Neil\My VPN\Sudoku.exe`, true)
	for _, required := range []string{`-Execute 'C:\Users\O''Neil\My VPN\Sudoku.exe'`, `-LogonType Interactive -RunLevel Limited`, `-AtLogOn -User $identity.Name`, `$identity.User.Value`, `-ExecutionTimeLimit ([TimeSpan]::Zero)`, `-MultipleInstances IgnoreNew`} {
		if !strings.Contains(script, required) {
			t.Fatalf("missing %q", required)
		}
	}
	disabled := windowsAutostartScript("ignored", false)
	if !strings.Contains(disabled, "Unregister-ScheduledTask") || strings.Contains(disabled, "Register-ScheduledTask ") || !strings.Contains(disabled, "$identity.User.Value") {
		t.Fatal("disable must only remove the current user's task")
	}
}
