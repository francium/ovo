#!/bin/env python3

import time
import signal
import sys

def on_signal(sig: int, frame):
    print("Caught signal:", sig)
    sys.exit(1)

# signal.SIGKILL can't be caught
signal.signal(signal.SIGINT, on_signal)
signal.signal(signal.SIGTERM, on_signal)

if __name__ == "__main__":
    print("Starting...")
    for i in range(1, 21):
        time.sleep(1)
        print(f"{i} second(s) have passed")
    print("Done")
