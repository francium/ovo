#!/bin/env python3

import atexit
from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"

def on_exit():
    print("At exist signal received")

if __name__ == "__main__":
    atexit.register(on_exit)
    app.run()
