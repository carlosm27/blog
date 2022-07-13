from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
import pandas as pd

options = Options()
options.headless = True

driver = webdriver.Chrome(options=options, executable_path="C:/Users/PC/chrome_driver/chromedriver.exe")

driver.get("https://www.macrotrends.net/2534/wheat-prices-historical-chart-data")

rows = driver.find_elements(By.TAG_NAME, "tr")


list_rows =[]

for row in rows:

    list_rows.append(str(row.text).replace(" ", ","))

    df = pd.DataFrame(list_rows)
    df.to_csv('wheat_prices.csv')


driver.quit()