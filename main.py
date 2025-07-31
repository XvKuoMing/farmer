import asyncio
from browser_use import Agent, BrowserSession
from browser_use import ChatOpenAI
# from playwright.async_api import async_playwright
from patchright.async_api import async_playwright   # stealth alternative
import os

llm = ChatOpenAI(
    model="gpt-4o-mini",
    temperature=0,
    api_key=os.getenv("OPENAI_API_KEY"),
    base_url=os.getenv("OPENAI_BASE_URL"),
)

async def main():
    async with async_playwright() as playwright:
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
            task="собери данные: топ 30 вопросов про ислам. Верни ответ в формате: вопрос - ответ",
            llm=llm,
            browser_session=browser_session,
        )
        history = await agent.run()
        result = history.final_result()
        with open("islam_questions.txt", "w", encoding="utf-8") as f:
            f.write(result)

asyncio.run(main())