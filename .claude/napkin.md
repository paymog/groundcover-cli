# Napkin

## Corrections
| Date | Source | What Went Wrong | What To Do Instead |
|------|--------|----------------|-------------------|
| 2026-07-06 | self | Ran the HAR generator directly on a narrow dashboard HAR and it overwrote the broad raw command registry with only 24 commands | For additive HAR integrations, preserve the generated baseline and add/merge new commands instead of replacing unrelated captured endpoints |

## User Preferences
- Always commit any changes made to this napkin file.

## Patterns That Work
- `raw grafana …` embedded-Grafana commands authenticate with a Grafana service account token (`glsa_…`), NOT the `gcsa_` API key. Wired in via `GROUNDCOVER_GRAFANA_SERVICE_ACCOUNT_TOKEN` / `--grafana-token`; the runner uses a plain http client for WebApp commands (the SDK transport would clobber Authorization with the gcsa bearer). Verified live: `raw grafana search` / `folders list` return real JSON. Token lives in ~/code/goldskyio/.envrc.

## Patterns That Don't Work
- (approaches that failed and why)

## Domain Notes
- `raw grafana …` commands CANNOT authenticate with the gcsa bearer token (RESOLVED: use a `glsa_` grafana service account token, see Patterns That Work). Embedded Grafana at app.groundcover.com/grafana/* is session-gated. A gcsa/bearer request to `/grafana/api/*` hits a catch-all returning the ~980KB Grafana SPA `index.html` (200 text/html), never JSON. Signal of a real Grafana response = `grafana-trace-id` response header + content-type application/json.
- The gcsa bearer IS valid against the real GC API (api.groundcover.com/api/*) — verified via GET /api/monitors/recurring-silences and GET /api/dashboards (both return JSON). So the grafana breakage is the proxy path, not the key.
- For dashboards over the CLI, use the GC-native SDK `dashboards` resource (/api/dashboards), NOT the embedded-Grafana `raw grafana …` passthrough.
