from playwright.sync_api import sync_playwright
import os
import time

with sync_playwright() as p:
    browser = p.chromium.launch(headless=False)
    page = browser.new_page()
    # page.goto(f"file:///{os.path.abspath(r'testing/_index.html')}")
    page.goto(r"http://localhost:9977")
    # page.goto(r"https://bot.sannysoft.com/")

    time.sleep(999999)
