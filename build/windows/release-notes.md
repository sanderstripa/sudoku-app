Windows 10/11 x64 · Sudoku 0.2.1

- Launch with Windows uses the standard per-user Run entry, verified after writing. Saving settings no longer depends on Task Scheduler registration.
- Legacy Sudoku scheduled tasks are removed in a bounded background migration.
- Startup reconciles the saved launch preference, including disabled autostart.
- Configuration writes use a temporary file and flush it before replacing the previous file.
- Regression checks cover an isolated Windows registry enable/disable round trip, quoted executable paths, saved settings after reload and replacing existing configuration. Application, installer and tray checks remain enabled.
- A real Windows sign-out/sign-in cycle remains a user validation step.
