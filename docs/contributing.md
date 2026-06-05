# Contributing

## Development Setup

1. Fork the repository
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/backend-ai-budget-splitter.git
cd backend-ai-budget-splitter
```

3. Add upstream remote:

```bash
git remote add upstream https://github.com/jackdafox/backend-ai-budget-splitter.git
```

4. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

## Making Changes

1. **Keep commits atomic** — one logical change per commit
2. **Write tests** for new functionality
3. **Run the linter** before pushing:

```bash
make lint
```

4. **Run tests**:

```bash
make test
```

5. **Build to verify**:

```bash
make build
```

## Code Style

- Follow Go idioms and `gofmt` formatting
- Document exported types and functions
- Use meaningful variable names
- Keep functions small and focused

## Commit Messages

```
type: short description

Longer explanation if needed.
Fixes #issue-number.
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

## Pull Request Process

1. Update documentation for any changed behavior
2. Ensure CI passes (lint + tests + build)
3. Request review from maintainers
4. Squash commits before merging

## Testing

```bash
# Run all tests with race detection
make test

# Run with coverage report
make test-coverage
```

## Building Docs

```bash
# Install MkDocs
pip install mkdocs-material

# Preview locally
mkdocs serve

# Build for deployment
mkdocs build
```

## Issues

Open issues for:
- Bug reports
- Feature requests
- Questions about the codebase
