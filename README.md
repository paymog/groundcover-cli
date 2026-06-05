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

```sh
export GROUNDCOVER_API_KEY=...
export GROUNDCOVER_BACKEND_ID=...
```

Also accepted: `GC_API_KEY`, `GC_BACKEND_ID`, `GROUNDCOVER_TENANT_UUID`, `GC_TENANT_UUID`, `GROUNDCOVER_BASE_URL`, `GC_BASE_URL`.

Defaults (override the env or flag for your own backend/tenant):

- Base URL: `https://api.groundcover.com`
- Backend ID: `groundcover` (`GROUNDCOVER_BACKEND_ID` / `--backend-id`)
- Tenant UUID: none. Set `GROUNDCOVER_TENANT_UUID` / `--tenant-uuid` to send `X-Tenant-UUID` on raw HAR-derived requests; the SDK transport already implies tenant from the API key.

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
```

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

## Release

Releases are tag-driven:

```sh
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions runs GoReleaser, then updates the source-build formula in `paymog/homebrew-tap`. The repo needs a `HOMEBREW_TAP_GITHUB_TOKEN` secret with write access to that tap.
