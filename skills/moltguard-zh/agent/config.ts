/**
 * MoltGuard configuration and API key management
 */

import type { OpenClawGuardConfig } from "./types.js";
import path from "node:path";
import os from "node:os";
import fs from "node:fs";

// =============================================================================
// API Configuration
// =============================================================================

export const DEFAULT_API_BASE_URL = "https://api.moltguard.com";

const CREDENTIALS_DIR = path.join(os.homedir(), ".openclaw/credentials/moltguard");
const CREDENTIALS_FILE = path.join(CREDENTIALS_DIR, "credentials.json");

// =============================================================================
// API Key Management
// =============================================================================

export function loadApiKey(): string | null {
  try {
    if (!fs.existsSync(CREDENTIALS_FILE)) {
      return null;
    }
    const data = JSON.parse(fs.readFileSync(CREDENTIALS_FILE, "utf-8"));
    return typeof data.apiKey === "string" ? data.apiKey : null;
  } catch {
    return null;
  }
}

export function saveApiKey(apiKey: string): void {
  if (!fs.existsSync(CREDENTIALS_DIR)) {
    fs.mkdirSync(CREDENTIALS_DIR, { recursive: true });
  }
  fs.writeFileSync(
    CREDENTIALS_FILE,
    JSON.stringify({ apiKey }, null, 2),
    "utf-8",
  );
}

export async function registerApiKey(agentName: string, baseUrl: string = DEFAULT_API_BASE_URL): Promise<string> {
  const response = await fetch(`${baseUrl}/api/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ agentName }),
  });

  if (!response.ok) {
    throw new Error(`Registration failed: ${response.status} ${response.statusText}`);
  }

  const data = (await response.json()) as { apiKey: string };
  if (!data.apiKey) {
    throw new Error("Registration response missing apiKey");
  }

  saveApiKey(data.apiKey);
  return data.apiKey;
}

// =============================================================================
// Default Configuration
// =============================================================================

export const DEFAULT_CONFIG: Required<OpenClawGuardConfig> = {
  enabled: true,
  sanitizePrompt: false,
  gatewayPort: 8900,
  gatewayAutoStart: true,
  blockOnRisk: true,
  apiKey: "",
  timeoutMs: 60000,
  logPath: path.join(os.homedir(), ".openclaw", "logs"),
  autoRegister: true,
  apiBaseUrl: DEFAULT_API_BASE_URL,
};

// =============================================================================
// Configuration Helpers
// =============================================================================

export function resolveConfig(config?: Partial<OpenClawGuardConfig>): Required<OpenClawGuardConfig> {
  return {
    enabled: config?.enabled ?? DEFAULT_CONFIG.enabled,
    sanitizePrompt: config?.sanitizePrompt ?? DEFAULT_CONFIG.sanitizePrompt,
    gatewayPort: config?.gatewayPort ?? DEFAULT_CONFIG.gatewayPort,
    gatewayAutoStart: config?.gatewayAutoStart ?? DEFAULT_CONFIG.gatewayAutoStart,
    blockOnRisk: config?.blockOnRisk ?? DEFAULT_CONFIG.blockOnRisk,
    apiKey: config?.apiKey ?? DEFAULT_CONFIG.apiKey,
    timeoutMs: config?.timeoutMs ?? DEFAULT_CONFIG.timeoutMs,
    logPath: config?.logPath ?? DEFAULT_CONFIG.logPath,
    autoRegister: config?.autoRegister ?? DEFAULT_CONFIG.autoRegister,
    apiBaseUrl: config?.apiBaseUrl ?? DEFAULT_CONFIG.apiBaseUrl,
  };
}
