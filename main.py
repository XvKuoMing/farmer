import asyncio
import os
import time
from browser_use import Agent, BrowserSession
from browser_use import ChatOpenAI
# from playwright.async_api import async_playwright
from patchright.async_api import async_playwright   # stealth alternative


temp = 0.0
api_key = os.getenv("OPENAI_API_KEY")
base_url = os.getenv("OPENAI_BASE_URL") or None

llm = ChatOpenAI(
    model="gpt-4o-mini",
    temperature=temp,
    api_key=api_key,
    base_url=base_url
)

planner_llm = ChatOpenAI(
    model='gpt-4o-mini',
    temperature=temp,
    api_key=api_key,
    base_url=base_url
)
    

async def wait_for_chrome():
    """Wait for Chrome container to be ready"""
    chrome_endpoint = os.getenv("CHROME_WS_ENDPOINT", "ws://chrome:9222")
    max_retries = 30
    retry_delay = 2
    
    print(f"Waiting for Chrome at {chrome_endpoint}...")
    
    for i in range(max_retries):
        try:
            async with async_playwright() as playwright:
                # Try to connect to Chrome
                browser = await playwright.chromium.connect_over_cdp(
                    endpoint_url=chrome_endpoint
                )
                print("Successfully connected to Chrome!")
                await browser.close()
                return True
        except Exception as e:
            print(f"Attempt {i+1}/{max_retries}: Chrome not ready yet ({e})")
            if i < max_retries - 1:
                await asyncio.sleep(retry_delay)
    
    raise Exception("Chrome container is not available after maximum retries")

async def main():
    # Wait for Chrome container to be ready
    await wait_for_chrome()
    
    async with async_playwright() as playwright:
        # Connect to remote Chrome instead of launching local browser
        chrome_endpoint = os.getenv("CHROME_WS_ENDPOINT", "ws://chrome:9222")
        print(f"Connecting to Chrome at {chrome_endpoint}")
        
        browser = await playwright.chromium.connect_over_cdp(
            endpoint_url=chrome_endpoint
        )
        
        context = await browser.new_context()
        page = await context.new_page()

        browser_session = BrowserSession(
            page=page,
            browser_context=context,  # all these are supported
            browser=browser,
            playwright=playwright,
        )

        agent = Agent(
            task="Открой сайт myspar.ru и собери мне топ 10 сыров",
            llm=llm,
            planner_llm=planner_llm,
            browser_session=browser_session,
        )
        history = await agent.run()
        result = history.final_result()
        with open("result.txt", "w", encoding="utf-8") as f:
            f.write(result)

asyncio.run(main())