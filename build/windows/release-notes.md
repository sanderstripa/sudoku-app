Windows 10/11 x64 · 0.1.9

- Left-clicking the tray icon opens and focuses the application. Right-clicking opens the themed tray menu.
- Domain exclusions now retain short-lived associations for real IPv4 addresses learned during direct DNS resolution. Later IP-only requests can reuse the matching domain exclusion.
- Associations expire and are checked against active rules. Unrelated addresses remain proxied. Shared addresses may follow the associated excluded domain during this short interval.
- Regression checks cover both resolved Ozon addresses, expiry, removed exceptions and unrelated IPs. These checks do not verify a signed-in Ozon browser session; website access still requires user validation.
- For Ozon, list ozon.ru, ozone.ru and ozonusercontent.com as separate exclusions. Save and reconnect after updating.
