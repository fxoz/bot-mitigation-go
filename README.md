# Anti-Bot Reverse Proxy Server in Go

One of my first projects in Go: A reverse proxy server that detects bots and automated browsers. Inspired by *Cloudflare*, [*Anubis*](https://github.com/TecharoHQ/anubis) and - to some extend - by the [*GrimAC* anti-cheat plugin](https://github.com/GrimAnticheat/Grim).

## Features

- Very quick, automated bot detection (JavaScript-based)
- Configurable using yml.
- Designed to prevent even automated browsers with additional bot detection protections like `playwright_stealth` and `undetected-chromedriver`
- Written in Go to ensure high performance and scalability

### Planned

- Webkit support, additional Microsoft Edge testing
- Various checks for `playwright_stealth` and `undetected-chromedriver`
- VPN/Proxy/datacenter/TOR detection
  - Planned using ASN and IP databases
- Caching
- Ratelimits
- Anti-DDoS
- Manual CAPTCHA
- SEO-friendiness (dummy pages for search engines)
- Admin UI

## Protection

The goal is to detect all of the following methods, especially automated browsers with additional bot detection protection.

- **PW** - [Playwright](https://playwright.dev/python/)
- **UC** - [ultrafunkamsterdam/undetected-chromedriver](https://github.com/ultrafunkamsterdam/undetected-chromedriver)
- **PS** - [AtuboDad/playwright_stealth](https://github.com/AtuboDad/playwright_stealth)

| Method                      | Detected? | Comment                         |
| --------------------------- | --------- | ------------------------------- |
| Simple `curl` requests etc. | ✅         | Protection needs JS to function |
| PW: Chrome                  | ✅         | As of 2025-03-21                |
| UC                          | ❌         | As of 2025-03-21                |
| PS: Webkit                  | ❌         | As of 2025-03-21                |

## Browser Compatibility

It's really important to ensure that the bot protection doesn't break the website for legitimate users, even on older browsers. Please note that it's incredibly intricate to get the balance between security and compatibility right and that testing several browsers and their older versions takes a lot of time.

| Browser                               | Working? | Comment          |
| ------------------------------------- | -------- | ---------------- |
| Windows: Brave 1.76.80 (Chromium 134) | ✅        | As of 2025-03-21 |
| Windows: ungoogled-chromium 123       | ✅        | As of 2025-03-21 |
