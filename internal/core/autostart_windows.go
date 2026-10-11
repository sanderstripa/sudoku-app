//go:build windows

package core

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

const (
	windowsRunKeyPath       = `Software\Microsoft\Windows\CurrentVersion\Run`
	windowsRunValueName     = "sudoku-app"
	windowsAutostartArgLine = "--autostart"
)

func setLaunchAtLogin(enabled bool) error {
	exe, err := currentExecutablePath()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", windowsAutostartScript(exe, enabled))
	applyManagedProcessSysProcAttr(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Не удалось изменить автозапуск Windows: %w: %s", err, strings.TrimSpace(string(output)))
	}
	// Remove only this app's legacy Run entry after the task is registered/removed.
	key, err := registry.OpenKey(registry.CURRENT_USER, windowsRunKeyPath, registry.SET_VALUE)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open legacy autostart key: %w", err)
	}
	defer key.Close()
	err = key.DeleteValue(windowsRunValueName)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	return err
}

// InteractiveToken + Highest runs in the user's desktop session, without storing
// a password. HKCU Run entries cannot reliably launch a requireAdministrator exe.
func windowsAutostartScript(exe string, enabled bool) string {
	lines := []string{
		"$ErrorActionPreference = 'Stop'",
		"$identity = [Security.Principal.WindowsIdentity]::GetCurrent()",
		"$taskName = 'Sudoku-App-' + $identity.User.Value",
	}
	if !enabled {
		return strings.Join(append(lines,
			"if (Get-ScheduledTask -TaskName $taskName -ErrorAction SilentlyContinue) { Unregister-ScheduledTask -TaskName $taskName -Confirm:$false }"), "\n")
	}
	return strings.Join(append(lines,
		"$action = New-ScheduledTaskAction -Execute '"+strings.ReplaceAll(exe, "'", "''")+"' -Argument '"+windowsAutostartArgLine+"'",
		"$trigger = New-ScheduledTaskTrigger -AtLogOn -User $identity.Name",
		"$principal = New-ScheduledTaskPrincipal -UserId $identity.Name -LogonType Interactive -RunLevel Highest",
		"$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -ExecutionTimeLimit ([TimeSpan]::Zero) -MultipleInstances IgnoreNew",
		"Register-ScheduledTask -TaskName $taskName -Action $action -Trigger $trigger -Principal $principal -Settings $settings -Force | Out-Null"), "\n")
}
