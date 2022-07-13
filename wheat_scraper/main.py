from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
options = Options()
options.headless = True
options.add_argument("--windows-size=1920,1080")
options.add_argument("start-maximized")

driver = webdriver.Chrome(options=options, executable_path="C:/Users/PC/chrome_driver/chromedriver.exe")
driver.get("https://www.macrotrends.net/2534/wheat-prices-historical-chart-data")


table = driver.find_element(By.TAG_NAME, "table")
rows = table.find_elements(By.TAG_NAME, "tr")

for row in rows:

    print(row.text)

driver.quit()