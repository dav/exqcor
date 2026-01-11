# README

## Mobile API auth

API clients authenticate with bearer tokens returned by `/api/sign_in` or `/api/sign_up`.

- Send `Authorization: Bearer <token>` with every `/api/*` request.
- Tokens expire after 30 days.
- If the API response includes `X-Auth-Token` and `X-Auth-Token-Expires-At`, the token has rotated. Replace the stored token with the new value.

This README would normally document whatever steps are necessary to get the
application up and running.

Things you may want to cover:

* Ruby version

* System dependencies

* Configuration

* Database creation

* Database initialization

* How to run the test suite

* Services (job queues, cache servers, search engines, etc.)

* Deployment instructions

* ...
