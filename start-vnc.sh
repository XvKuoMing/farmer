#!/bin/bash

# Kill any existing VNC servers
vncserver -kill $DISPLAY > /dev/null 2>&1 || true

# Wait a moment
sleep 2

# Start Xvfb
Xvfb $DISPLAY -screen 0 ${SCREEN_WIDTH}x${SCREEN_HEIGHT}x${SCREEN_DEPTH} -ac +extension GLX +render -noreset &

# Wait for X server to start
sleep 3

# Start window manager (fluxbox)
fluxbox &

# Start VNC server
x11vnc -display $DISPLAY -nopw -listen localhost -xkb -ncache 10 -ncache_cr -forever -shared