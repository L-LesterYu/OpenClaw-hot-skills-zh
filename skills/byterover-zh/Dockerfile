FROM node:20-bookworm-slim

# Install build dependencies for better-sqlite3
RUN apt-get update && apt-get install -y \
    python3 \
    build-essential \
    libsecret-1-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Install ByteRover CLI
RUN npm install -g byterover-cli

# Set up working directory
WORKDIR /workspace

# Set environment variables for headless mode
ENV BRV_HEADLESS=true
ENV BRV_FORMAT=json

# Default command - keep container alive
CMD ["tail", "-f", "/dev/null"]
