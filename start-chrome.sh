#!/bin/bash

# Wait for X server to be ready
export DISPLAY=:99
while ! xdpyinfo -display $DISPLAY >/dev/null 2>&1; do
    echo "Waiting for X server to start..."
    sleep 1
done

echo "X server is ready, starting Chrome..."

# Chrome arguments for CDP and GUI mode with enhanced stealth
CHROME_ARGS=(
    # Basic Chrome setup
    --no-sandbox
    --disable-dev-shm-usage
    --disable-gpu-sandbox
    --disable-software-rasterizer
    --user-data-dir=/home/chrome/chrome-data
    --remote-debugging-port=9223
    --remote-debugging-address=127.0.0.1
    --remote-allow-origins=*
    --enable-remote-extensions
    
    # Core stealth - disable automation detection
    --disable-blink-features=AutomationControlled
    --disable-automation
    --exclude-switches=enable-automation
    --enable-automation=false
    --disable-infobars
    
    # Window and display
    --window-size=${SCREEN_WIDTH},${SCREEN_HEIGHT}
    --window-position=0,0
    --start-maximized
    
    # Disable telemetry and reporting
    --disable-background-timer-throttling
    --disable-backgrounding-occluded-windows
    --disable-renderer-backgrounding
    --disable-background-networking
    --disable-sync
    --disable-translate
    --disable-background-downloads
    --disable-add-to-shelf
    --disable-client-side-phishing-detection
    --disable-datasaver-prompt
    --disable-domain-reliability
    --disable-component-update
    --disable-component-extensions-with-background-pages
    --disable-default-apps
    --disable-extensions
    --disable-extensions-file-access-check
    --disable-extensions-http-throttling
    
    # Performance and memory
    --memory-pressure-off
    --max_old_space_size=4096
    --aggressive-cache-discard
    
    # Disable various web APIs for detection evasion
    --disable-features=VizDisplayCompositor,AudioServiceOutOfProcess,UserAgentClientHint,VizHitTestSurfaceLayer,TranslateUI
    
    # First run and defaults
    --no-first-run
    --no-default-browser-check
    --no-service-autorun
    --password-store=basic
    --use-mock-keychain
    
    # Permissions and prompts
    --deny-permission-prompts
    --disable-notifications
    --disable-geolocation
    --disable-popup-blocking
    --disable-save-password-bubble
    
    # Security (reduces detection but may affect functionality)
    --disable-web-security
    --allow-running-insecure-content
    --ignore-certificate-errors
    --ignore-ssl-errors
    --ignore-certificate-errors-spki-list
    
    # Behavior modifications
    --disable-field-trial-config
    --disable-hang-monitor
    --disable-prompt-on-repost
    --disable-ipc-flooding-protection
    --autoplay-policy=no-user-gesture-required
    
    # Proxy detection evasion
    --proxy-bypass-list=*
    --proxy-server="direct://"
    
    # Logging and debugging
    --disable-logging
    --silent-debugger-extension-api
    --log-level=3
)

# Start Chrome on localhost:9223
echo "Starting Chrome on localhost:9223..."
/usr/bin/google-chrome-stable "${CHROME_ARGS[@]}" "$@" &

# Start ncat proxy to make Chrome accessible on all interfaces
echo "Starting ncat proxy on 0.0.0.0:9222 -> localhost:9223..."
ncat \
    --sh-exec "ncat localhost 9223" \
    -l 9222 \
    --keep-open &

echo "Chrome DevTools accessible on port 9222"

# Wait for processes
wait