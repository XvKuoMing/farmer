#!/bin/bash

echo "Building Chrome noVNC Docker image..."
docker build -t chrome-novnc .

echo "Build complete!"
echo ""
echo "To run the container:"
echo "docker-compose up -d"
echo ""
echo "Or run directly:"
echo "docker run -d --name chrome-novnc -p 6080:6080 -p 9222:9222 -v chrome-data:/home/chrome/.config/google-chrome chrome-novnc"
echo ""
echo "Access URLs:"
echo "- noVNC Web Interface: http://localhost:6080"
echo "- Chrome DevTools: http://localhost:9222"
echo "- VNC Direct: localhost:5901"