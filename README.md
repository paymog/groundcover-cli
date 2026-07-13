# groundcover-cli

Go CLI for Groundcover APIs.

It uses the official `github.com/groundcover-com/groundcover-sdk-go` where that SDK has a stable contract, and keeps HAR-derived raw commands for webapp endpoints the SDK does not expose.

## Install

```sh
brew tap paymog/tap
brew install groundcover
```

The formula builds from source, so Homebrew installs a Go toolchain as a build dependency.

Go install:

```sh
go install github.com/paymog/groundcover-cli/cmd/groundcover@latest
```

Local build:

```sh
go build ./cmd/groundcover
```

## Auth

The simplest setup is environment variables (recommended for CI):

```sh
export GROUNDCOVER_API_KEY=...
export GROUNDCOVER_BACKEND_ID=...
```

Also accepted: `GC_API_KEY`, `GC_BACKEND_ID`, `GROUNDCOVER_TENANT_UUID`, `GC_TENANT_UUID`, `GROUNDCOVER_BASE_URL`, `GC_BASE_URL`, `GROUNDCOVER_GRAFANA_SERVICE_ACCOUNT_TOKEN`, `GC_GRAFANA_SERVICE_ACCOUNT_TOKEN`.

### Stored profiles

For day-to-day use you can store credentials in named profiles. The API key is
kept in your OS keyring (macOS Keychain, Linux libsecret/secret-service, Windows
Credential Manager); only non-secret metadata (backend ID, base URL, tenant UUID)
is written to `~/.config/groundcover/profiles.yaml` (XDG-aware; `%AppData%` on
Windows).

```sh
groundcover auth login              # prompts for the API key; profile "default"
groundcover auth login prod --backend-id my-backend --key lin_...
groundcover auth list               # table; * marks the default
groundcover auth default prod       # set the default profile
groundcover --profile prod monitors list   # use a profile for one command
groundcover auth status             # show which source resolved
groundcover auth token              # print the resolved API key
groundcover auth logout prod        # remove a profile (-f to skip confirm)
```

`login` validates the key against the API before saving it. The first profile
added becomes the default.

**Credential precedence** (highest first):

1. `--api-key` flag / `GROUNDCOVER_API_KEY` env var
2. `--profile <name>` → stored profile
3. default profile → stored profile

Combining an explicit key with `--profile` is rejected as ambiguous. If no keyring
is available, fall back to `GROUNDCOVER_API_KEY`.

Defaults (override the env or flag for your own backend/tenant):

- Base URL: `https://api.groundcover.com`
- Backend ID: `groundcover` (`GROUNDCOVER_BACKEND_ID` / `--backend-id`)
- Tenant UUID: none. Set `GROUNDCOVER_TENANT_UUID` / `--tenant-uuid` to send `X-Tenant-UUID` on raw HAR-derived requests; the SDK transport already implies tenant from the API key.
- Grafana service account token: none. The embedded Grafana (`raw grafana …`) endpoints are session-gated and ignore the API key, so they need a Grafana service account token (`glsa_…`). Set `GROUNDCOVER_GRAFANA_SERVICE_ACCOUNT_TOKEN` / `--grafana-token`. It is only used for `raw grafana …` commands. Running one without a token prints setup steps; generate a token with groundcover's official CLI (`~/.groundcover/bin/groundcover auth generate-service-account-token`, needs tenant admin).

## SDK-backed commands

