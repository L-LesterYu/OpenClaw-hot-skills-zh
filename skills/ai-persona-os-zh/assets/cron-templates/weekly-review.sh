# Weekly Review — Cron Job Template
# ⚠️  OPT-IN ONLY: This template is NOT auto-installed.
# The user must explicitly request cron setup ("set up cron jobs")
# and manually copy-paste this command into their terminal.
# This skill NEVER runs this command automatically.
#
# Requires: openclaw CLI (pre-installed with OpenClaw)
# Effect: Creates a scheduled job that runs every Monday at 9 AM
# Scope: Runs in an isolated session — reads/writes workspace files only
# Network: No network activity — reads local files only
#
# Deep review of the past week: learnings, archiving, pattern recognition
# Schedule: Monday at 9 AM (adjust timezone)
# Uses Opus model for deeper analysis
#
# Usage:
#   Copy and paste the command below into your terminal.
#   Change --tz to your timezone.
#   Remove --model opus if you prefer your default model.

openclaw cron add \
  --name "ai-persona-weekly-review" \
  --cron "0 9 * * 1" \
  --tz "America/Los_Angeles" \
  --session isolated \
  --model opus \
  --message "每周回顾 protocol:

1. Scan memory/ for the past 7 days. Summarize key themes, decisions, and outcomes.

2. Review .learnings/LEARNINGS.md — promote items with 3+ repetitions to MEMORY.md or AGENTS.md.

3. Archive logs older than 90 days to memory/archive/.

4. Check MEMORY.md size — prune if approaching 4KB.

5. Review WORKFLOWS.md — any new recurring patterns worth documenting?

6. Deliver a weekly summary: wins, issues resolved, lessons learned, and focus areas for the coming week.

Use 🟢🟡🔴 indicators for overall system health." \
  --announce
