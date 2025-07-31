#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Browser Automation Container${NC}"

# Function to check if a process is running
is_running() {
    pgrep -f "$1" > /dev/null 2>&1
}

# Function to wait for a service to be ready
wait_for_service() {
    local service_name="$1"
    local check_command="$2"
    local max_attempts=30
    local attempt=1
    
    echo -e "${YELLOW}⏳ Waiting for $service_name to start...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if eval "$check_command"; then
            echo -e "${GREEN}✅ $service_name is ready!${NC}"
            return 0
        fi
        
        echo -e "${YELLOW}   Attempt $attempt/$max_attempts - $service_name not ready yet...${NC}"
        sleep 1
        ((attempt++))
    done
    
    echo -e "${RED}❌ $service_name failed to start after $max_attempts attempts${NC}"
    return 1
}

# Start Xvfb (Virtual Frame Buffer)
echo -e "${BLUE}🖥️  Starting virtual display (Xvfb)...${NC}"
if ! is_running "Xvfb"; then
    Xvfb :99 -screen 0 1920x1080x24 -ac +extension GLX +render -noreset -nolisten tcp &
    XVFB_PID=$!
    echo -e "${GREEN}✅ Xvfb started with PID: $XVFB_PID${NC}"
else
    echo -e "${YELLOW}⚠️  Xvfb is already running${NC}"
fi

# Export display variable
export DISPLAY=:99
echo -e "${BLUE}📺 Display set to: $DISPLAY${NC}"

# Wait for Xvfb to be ready
wait_for_service "Xvfb" "xdpyinfo -display :99 >/dev/null 2>&1"

# Start VNC server if requested
if [ "${ENABLE_VNC:-true}" = "true" ]; then
    echo -e "${BLUE}🔗 Starting VNC server...${NC}"
    if ! is_running "x11vnc"; then
        x11vnc -display :99 -nopw -listen localhost -xkb -ncache 10 -ncache_cr -forever -shared -bg -o /var/log/x11vnc.log
        echo -e "${GREEN}✅ VNC server started on port 5900${NC}"
        echo -e "${BLUE}📋 Connect with VNC viewer to: localhost:5900${NC}"
    else
        echo -e "${YELLOW}⚠️  VNC server is already running${NC}"
    fi
else
    echo -e "${YELLOW}⏭️  VNC server disabled (set ENABLE_VNC=true to enable)${NC}"
fi

# Function to cleanup on exit
cleanup() {
    echo -e "${YELLOW}🧹 Cleaning up processes...${NC}"
    if [ ! -z "$XVFB_PID" ]; then
        kill $XVFB_PID 2>/dev/null || true
    fi
    pkill -f "x11vnc" 2>/dev/null || true
    echo -e "${GREEN}✅ Cleanup complete${NC}"
}

# Set trap for cleanup on exit
trap cleanup EXIT INT TERM

# Display system information
echo -e "${BLUE}🔍 System Information:${NC}"
echo -e "   Display: $DISPLAY"
echo -e "   Resolution: $(xdpyinfo -display :99 2>/dev/null | grep dimensions | awk '{print $2}' || echo 'Unknown')"
echo -e "   Python: $(python --version 2>&1)"
echo -e "   Working Directory: $(pwd)"

# Check if browsers are available
echo -e "${BLUE}🌐 Browser Check:${NC}"
if command -v chromium-browser >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Chromium browser available${NC}"
elif command -v google-chrome >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Google Chrome browser available${NC}"
else
    echo -e "${YELLOW}⚠️  No browser found, patchwright will download one${NC}"
fi

# Additional environment setup
echo -e "${BLUE}⚙️  Environment Setup:${NC}"
echo -e "   PYTHONUNBUFFERED: ${PYTHONUNBUFFERED:-1}"
echo -e "   ANONYMIZED_TELEMETRY: ${ANONYMIZED_TELEMETRY:-false}"

# Test the display
echo -e "${BLUE}🧪 Testing display...${NC}"
if xdpyinfo -display :99 >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Display test passed${NC}"
else
    echo -e "${RED}❌ Display test failed${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Container startup complete!${NC}"
echo -e "${BLUE}🚀 Starting application...${NC}"

# Execute the main command
exec "$@"