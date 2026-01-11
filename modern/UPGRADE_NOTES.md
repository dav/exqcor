# Upgrade Notes

- 2026-01-11: Generated Rails 8.0.4 app in `modern/` with Ruby 3.4.8; initial `bundle install`, `bin/rails test`, and `bin/rails server` boot completed.
- 2026-01-11: Added core schema/migrations/models for Scripts (renamed from Plays) plus Writers; SubSection now optionally belongs to Writer for anonymous mode; Character.actor remains optional.
- 2026-01-11: Added JSON-capable admin setup controllers/routes for Scripts, Sections, Characters, Props, and CharacterSections; `/api` scope mirrors the same routes with JSON defaults.
- 2026-01-11: Removed read-side VOSD mutation; added Script#ensure_vosd! and Section#ensure_vosd! with controller hooks; added unique indexes for Character name per Script and CharacterSection join pairs.
- 2026-01-11: Added `characters.role` with `vosd` role to decouple VOSD from name; added model validations and tests covering role-based VOSD behavior and idempotent VOSD joins.
- 2026-01-11: Pinned `minitest` to `~> 5.25` to avoid Rails test runner incompatibility with Minitest 6.
- 2026-01-11: Added API integration tests to verify script/section create flows ensure VOSD behavior.
- 2026-01-11: Added partial unique index to enforce one VOSD per script (Postgres compatible).
- 2026-01-11: Added HTML views/nav for admin setup and made controllers respond to HTML or JSON; added demo seed data for local browser use.
- 2026-01-11: Added writer-facing SubSection UI (HTML) and Lines controller for submitting lines; added sub_section routes and demo writers/seed line.
- 2026-01-11: Added Rails 8 authentication generator output plus API token auth (UserToken), web signup, and API sign-up/sign-in endpoints with bearer token support.
- 2026-01-11: Added admin user management UI and 30-day API token expiration with rotation (7-day window) plus bearer token headers for rotation.
- 2026-01-11: Documented mobile auth rotation headers in README and added a current-user badge to the admin nav.
- 2026-01-11: Added sign-in/sign-up links on the scripts index for quick access.
- 2026-01-11: Switched sign-out links to form buttons to avoid GET /session when Turbo method is unavailable.
