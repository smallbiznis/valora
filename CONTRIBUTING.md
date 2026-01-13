# Contributing to Valora OSS

Thank you for your interest in contributing to Valora!
To maintain a high-quality codebase and a stable release pipeline, we follow a strict **Feature Branch Workflow**.

## 🚀 Branching Strategy

We do **NOT** commit directly to `main`. All changes must go through a Pull Request.

### Branch Types

Please use the following prefixes for your branches:

- **`feat/`**: New features (e.g., `feat/auth-provider-google`)
- **`fix/`**: Bug fixes (e.g., `fix/docker-build-error`)
- **`chore/`**: Maintenance, config tasks, or dependency updates (e.g., `chore/upgrade-go-1.25`)
- **`docs/`**: Documentation updates (e.g., `docs/update-readme`)
- **`refactor/`**: Code restructuring without behavior change (e.g., `refactor/centralize-config`)

### The Golden Rule
> **The `main` branch is always clean, buildable, and deployable.**

## 🛠️ Development Workflow

1.  **Sync**: Always start from up-to-date main.
    ```bash
    git checkout main
    git pull origin main
    ```
2.  **Branch**: Create your feature branch.
    ```bash
    git checkout -b feat/my-awesome-feature
    ```
3.  **Work**: Write code, run tests, and commit locally.
    - We encourage [Conventional Commits](https://www.conventionalcommits.org/):
      - `feat: add google oauth support`
      - `fix: resolve docker build failure`
4.  **Push**: Push your branch to origin.
    ```bash
    git push -u origin feat/my-awesome-feature
    ```
5.  **Pull Request (PR)**:
    - Open a PR on GitHub targetting `main`.
    - Fill in the description.
    - Request review.
6.  **Merge**:
    - Once approved and CI passes, use **Squash and Merge** (preferred) to keep history linear.
    - Delete the branch after merging.

## 📦 Release Process

Releases are **Tag-Driven**.

1.  Ensure all features for the proper version are merged into `main`.
2.  Create and push a tag:
    ```bash
    git tag -a v1.0.0 -m "Release v1.0.0"
    git push origin v1.0.0
    ```
3.  GitHub Actions will automatically build Docker images and create a GitHub Release.
