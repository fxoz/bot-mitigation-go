import time
from playwright.sync_api import sync_playwright
from playwright_stealth import stealth_sync

with sync_playwright() as p:
    # browser_type = p.firefox
    browser_type = p.webkit
    # browser_type = p.chromium

    browser = browser_type.launch(
        headless=False, args=["--no-sandbox", "--disable-dev-shm-usage"]
    )
    page = browser.new_page()
    stealth_sync(page)

    page.goto("http://localhost:9977/")
    # page.goto("https://bot.sannysoft.com/")

    time.sleep(9999999)
    time.sleep(0.5)

    page.screenshot(path=f"testing/screenshots/{browser_type.name}.png")
    browser.close()
