#!/bin/bash

# Exit on any error
set -e

# Function to cleanup on exit
cleanup() {
    echo "Cleaning up..."
    jobs -p | xargs -r kill
}
trap cleanup EXIT

# Start Xvfb (Virtual display)
echo "Starting Xvfb..."
Xvfb :99 -screen 0 1920x1080x24 -ac +extension GLX +render -noreset &
XVFB_PID=$!
export DISPLAY=:99

# Wait for Xvfb to start
sleep 3

# Check if Xvfb is running
if ! ps -p $XVFB_PID > /dev/null; then
    echo "ERROR: Xvfb failed to start"
    exit 1
fi

# Start window manager (fluxbox) for a proper desktop environment
echo "Starting window manager..."
fluxbox > /dev/null 2>&1 &
FLUXBOX_PID=$!

# Start VNC server for remote viewing (optional)
if [ "${ENABLE_VNC:-true}" = "true" ]; then
    echo "Starting VNC server..."
    x11vnc -display :99 -nopw -listen 0.0.0.0 -xkb -ncache 10 -ncache_cr -forever -shared -quiet &
    VNC_PID=$!
fi

# Wait for desktop environment to settle
sleep 2

# Create and setup Chrome user data directory
echo "Setting up Chrome user data directory..."
CHROME_DATA_DIR="/home/chrome/.config/google-chrome"
rm -rf "$CHROME_DATA_DIR"
mkdir -p "$CHROME_DATA_DIR"
chmod 755 "$CHROME_DATA_DIR"

# Create a basic Chrome preferences file to avoid first-run dialogs
cat > "$CHROME_DATA_DIR/First Run" << EOF
EOF

# Start Chrome with remote debugging
echo "Starting Chrome with remote debugging..."
echo "Chrome will be available for remote debugging on port 9222"
echo "VNC server available on port 5900 (if enabled)"

# Start Google Chrome with comprehensive flags for container environment
google-chrome \
    --no-sandbox \
    --disable-dev-shm-usage \
    --disable-gpu \
    --disable-software-rasterizer \
    --disable-background-timer-throttling \
    --disable-backgrounding-occluded-windows \
    --disable-renderer-backgrounding \
    --disable-features=TranslateUI,VizDisplayCompositor \
    --disable-extensions \
    --disable-plugins \
    --disable-sync \
    --disable-default-apps \
    --disable-background-networking \
    --disable-background-mode \
    --disable-client-side-phishing-detection \
    --disable-component-extensions-with-background-pages \
    --disable-component-update \
    --disable-domain-reliability \
    --disable-ipc-flooding-protection \
    --disable-prompt-on-repost \
    --disable-hang-monitor \
    --disable-web-security \
    --disable-features=VizDisplayCompositor \
    --remote-debugging-address=0.0.0.0 \
    --remote-debugging-port=9222 \
    --user-data-dir="$CHROME_DATA_DIR" \
    --window-size=1920,1080 \
    --start-maximized \
    --no-first-run \
    --no-default-browser-check \
    --disable-infobars \
    --disable-session-crashed-bubble \
    --disable-restore-session-state \
    about:blank &

CHROME_PID=$!

# Wait for Chrome to start and begin listening on port 9222
echo "Waiting for Chrome to start remote debugging server..."
for i in {1..30}; do
    if netstat -tuln 2>/dev/null | grep -q ":9222 "; then
        echo "Chrome remote debugging server is ready on port 9222!"
        break
    elif [ $i -eq 30 ]; then
        echo "ERROR: Chrome remote debugging server failed to start after 30 seconds"
        ps aux | grep chrome
        netstat -tuln | grep 9222 || echo "Port 9222 not listening"
        exit 1
    else
        echo "Attempt $i/30: Waiting for Chrome remote debugging server..."
        sleep 1
    fi
done

# Keep the script running and monitor Chrome
echo "Chrome container is ready!"
echo "Remote debugging: http://localhost:9222"
echo "VNC access: localhost:5900"

# Monitor Chrome process and restart if it crashes
while true; do
    if ! ps -p $CHROME_PID > /dev/null; then
        echo "WARNING: Chrome process died, restarting..."
        google-chrome \
            --no-sandbox \
            --disable-dev-shm-usage \
            --disable-gpu \
            --disable-software-rasterizer \
            --disable-background-timer-throttling \
            --disable-backgrounding-occluded-windows \
            --disable-renderer-backgrounding \
            --disable-features=TranslateUI,VizDisplayCompositor \
            --disable-extensions \
            --disable-plugins \
            --disable-sync \
            --disable-default-apps \
            --disable-background-networking \
            --disable-background-mode \
            --disable-client-side-phishing-detection \
            --disable-component-extensions-with-background-pages \
            --disable-component-update \
            --disable-domain-reliability \
            --disable-ipc-flooding-protection \
            --disable-prompt-on-repost \
            --disable-hang-monitor \
            --disable-web-security \
            --disable-features=VizDisplayCompositor \
            --remote-debugging-address=0.0.0.0 \
            --remote-debugging-port=9222 \
            --user-data-dir="$CHROME_DATA_DIR" \
            --window-size=1920,1080 \
            --start-maximized \
            --no-first-run \
            --no-default-browser-check \
            --disable-infobars \
            --disable-session-crashed-bubble \
            --disable-restore-session-state \
            about:blank &
        CHROME_PID=$!
    fi
    sleep 5
done