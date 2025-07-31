# Chrome with noVNC Docker Container

This Docker container runs Google Chrome with noVNC, allowing you to access Chrome through a web browser and connect via Chrome DevTools Protocol (CDP).

## Features

- ✅ Chrome browser (stable channel) with GUI (headless=false)
- ✅ noVNC web interface for browser access
- ✅ Chrome DevTools Protocol (CDP) access on port 9222
- ✅ Persistent user data and downloads
- ✅ VNC direct connection support
- ✅ Configurable screen resolution

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

### Using Docker directly

```bash
# Build the image
./build.sh

# Run the container
docker run -d \
  --name chrome-novnc \
  -p 6080:6080 \
  -p 9222:9222 \
  -p 5901:5901 \
  -v chrome-data:/home/chrome/.config/google-chrome \
  -v chrome-downloads:/home/chrome/Downloads \
  chrome-novnc
```

## Access Points

- **noVNC Web Interface**: http://localhost:6080
- **Chrome DevTools**: http://localhost:9222
- **VNC Direct**: localhost:5901 (no password required)

## Chrome DevTools Protocol (CDP) Usage

You can connect to Chrome via CDP using various tools:

### Python Example with pychrome

```python
import pychrome

# Connect to Chrome
browser = pychrome.Browser(url="http://localhost:9222")

# Create a new tab
tab = browser.new_tab()

# Navigate to a page
tab.Page.navigate(url="https://example.com")

# Wait for page to load
tab.Page.loadEventFired()

# Get page content
result = tab.Runtime.evaluate(expression="document.title")
print(result['result']['value'])
```

### JavaScript Example with puppeteer

```javascript
const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.connect({
    browserURL: 'http://localhost:9222'
  });
  
  const page = await browser.newPage();
  await page.goto('https://example.com');
  
  const title = await page.title();
  console.log('Page title:', title);
  
  await browser.disconnect();
})();
```

## Configuration

### Environment Variables

- `SCREEN_WIDTH`: Screen width (default: 1920)
- `SCREEN_HEIGHT`: Screen height (default: 1080)
- `SCREEN_DEPTH`: Color depth (default: 24)

### Persistence

The container persists:
- Chrome user data: `/home/chrome/.config/google-chrome`
- Downloads: `/home/chrome/Downloads`

These are mapped to Docker volumes for persistence across container restarts.

## Advanced Usage

### Custom Chrome Arguments

You can modify `start-chrome.sh` to add custom Chrome arguments as needed.

### VNC Password

By default, VNC runs without a password for localhost connections. To add a password, modify the `start-vnc.sh` script.

### Screen Resolution

Change screen resolution by setting environment variables:

```bash
docker run -d \
  -e SCREEN_WIDTH=1280 \
  -e SCREEN_HEIGHT=720 \
  chrome-novnc
```

## Troubleshooting

### Check Container Logs

```bash
docker-compose logs chrome-novnc
```

### Access Container Shell

```bash
docker exec -it chrome-novnc bash
```

### Chrome Process Issues

If Chrome doesn't start properly, check the supervisor logs:

```bash
docker exec -it chrome-novnc cat /var/log/supervisor/chrome.log
```

## Architecture

The container uses:
- **Ubuntu 22.04** as base
- **Xvfb** for virtual display
- **x11vnc** for VNC server
- **noVNC** for web interface
- **Fluxbox** as window manager
- **Supervisor** for process management
- **Google Chrome Stable** browser

## Security Notes

- The container runs Chrome with `--no-sandbox` for compatibility
- VNC is configured for localhost access only
- CDP is exposed on all interfaces for external access
- Consider using authentication for production environments
