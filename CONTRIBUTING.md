# Contributing to token-audit

Thanks for your interest in contributing!

## How to Contribute

### Reporting Bugs

- Use the [Bug Report](https://github.com/mrly2333/token-audit/issues/new?template=bug_report.md) template
- Include steps to reproduce, expected vs actual behavior
- Attach logs or screenshots if possible

### Suggesting Features

- Use the [Feature Request](https://github.com/mrly2333/token-audit/issues/new?template=feature_request.md) template
- Explain the use case and why it would be valuable

### Submitting Code

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Ensure the code compiles: `go build ./...`
5. Commit with a clear message
6. Push and open a Pull Request

## Development Setup

```bash
# Clone
git clone https://github.com/mrly2333/token-audit.git
cd token-audit

# Start PostgreSQL
docker compose up -d postgres

# Build
go build -o bin/newapi-audit-proxy ./cmd/newapi-audit-proxy

# Run
./bin/newapi-audit-proxy -config config.yaml
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add comments for non-obvious logic
- Update documentation if changing public behavior

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
