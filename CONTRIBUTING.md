# Contributing to Hata

Thank you for your interest in contributing! Here's how to get started.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone git@github.com:<your-username>/hata.git
   cd hata
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build and verify:
   ```bash
   go build ./...
   ```

## Development Workflow

1. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature
   ```
2. Make your changes
3. Build and test:
   ```bash
   go build ./...
   go vet ./...
   ```
4. Commit with a clear message (see [Commit Convention](#commit-convention))
5. Push and open a Pull Request

## Commit Convention

Use the format: `type: short description`

| Type | When to use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `refactor` | Code change with no feature/fix |
| `test` | Adding or updating tests |
| `chore` | Tooling, deps, CI |

Examples:
```
feat: add support for nested base.json
fix: handle duplicate keys in sheet
docs: update OAuth setup guide
```

## Project Structure

```
hata/
├── cmd/              # CLI commands (init, push, pull, diff)
├── internal/
│   ├── auth/         # OAuth and service account authentication
│   ├── config/       # Config file loading and saving
│   ├── diff/         # Key comparison logic
│   ├── i18n/         # base.json reading and locale file generation
│   ├── locale/       # Locale list and interactive selector
│   └── sheet/        # Google Sheets API client
├── example/          # Example project with base.json and config
├── main.go
└── go.mod
```

## Pull Request Guidelines

- Fill in the PR template
- Keep PRs focused — one feature or fix per PR
- Update `example/` or docs if your change affects usage
- Ensure `go build ./...` passes before submitting

## Reporting Issues

Please include:
- What you expected vs what happened
- Exact error message
- Your OS and Go version (`go version`)
- Your `i18n.config.yml` (with Sheet ID redacted)

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
