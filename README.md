# Anti-Bot Reverse Proxy Server in Go

One of my first projects in Go.

## Features

## Protection

The goal is to detect all of the following methods, especially automated browsers with additional bot detection protection.

- [ultrafunkamsterdam/undetected-chromedriver](https://github.com/ultrafunkamsterdam/undetected-chromedriver)
- [AtuboDad/playwright_stealth](https://github.com/AtuboDad/playwright_stealth)

| Method                    | Detected? | Comment                         |
| ------------------------- | --------- | ------------------------------- |
| Simple curl requests etc. | ✅         | Protection needs JS to function |
| Playwright: Chrome        | ✅         | As of 2025-03-21                |
| Undetected Chromedriver   | ❌         | As of 2025-03-21                |
| `playwright_stealth`      | ❌         | As of 2025-03-21                |

## Browser Compatibility

It's really important to ensure that the bot protection doesn't break the website for legitimate users, even on older browsers. Please note that it's incredibly intricate to get the balance between security and compatibility right and that testing several browsers and their older versions takes a lot of time.

| Browser                               | Working? | Comment          |
| ------------------------------------- | -------- | ---------------- |
| Windows: Brave 1.76.80 (Chromium 134) | ✅        | As of 2025-03-21 |
| Windows: ungoogled-chromium 123       | ✅        | As of 2025-03-21 |
