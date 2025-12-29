# Changelog

All notable changes to Valora OSS will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

### Security

## [0.1.0]

### Added

- Authorization module with policy model, service implementation, and tests.
- Audit logging domain/service plus admin audit log API endpoints.
- API key lifecycle endpoints (list/create/rotate/revoke) with password confirmation and rate limiting.
- OAuth/password flows and auth middleware additions for local and oauth2 providers.
- Migration, seed, scheduler, and cloud metrics components to support billing operations.
- Release docs/workflows (`RELEASE.md`, changelog, and GitHub Actions release automation).

### Changed

- Refactor billing cycle, invoice, rating, subscription, and usage services/models.
- Update server middleware, auth handling, subscription endpoints, and config wiring.
- Refresh UI routes/pages for login, invoices, subscriptions, API keys, and audit logs.
- Update Dockerfile and README to reflect release process.

### Fixed

### Security
