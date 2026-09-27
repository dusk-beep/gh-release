# gh-release

A concurrent release tracker TUI. The application is provider-oriented so GitHub is an implementation, not the application model.

## Requirements

- Go 1.26+

## Run

```sh
cp config.example.toml ~/.config/gh-release/config.toml
go run ./cmd/gh-release
```

Optional GitHub authentication can be configured with `github.token`. Authentication is useful for avoiding the low unauthenticated API rate limit.

Keys:

- `r` refresh
- `q` / `Ctrl-C` quit
- arrow keys / `j` / `k` navigate the table
