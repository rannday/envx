# envx
Go environment loading module.

## Behavior
- `Load` reads configuration into a struct from environment variables (`os.LookupEnv`).
- Optional `.env` loading is available via `Options.DotEnvPath` as a fallback source.
- Precedence is:
1. Real process environment
2. Optional `.env` file values
3. Struct `default` tag
- `required:"true"` requires a non-empty value.
- Use `allowempty:"true"` with `required:"true"` to permit explicit empty string values.

This is designed to work well with production systems like `systemd EnvironmentFile=...`, where the service manager provides the process environment and the app reads env directly.
