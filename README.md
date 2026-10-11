<img align="right" src="vscode/images/icon.png" alt="DarijaScript Logo" width="150" height="150">

# DarijaScript

DarijaScript is an experimental programming language inspired by Moroccan Arabic, built for fun and learning.

<br clear="right"/>

## Getting Started

Install [Go](https://go.dev/), then install the CLI:

```bash
go install github.com/rabraghib/darijascript@latest
```

Add your Go bin directory to `PATH`. Run a `.ds` file or start an interactive session:

```bash
darijascript run my_program.ds
darijascript interactive
```

Check syntax and format files:

```bash
darijascript check my_program.ds
darijascript format --write my_program.ds
```

- [Online playground](https://ds.rabraghib.com/)
- [Language reference](vscode/LANGUAGE.md)
- [Examples](examples)

The [VS Code extension](vscode) provides highlighting, completions, diagnostics, and formatting. It uses the CLI on `PATH`, or the executable configured in `darijascript.interpreterPath`.

Open the guide from the Command Palette with **DarijaScript: Open Language Reference**.

In the VS Code Markdown editor and preview, use `darijascript` or `ds` on fenced code blocks for syntax highlighting:

````markdown
```darijascript
9ayed name = "Salam";
golih(name);
```
````

## Development

Run `go build .` to build the interpreter and CLI.

For the website, use Linux or WSL with Node.js 22+, Go, and [just](https://just.systems). From the repository root:

```sh
npm --prefix website ci
just wasm
npm --prefix website run dev
```

After changing Go code, rerun `just wasm` and reload the playground. After editing examples, run `npm --prefix website run sync-examples`.

Verify changes:

```sh
go test ./...
npm --prefix website run typecheck
npm --prefix website run lint
```

Use `just build-web` for a production website build, or `just build <tag>` for release binaries and playground WASM artifacts.

## 🤗 Before you go!

Contributions are always welcome!

- [Give the repo a star ⭐](https://github.com/rabraghib/darijascript).
- Follow [@rabraghib](https://x.com/rabraghib) on X for updates.

---

<p align="center">Made with ❤️ by <a href="https://rabraghib.com/">Rabyâ RAGHIB</a></p>
