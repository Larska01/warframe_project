import schedule
import api_calls
import time
from datetime import datetime 

def job():
    '''job
    Function that updates the warframe api cache
    '''
    print(f"{datetime.now()} Updating cache")
    api_calls.get_all_data()

schedule.every(1).minutes.do(job)

while True:
    schedule.run_pending()
    time.sleep(1)