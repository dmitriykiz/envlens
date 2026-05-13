# envlens

> Audit `.env` files across a monorepo for missing, duplicate, or sensitive keys.

---

## Installation

```bash
go install github.com/yourorg/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/envlens.git && cd envlens && go build -o envlens .
```

---

## Usage

Run `envlens` from the root of your monorepo:

```bash
envlens --dir ./services --baseline .env.example
```

### Example Output

```
[services/api/.env]
  ⚠  MISSING KEY   DATABASE_URL   (defined in .env.example)
  🔑 SENSITIVE KEY  AWS_SECRET_KEY (matches sensitive pattern)

[services/worker/.env]
  ✖  DUPLICATE KEY  LOG_LEVEL      (also found in services/api/.env)

Summary: 3 issues found across 5 .env files
```

### Flags

| Flag | Description |
|------|-------------|
| `--dir` | Root directory to scan (default: `.`) |
| `--baseline` | Reference `.env` file to check for missing keys |
| `--sensitive` | Path to custom sensitive key patterns file |
| `--json` | Output results as JSON |

---

## Why envlens?

Managing environment variables across multiple services is error-prone. `envlens` gives you a single command to catch configuration drift, accidental secret exposure, and inconsistencies before they reach production.

---

## License

MIT © [yourorg](https://github.com/yourorg)