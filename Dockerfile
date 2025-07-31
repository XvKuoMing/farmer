FROM python:3.12-slim-bookworm

# Install minimal system dependencies (no browser dependencies since Chrome runs separately)
RUN apt-get update && apt-get install -y \
    # Basic utilities
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy UV binary
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

# Copy the project into the image
ADD . /app

# Sync the project into a new environment, asserting the lockfile is up to date
WORKDIR /app
RUN uv sync --locked

# No need to install browsers - Chrome runs in separate container

# Set environment variables
ENV PYTHONUNBUFFERED=1
CMD ["uv", "run", "main.py"]