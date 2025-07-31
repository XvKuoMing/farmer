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
    

async def main():    
    async with async_playwright() as playwright:
        # Connect to remote Chrome instead of launching local browser
        browser = await playwright.chromium.launch(
            channel="chrome",
            headless=False,
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