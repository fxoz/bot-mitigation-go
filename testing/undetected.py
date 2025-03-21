# undetected chromedriver
import time
import undetected_chromedriver as uc
from selenium import webdriver

options = webdriver.ChromeOptions()
# options.add_argument('--headless')
options.add_argument("--no-sandbox")
options.add_argument("--disable-dev-shm-usage")

driver = uc.Chrome(options=options)
driver.get("http://localhost:9977")
# driver.get("https://bot.sannysoft.com/")

time.sleep(999999)

# 36 - Audio context
# sampleRate : 96000
# state : running
