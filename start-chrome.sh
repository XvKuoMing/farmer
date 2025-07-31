#!/bin/bash

# Wait for X server to be ready
export DISPLAY=:99
while ! xdpyinfo -display $DISPLAY >/dev/null 2>&1; do
    echo "Waiting for X server to start..."
    sleep 1
done

echo "X server is ready, starting Chrome..."

# Chrome arguments for CDP and GUI mode
CHROME_ARGS=(
    --no-sandbox
    --disable-dev-shm-usage
    --disable-gpu-sandbox
    --disable-software-rasterizer
    --remote-debugging-address=0.0.0.0
    --remote-debugging-port=9222
    --enable-remote-extensions
    --user-data-dir=/home/chrome/chrome-data
    --disable-background-timer-throttling
    --disable-backgrounding-occluded-windows
    --disable-renderer-backgrounding
    --disable-features=TranslateUI
    --disable-ipc-flooding-protection
    --enable-automation
    --password-store=basic
    --use-mock-keychain
    --no-first-run
    --no-default-browser-check
    --disable-default-apps
    --disable-popup-blocking
    --disable-translate
    --disable-background-networking
    --disable-background-timer-throttling
    --disable-backgrounding-occluded-windows
    --disable-renderer-backgrounding
    --disable-field-trial-config
    --disable-hang-monitor
    --disable-prompt-on-repost
    --disable-domain-reliability
    --disable-component-extensions-with-background-pages
    --disable-extensions
    --window-size=${SCREEN_WIDTH},${SCREEN_HEIGHT}
    --window-position=0,0
    # Reduce D-Bus and service-related errors
    --disable-features=VizDisplayCompositor
    --disable-logging
    --silent-debugger-extension-api
    --disable-web-security
    --disable-features=VizDisplayCompositor,VizHitTestSurfaceLayer
    --log-level=3
)

# Start Chrome with channel=chrome (using the installed Chrome stable)
exec /usr/bin/google-chrome-stable "${CHROME_ARGS[@]}" "$@"