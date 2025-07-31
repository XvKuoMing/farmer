FROM ubuntu:22.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
RUN apt-get update && apt-get install -y \
    # Basic utilities
    wget \
    curl \
    unzip \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    # X11 and VNC dependencies
    xvfb \
    x11vnc \
    tigervnc-standalone-server \
    tigervnc-common \
    fluxbox \
    dbus-x11 \
    xfonts-base \
    xfonts-75dpi \
    xfonts-100dpi \
    # Audio support (optional)
    pulseaudio \
    # Python for noVNC
    python3 \
    python3-pip \
    python3-numpy \
    # Process management
    supervisor \
    # Networking tools
    net-tools \
    && rm -rf /var/lib/apt/lists/*

# Install Google Chrome
RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update \
    && apt-get install -y google-chrome-stable \
    && rm -rf /var/lib/apt/lists/*

# Install noVNC
RUN mkdir -p /opt/noVNC/utils/websockify \
    && wget -qO- https://github.com/novnc/noVNC/archive/v1.4.0.tar.gz | tar xz --strip 1 -C /opt/noVNC \
    && wget -qO- https://github.com/novnc/websockify/archive/v0.10.0.tar.gz | tar xz --strip 1 -C /opt/noVNC/utils/websockify \
    && chown -R root:root /opt/noVNC

# Create user for running Chrome
RUN useradd -m -s /bin/bash chrome \
    && usermod -aG audio chrome

# Create directories for persistence
RUN mkdir -p /home/chrome/.config/google-chrome \
    && mkdir -p /home/chrome/chrome-data \
    && mkdir -p /home/chrome/Downloads \
    && chown -R chrome:chrome /home/chrome

# Set up VNC directory
RUN mkdir -p /home/chrome/.vnc \
    && chown -R chrome:chrome /home/chrome/.vnc

# Environment variables
ENV DISPLAY=:1
ENV VNC_PORT=5901
ENV NOVNC_PORT=6080
ENV CDP_PORT=9222
ENV SCREEN_WIDTH=1920
ENV SCREEN_HEIGHT=1080
ENV SCREEN_DEPTH=24

# Copy configuration files
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY start-chrome.sh /usr/local/bin/start-chrome.sh
COPY start-vnc.sh /usr/local/bin/start-vnc.sh

# Make scripts executable
RUN chmod +x /usr/local/bin/start-chrome.sh \
    && chmod +x /usr/local/bin/start-vnc.sh

# Expose ports
EXPOSE $VNC_PORT $NOVNC_PORT $CDP_PORT

# Create volumes for persistence
VOLUME ["/home/chrome/.config/google-chrome", "/home/chrome/Downloads"]

# Start supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]