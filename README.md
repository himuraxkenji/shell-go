# shell-go

A POSIX-style shell written in Go. It reads commands in a REPL loop, parses
them with quote-aware tokenization (single and double quotes), runs builtin
commands and executes external programs found in `PATH`.

## Features

- Interactive REPL with `$ ` prompt
- Builtins: `exit`, `echo`, `type`
- Quote-aware argument parsing (`'single'` and `"double"` quotes)
- External command execution via `PATH` lookup

## Requirements

- Go 1.26+

## Usage

```sh
make run
```

## Development

```sh
make format   # gofmt all sources
make lint     # go vet + gofmt check
make test     # run tests
```
