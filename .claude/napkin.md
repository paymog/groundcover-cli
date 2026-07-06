# Napkin

## Corrections
| Date | Source | What Went Wrong | What To Do Instead |
|------|--------|----------------|-------------------|
| 2026-07-06 | self | Ran the HAR generator directly on a narrow dashboard HAR and it overwrote the broad raw command registry with only 24 commands | For additive HAR integrations, preserve the generated baseline and add/merge new commands instead of replacing unrelated captured endpoints |

## User Preferences
- Always commit any changes made to this napkin file.

## Patterns That Work
- (approaches that succeeded)

## Patterns That Don't Work
- (approaches that failed and why)

## Domain Notes
- (project/domain context that matters)
