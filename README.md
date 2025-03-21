# Anti-Bot Reverse Proxy Server in Go

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/fxoz/bot-mitigation-go/codeql.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/fxoz/bot-mitigation-go)
![GitHub License](https://img.shields.io/github/license/fxoz/bot-mitigation-go)

One of my first projects in Go: A reverse proxy server that detects bots and automated browsers. Inspired by *Cloudflare*, [*Anubis*](https://github.com/TecharoHQ/anubis) and - to some extend - by the [*GrimAC* anti-cheat plugin](https://github.com/GrimAnticheat/Grim).

With LLM scrapers gaining popularity, bot mitigation is more important than ever. This experimental project aims to provide a high-performance, scalable and easy-to-use solution for bot detection and mitigation.

## Features

- Very quick, automated bot detection (JavaScript-based)
- Configurable using yml.
- Designed to prevent even automated browsers with additional bot detection protections like `playwright_stealth` and `undetected-chromedriver`
- Written in Go to ensure high performance and scalability

### Planned

- Webkit support, additional Microsoft Edge testing
- Monitor input fields & scrolling
- Various checks for `playwright_stealth` and `undetected-chromedriver`
- VPN/Proxy/datacenter/TOR detection
  - Planned using ASN and IP databases
- SSL(?)
- Caching
- Ratelimits
- Anti-DDoS
- Manual CAPTCHA
  - Accessible (keyboard navigation, screen readers for visually impaired) but still robust; also low resource usage
- SEO-friendiness (dummy pages for search engines)
- Load balancing
- Admin UI

## Protection

The goal is to detect all of the following methods, especially automated browsers with additional bot detection protection.

- **PW** - [Playwright](https://playwright.dev/python/)
- **UCD** - [ultrafunkamsterdam/undetected-chromedriver](https://github.com/ultrafunkamsterdam/undetected-chromedriver)
- **PWS** - [AtuboDad/playwright_stealth](https://github.com/AtuboDad/playwright_stealth)

### Regular browsers

| Client                     | Detected? | Note                            |
| -------------------------- | --------- | ------------------------------- |
| Simple `curl` request etc. | ✅         | Protection needs JS to function |
| PW: Chrome                 | ✅         | As of 2025-03-21                |

### Browsers with evasion techniques

| Client       | Detected? | Note                                                                                                |
| ------------ | --------- | --------------------------------------------------------------------------------------------------- |
| PWS: Chrome  | ✅         | Fails [`navigator.webdriver`](https://developer.mozilla.org/en-US/docs/Web/API/Navigator/webdriver) |
| PWS: Firefox | ✅         | Fails [`navigator.userAgent`](https://caniuse.com/?search=navigator.userAgent)                      |
| PWS: WebKit  | ✅         | Fails [`navigator.userAgent`](https://caniuse.com/?search=navigator.userAgent)                      |
| UCD          | ❌         | As of 2025-03-21                                                                                    |

## Browser Compatibility

It's really important to ensure that the bot protection doesn't break the website for legitimate users, even on older browsers. Please note that it's incredibly intricate to get the balance between security and compatibility right and that testing several browsers and their older versions takes a lot of time.

| Browser                            | Passing? | Note             |
| ---------------------------------- | -------- | ---------------- |
| Windows: Brave 1.76.80 (C 134)     | ✅        | As of 2025-03-21 |
| Windows: Chrome 134                | ✅        | As of 2025-03-21 |
| Windows: ungoogled-chromium 123    | ✅        | As of 2025-03-21 |
| Windows: Firefox 136               | ✅        | As of 2025-03-21 |
| Android 14: Chrome 126             | ✅        | As of 2025-03-21 |
| Android 14: Brave 1.75.181 (C 133) | ✅        | As of 2025-03-21 |
| Android 14: Firefox 135            | ✅        | As of 2025-03-21 |

***

Developement of this project started on 2025-03-20.
