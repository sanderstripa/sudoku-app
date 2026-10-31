Windows 10/11 x64 · Sudoku 0.2.2

- Update protocol core and link/config parsing to upstream Sudoku v0.5.2-1.
- Include upstream SOCKS5 UDP loopback/IPv6 fixes and disable MPTCP by default.
- Port desktop traffic counters, DIRECT exceptions and TUN bypass to the new client architecture.
- Keep upstream HTTP CONNECT buffering, UDP handling, existing subscriptions, credentials and saved settings.
- Build with Go 1.27.3; test upstream client regressions and desktop settings/import behavior before packaging.
