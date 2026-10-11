<p align="center"><img src="build/appicon.png" width="112" height="112" alt="Sudoku VPN icon"></p>

# Sudoku VPN

A minimalist Sudoku VPN client for Windows 10/11 x64. Built with Go, Wails and Vue, with bundled Montserrat typography and light, dark and system themes.

**[Download the Windows installer](https://github.com/sanderstripa/sudoku-app/releases/latest/download/Sudoku-Setup.exe)** · **[Releases](https://github.com/sanderstripa/sudoku-app/releases)**

## Features

- Device-wide VPN with a single connection button.
- Domain-based split tunneling: route all traffic through VPN, or bypass VPN for listed domains and their subdomains.
- Import Sudoku connections and subscriptions exported by [Sudoku UI](https://github.com/sanderstripa/sudoku-ui).
- A themed tray menu for connections, profiles and settings.
- Optional launch with Windows and automatic connection on startup.
- A matching installer with light/dark themes and Russian/English switching.

## Getting started

1. Download and run **Sudoku-Setup.exe**. It includes the application and installs it into `Program Files\Sudoku`.
2. Open Sudoku, select **Add connection** and paste a `sudoku://` link or an HTTP(S) Sudoku subscription link.
3. Click the circular power button to connect. Automatic connection is optional and disabled by default.
4. To bypass selected websites, choose **All except the list** in Settings, enter one domain per line and save.

The application currently uses Russian interface labels. The installer supports Russian and English.

The interface runs with normal user permissions. Windows requests elevation for a separate network helper; manual administrator launch is unnecessary. Microsoft Edge WebView2 Runtime is required.

Closing the window hides it in the tray and keeps the connection active. **Exit** in the tray disconnects VPN and closes the application. Saving connection settings during an active connection triggers a reconnect. Theme changes apply immediately.

## Split tunneling

An entry such as `hh.ru` includes its subdomains, including `ekaterinburg.hh.ru`. Enter domains only, one per line. Application-based and arbitrary IP-based exclusions are not supported by this list.

Websites may use additional domains for authentication, APIs or media. Excluding a main domain does not exclude unrelated domains automatically. Existing browser connections can retain their previous route until the browser is restarted. A website's own access restrictions can still prevent it from opening; an exclusion does not guarantee acceptance by that website.

## Profiles and updates

Connections and settings are stored locally in `%APPDATA%\sudoku-app\config.json`. It contains connection credentials; do not share it publicly. Reimporting the same subscription updates its existing connection. Updates and uninstalling preserve profile data.

Only the Sudoku protocol is supported. VLESS, VMess and Trojan links are not supported. Sudoku, HEV and Wintun are bundled and extracted locally at startup. Updating the application updates the bundled core. In-app core updates and browser WebRTC controls are not implemented.

## Build and verification

Requirements: Windows x64, Go 1.26.4 or newer, Node.js 24 or newer, Git and internet access for dependencies.

```powershell
./scripts/build-windows.ps1
./scripts/build-installer.ps1 -SkipApp
```

Outputs: `build/bin/Sudoku.exe` and `build/bin/Sudoku-Setup.exe`. Use `-SkipCore` for an application build when the core has already been built.

```powershell
go test ./internal/core
cd frontend
npm ci
npx playwright install chromium
npm test
npm run build
```

GitHub Actions checks and builds the application and installer. Tagged releases publish both executables to GitHub Releases. Executables are not code-signed.

The build uses Sudoku `v0.5.1` with source patches, HEV `2.17.1`, Wails `v3.0.0-alpha.74` and frontend runtime `3.0.0-alpha.79`. Dependencies are recorded in `go.sum` and `frontend/package-lock.json`. Windows connections use TUN; the current client does not modify system proxy or WinHTTP settings.

## Credits and license

Based on [SUDOKU-ASCII/sudoku-desktop](https://github.com/SUDOKU-ASCII/sudoku-desktop), commit `b11605f5b8a205fdeea845ffc3db052a4c7a9f7e`. The original Go module path is retained for Wails bindings. Subscription compatibility follows [sanderstripa/sudoku-ui](https://github.com/sanderstripa/sudoku-ui), commit `d7695238975da95451903197b3fec4c20e622859`.

Licensed under **GPL-3.0-or-later**; see [LICENSE](LICENSE). Montserrat uses the **SIL Open Font License**; see [Montserrat-OFL.txt](frontend/src/assets/fonts/Montserrat-OFL.txt). Component notices are listed in [THIRD_PARTY.md](THIRD_PARTY.md).
