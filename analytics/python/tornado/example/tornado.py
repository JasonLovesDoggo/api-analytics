import os
# import sys
# sys.path.insert(0, os.path.abspath('../'))

import asyncio
from tornado.web import Application

from api_analytics.tornado import Analytics

from dotenv import load_dotenv

load_dotenv()


class MainHandler(Analytics):
    def __init__(self, app, res):
        api_key = os.environ.get("API_KEY")
        super().__init__(app, res, api_key)

    def get(self):
        self.write({'message': 'Hello World!'})


def make_app():
    return Application([
        (r"/", MainHandler),
    ])


async def main():
    app = make_app()
    app.listen(8080)
    await asyncio.Event().wait()

if __name__ == "__main__":
    asyncio.run(main())
