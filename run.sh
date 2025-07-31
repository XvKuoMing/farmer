#!/bin/bash

# Chrome Container Setup - Quick Start Script

echo "🚀 Starting Chrome Container Setup..."
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating template..."
    cat > .env << EOF
# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_BASE_URL=

# VNC Configuration
ENABLE_VNC=true
EOF
    echo "📝 Please edit .env file with your OpenAI API key before running again."
    exit 1
fi

# Build and start containers
echo "🔨 Building and starting containers..."
docker-compose up --build -d

echo ""
echo "✅ Setup complete!"
echo ""
echo "📊 Container Status:"
docker-compose ps

echo ""
echo "🌐 Access Points:"
echo "  • Chrome DevTools: http://localhost:9222"
echo "  • VNC Viewer: localhost:5900 (to see Chrome browser)"
echo ""
echo "📋 Useful Commands:"
echo "  • View logs: docker-compose logs -f"
echo "  • Stop services: docker-compose down"
echo "  • Restart: docker-compose restart"
echo ""
echo "🔍 Troubleshooting:"
echo "  • Check Chrome logs: docker-compose logs chrome"
echo "  • Check app logs: docker-compose logs aisearch"
echo ""