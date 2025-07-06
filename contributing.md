# Contributing to Ngingo

Welcome, warrior of logic. Whether you bring new features, refactor old runes, or raise flags of bug fixes — your hands are needed in the forge of Ngingo.

**Ngingo** is an AI-driven Reverse Proxy and Docker Container Orchestrator written in Go. Precision matters. Consistency matters. Contribution, however small, echoes across every container spun and every packet balanced.

---

## ⚙️ Project Setup

Before you write the first line of code, ensure your local environment is armed and ready.

### Requirements

- Go 1.20+ installed
- Docker (for container-based features)
- Basic knowledge of `git`, `go`, and CLI tooling

### Install dependencies

```bash
git clone https://github.com/keyzarrasya/ngingo.git
cd ngingo
go mod tidy
```

---

## 🧪 Testing Your Code

Run the test suite before committing your changes:

```bash
go test ./...
```

Want to lint your code?

```bash
go vet ./...
golint ./...
```

You are encouraged to install [`golangci-lint`](https://golangci-lint.run/) for comprehensive analysis.

---

## ✍️ Commit Message Guidelines

We use **[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)** to maintain clear, readable commit history and enable automated releases.

### 💬 Commit Format

```
<type>(optional scope): <short summary>

[optional body]

[optional footer]
```

### 🧱 Allowed Types

| Type        | Description                                                                 |
|-------------|-----------------------------------------------------------------------------|
| `feat`      | A new feature                                                               |
| `fix`       | A bug fix                                                                   |
| `docs`      | Documentation changes only                                                  |
| `style`     | Changes that do not affect code logic (formatting, spacing, etc.)          |
| `refactor`  | Code changes that neither fix a bug nor add a feature                      |
| `perf`      | Code changes that improve performance                                       |
| `test`      | Adding missing tests or refactoring existing tests                         |
| `chore`     | Routine tasks like dependency updates, config changes, etc.                |
| `ci`        | Changes to CI configuration files/scripts (e.g., GitHub Actions)           |
| `build`     | Changes that affect the build system or dependencies                       |
| `revert`    | Reverts a previous commit                                                   |

### ✅ Examples

```bash
feat(proxy): support multiple upstream rules per route
fix(container): fix nil pointer on resource stats
docs: update README with Docker orchestration guide
style: format handlers and remove unused imports
refactor(balancer): simplify load distribution logic
test: add unit tests for container health checker
```

---

## 🔁 Workflow for Contributors

1. **Fork** this repository
2. Create a branch:  
   `git checkout -b feat/your-new-feature`
3. Write clean, idiomatic Go code
4. Add and run tests
5. Commit your changes using the format above
6. Push to your fork
7. Open a Pull Request describing what you did and why

We'll review your PR, discuss if needed, and merge it if it aligns with the vision.

---

## 💡 Code Style & Best Practices

- Use [`gofmt`](https://golang.org/cmd/gofmt/), always.
- Avoid global state unless you have a really good reason.
- Keep functions small and expressive.
- Favor interfaces over concrete implementations for abstractions.
- Avoid premature optimization — clarity trumps cleverness.
- Document public methods and exported structs.

---

## 🐣 New to Open Source?

You're not just welcome — you're needed. Look for issues labeled:

- `good first issue`
- `help wanted`

Or open a draft PR if you're unsure — we’d rather guide you than lose your contribution.

---

> Let the architecture be elegant,  
> let the commits be poetic,  
> and may every PR  
> echo with purpose.

See you in the commit logs,  
— The Ngingo Maintainers 🛡️
