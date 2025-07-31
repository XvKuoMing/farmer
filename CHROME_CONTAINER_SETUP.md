# Chrome Container Setup

This setup runs Chrome browser in a separate container from your main application, allowing you to connect via Chrome DevTools Protocol (CDP).

## Architecture

- **chrome**: Separate container running Chrome with remote debugging enabled
- **aisearch**: Main application container that connects to Chrome via CDP
- Both containers communicate over a Docker network

## Features

- ✅ Chrome runs in head mode (visible browser)
- ✅ Chrome runs as a separate container
- ✅ VNC access to view the browser remotely
- ✅ Remote debugging enabled on port 9222
- ✅ Automatic connection retry logic
- ✅ Proper container orchestration with docker-compose

## Usage

### Start the services
```bash
docker-compose up --build
```

### View the browser remotely (VNC)
Connect to `localhost:5900` with any VNC client to see the Chrome browser running.

### Access Chrome DevTools
Open `http://localhost:9222` in your local browser to access Chrome DevTools.

### Environment Variables

- `CHROME_WS_ENDPOINT`: WebSocket endpoint for Chrome (default: `ws://chrome:9222`)
- `ENABLE_VNC`: Enable VNC server (default: `true`)
- `OPENAI_API_KEY`: Your OpenAI API key
- `OPENAI_BASE_URL`: Optional custom OpenAI base URL

## Container Details

### Chrome Container (`chrome`)
- Based on Debian with Google Chrome stable
- Runs Xvfb for virtual display
- Exposes port 9222 for remote debugging
- Exposes port 5900 for VNC access
- Runs as non-root user for security

### Application Container (`aisearch`)
- Minimal Python container
- No browser dependencies (connects to remote Chrome)
- Waits for Chrome container to be ready before starting

## Troubleshooting

### Chrome not connecting
- Check if Chrome container is running: `docker-compose ps`
- View Chrome container logs: `docker-compose logs chrome`
- Verify network connectivity: `docker-compose exec aisearch ping chrome`

### VNC not working
- Ensure port 5900 is not blocked by firewall
- Try different VNC clients (TigerVNC, RealVNC, etc.)
- Check VNC is enabled: `ENABLE_VNC=true`

### Application errors
- Check application logs: `docker-compose logs aisearch`
- Verify environment variables are set correctly
- Ensure Chrome container started successfully

## Benefits of This Setup

1. **Isolation**: Browser crashes don't affect the main application
2. **Scalability**: Multiple app instances can connect to the same Chrome
3. **Development**: Easy to view browser activity via VNC
4. **Debugging**: Direct access to Chrome DevTools
5. **Resource Management**: Separate resource allocation for browser vs app