```sh
groundcover dashboards list
groundcover dashboards get <id>
groundcover dashboards create --body-file dashboard.json
groundcover dashboards update <id> --body-file dashboard.json
groundcover dashboards delete <id>

groundcover monitors list --query 'monitor_name = "cpu"'
groundcover monitors get <id>
groundcover monitors create --body-file monitor.yaml
groundcover monitors update <id> --body-file monitor.yaml
groundcover monitors delete <id>

groundcover silences list --active
groundcover silences create --body-file silence.json
groundcover silences delete <id>

groundcover dashboards archive <id>
groundcover dashboards restore <id>

groundcover recurring-silences list
groundcover recurring-silences create --body-file silence.json

groundcover connected-apps list --query 'type:slack-webhook'
groundcover connected-apps get <id>
groundcover connected-apps create --body-file app.json
groundcover connected-apps update <id> --body-file app.json
groundcover connected-apps delete <id>

groundcover notification-routes list --query 'prod'
groundcover notification-routes get <id>
groundcover notification-routes create --body-file route.json
groundcover notification-routes update <id> --body-file route.json
groundcover notification-routes delete <id>

# Auth / RBAC
groundcover api-keys list
groundcover api-keys create --body-file key.json
groundcover service-accounts list
groundcover service-accounts create --body-file sa.json
groundcover ingestion-keys list
groundcover policies list
groundcover policies apply --body-file policy.json
groundcover policies audit-trail <id>

# Synthetics, secrets, workflows
groundcover synthetics list
groundcover synthetics create --body-file test.json
groundcover secrets create --body-file secret.json
groundcover secrets hash <id>
groundcover workflows list
groundcover workflows create --body-file workflow.yaml

# Pipeline configs (singleton get/create/update/delete)
groundcover logs-pipeline get
groundcover metrics-pipeline get
groundcover traces-pipeline get
groundcover metrics-aggregator get

# Integrations (typed)
groundcover integrations list
groundcover integrations describe <type>
groundcover integrations create <type> --body-file config.json
groundcover integrations update <type> <id> --body-file config.json

# Read/query
groundcover logs search --body-file search.json
groundcover traces search --body-file search.json
groundcover metrics query --body-file query.json
groundcover metrics names --body-file body.json
groundcover search discovery --body-file body.json
groundcover k8s clusters --body-file body.json
groundcover k8s workloads --body-file body.json
groundcover k8s events-search --body-file body.json
```

## Raw HAR-derived commands

```sh
groundcover raw list
groundcover raw list k8s
groundcover raw k8s clusters list
groundcover raw dashboards get --dashboard-id <id>
groundcover raw metrics query-range --body-file body.json
groundcover raw prometheus api query --query query='up'
groundcover raw grafana dashboards get --dashboard-uid <uid>
groundcover raw grafana dashboards save --body-file dashboard.json
groundcover raw grafana folders list
groundcover raw grafana ds query --body-file query.json
```

Storage management uses one endpoint per data type:

```sh
groundcover raw storage-management get --data-type logs --raw \
  | jq '{retention,version,cold_move_duration,cold_volume,custom_rules:(.custom_rules // [])}' \
  > storage.json
# Edit storage.json, preserving every writable field and the complete rule list.
groundcover raw storage-management update --data-type logs --body-file storage.json
```

Supported data types are `logs`, `traces`, `events`, `measurements`, and `monitor_instance`. Updates use optimistic concurrency and replace the writable settings document: start from `get`, pass its current `version`, and preserve `retention`, `cold_move_duration`, `cold_volume`, and the full `custom_rules` list. Omitting `custom_rules` removes existing rules.


Raw commands support:

- `--body-json '<json>'`
- `--body-file path.json`
- `--body-file path.yaml`
- `--set dotted.path=value`
- `--query key=value`
- path flags generated from captured UUIDs, for example `--dashboard-id <id>`
- `--raw`

## Regenerate raw commands

```sh
go run ./scripts/generate-commands.go ~/Downloads/app.groundcover.com.har
```

Generated commands are written to `internal/raw/commands_generated.go`.

## Shape

- `internal/sdkcmd`: first-class commands backed by the official SDK.
- `internal/raw`: best-effort command registry and runner for HAR-derived endpoints.
- `internal/config`: shared auth, base URL, timeout, and SDK transport setup.

## Claude Code skill

This repo ships a [Claude Code](https://claude.com/claude-code) skill that teaches the
agent how to drive the CLI (auth, the SDK-vs-raw split, and ready-made request-body
templates for logs/traces/metrics/k8s). It lives in [`skills/groundcover-cli`](skills/groundcover-cli).

Install it with [`npx skills`](https://github.com/vercel-labs/skills) (Vercel's agent-skills tool):

```sh
# install into the current project (.claude/skills/)
npx skills add paymog/groundcover-cli

# or install globally for your user, skipping prompts
npx skills add paymog/groundcover-cli --global --yes
```

Useful flags: `--list` to preview without installing, `--skill groundcover-cli` to target it
explicitly, `-a claude-code` to pick the agent.

Or install manually:

```sh
git clone https://github.com/paymog/groundcover-cli
cp -r groundcover-cli/skills/groundcover-cli ~/.claude/skills/groundcover-cli
```

Then ask Claude Code to query logs, manage monitors, debug a prod issue, etc., and it will
invoke `groundcover`.

## Release

Releases are tag-driven:

```sh
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions runs GoReleaser, then updates the source-build formula in `paymog/homebrew-tap`. The workflow pushes to the tap over SSH using a deploy key stored as the `HOMEBREW_TAP_DEPLOY_KEY` secret (the public half is a read-write deploy key on the tap repo).
