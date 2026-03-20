/**
 * Agent Runner - MoltGuard API-based analysis
 *
 * Sends content to the MoltGuard API for prompt injection detection
 * instead of calling an LLM directly.
 */

import type {
  AnalysisTarget,
  AnalysisVerdict,
  Finding,
  Logger,
  MoltGuardApiResponse,
} from "./types.js";
import {
  DEFAULT_API_BASE_URL,
  loadApiKey,
  registerApiKey,
} from "./config.js";
import { sanitizeContent } from "./sanitizer.js";

// =============================================================================
// Runner Config
// =============================================================================

export type RunnerConfig = {
  apiKey: string;
  timeoutMs: number;
  autoRegister: boolean;
  apiBaseUrl: string;
};

// =============================================================================
// API Key Resolution
// =============================================================================

async function ensureApiKey(
  configKey: string,
  autoRegister: boolean,
  apiBaseUrl: string,
  log: Logger,
): Promise<string> {
  // 1. Use explicitly configured key
  if (configKey) {
    return configKey;
  }

  // 2. Try loading from credentials file
  const savedKey = loadApiKey();
  if (savedKey) {
    return savedKey;
  }

  // 3. Auto-register (if enabled)
  if (!autoRegister) {
    throw new Error(
      "No API key configured and autoRegister is disabled. " +
      "Please set apiKey in your MoltGuard plugin config or enable autoRegister.",
    );
  }

  log.info("No API key found — registering with MoltGuard...");

  try {
    // Add 10-second timeout for registration to prevent hanging
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 10000);

    const newKey = await registerApiKey("openclaw-agent", apiBaseUrl);
    clearTimeout(timeoutId);

    log.info("Registered with MoltGuard. API key saved to ~/.openclaw/credentials/moltguard/credentials.json");
    return newKey;
  } catch (error) {
    throw new Error(
      `Failed to auto-register API key: ${error instanceof Error ? error.message : String(error)}. ` +
      "Please check your network connection or set apiKey manually in config."
    );
  }
}

// =============================================================================
// API Response Mapping
// =============================================================================

export function mapApiResponseToVerdict(apiResponse: MoltGuardApiResponse): AnalysisVerdict {
  const verdict = apiResponse.verdict;

  const findings: Finding[] = (verdict.findings ?? []).map((f) => ({
    suspiciousContent: f.suspiciousContent,
    reason: f.reason,
    confidence: f.confidence,
  }));

  return {
    isInjection: verdict.isInjection,
    confidence: verdict.confidence,
    reason: verdict.reason,
    findings,
    chunksAnalyzed: 1, // API handles chunking server-side
  };
}

// =============================================================================
// Main Analysis Function
// =============================================================================

export async function runGuardAgent(
  target: AnalysisTarget,
  config: RunnerConfig,
  log: Logger,
): Promise<AnalysisVerdict> {
  const startTime = Date.now();

  log.info(`Analyzing content: ${target.content.length} chars`);

  // Sanitize content locally before sending to API
  const { sanitized, redactions, totalRedactions } = sanitizeContent(target.content);
  if (totalRedactions > 0) {
    log.info(`Sanitized ${totalRedactions} sensitive items: ${Object.entries(redactions).map(([k, v]) => `${v} ${k}`).join(", ")}`);
  }

  // Ensure we have an API key
  const baseUrl = config.apiBaseUrl || DEFAULT_API_BASE_URL;
  const apiKey = await ensureApiKey(config.apiKey, config.autoRegister, baseUrl, log);

  // Call MoltGuard API
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), config.timeoutMs);

  try {
    const response = await fetch(
      `${baseUrl}/api/check/tool-call`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${apiKey}`,
        },
        body: JSON.stringify({
          content: sanitized,
          async: false,
        }),
        signal: controller.signal,
      },
    );

    if (!response.ok) {
      throw new Error(`MoltGuard API error: ${response.status} ${response.statusText}`);
    }

    const apiResponse = (await response.json()) as MoltGuardApiResponse;

    if (!apiResponse.ok) {
      throw new Error(`MoltGuard API returned error: ${apiResponse.error ?? "unknown"}`);
    }

    const verdict = mapApiResponseToVerdict(apiResponse);

    const durationMs = Date.now() - startTime;
    log.info(`Analysis complete in ${durationMs}ms: ${verdict.isInjection ? "INJECTION DETECTED" : "SAFE"}`);

    return verdict;
  } catch (error) {
    if ((error as Error).name === "AbortError") {
      log.warn("Analysis timed out");
      return {
        isInjection: false,
        confidence: 0,
        reason: "Timeout",
        findings: [],
        chunksAnalyzed: 0,
      };
    }
    throw error;
  } finally {
    clearTimeout(timeoutId);
  }
}
