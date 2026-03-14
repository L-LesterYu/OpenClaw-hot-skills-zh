#!/bin/bash
# ByteRover - Quick status check

CONFIG=$(cat ~/.clawdbot/byterover-config.json)
API_KEY=$(echo "$CONFIG" | jq -r '.apiKey')
TEAM=$(echo "$CONFIG" | jq -r '.team')
SPACE=$(echo "$CONFIG" | jq -r '.space')

cd ~/.openclaw/workspace/skills/byterover
BRV_API_KEY="$API_KEY" BRV_TEAM="$TEAM" BRV_SPACE="$SPACE" \
  docker-compose exec -T byterover brv status --headless --format json
