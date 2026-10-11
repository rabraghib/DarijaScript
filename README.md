<img align="right" src="language-service/images/icon.png" alt="DarijaScript Logo" width="150" height="150">

# DarijaScript

DarijaScript is a programming language inspired by the Moroccan Arabic language.

<br clear="right"/>

## Note

DarijaScript is not intended for serious projects or production use.
It's a lighthearted experiment in language design and cultural expression.
Have fun exploring and playing around with DarijaScript!

## Getting Started

You can get started with DarijaScript by using the [DarijaScript Starter](https://github.com/rabraghib/darijascript-starter) template.

### Installation

To install DarijaScript, you need to have [Go](https://go.dev/) installed.
After installing Go, you can use the following command to install DarijaScript:

```bash
go install github.com/rabraghib/darijascript@latest
```

Make sure your Go bin directory is in your system's PATH.

### Usage

You can either write your DarijaScript code in a file and run it:

```bash
darijascript run my_program.ds
```

Or you can run DarijaScript in interactive mode:

```bash
darijascript interactive
```

### Examples

Discover various DarijaScript code examples in the [examples](examples) folder.

## Documentation

Try programs in the [DarijaScript Online Playground](https://ds.rabraghib.com/).

### Development

Run `go build .` to build the interpreter and CLI.

The `justfile` targets Linux (or WSL). For the website, install Node.js 22 or newer, Go, and [just](https://just.systems), then run from the repository root:

```sh
npm --prefix website ci
just dev
```

The `just dev` and `just build-web` commands build the WASM binary and copy the unmodified `wasm_exec.js` from the same Go compiler. After changing Go code while the development server is running, run `just wasm` and reload the playground.

Run `just build-web` to create a production website build. Deployment uses `npm run build` inside `website` with the WASM artifacts supplied by CI.

CI runs `just build <tag>` to build the Windows, Linux, and macOS release binaries along with the playground WASM artifacts.

## 🤗 Before you go!

Contributions are always welcome!

If you want to show some support, you can:

- Give [DarijaScript repo](https://github.com/rabraghib/darijascript) a star ⭐ on Github.
- Follow [@rabraghib](https://twitter.com/rabraghib) on Twitter to stay up to date with all the cool stuff happening.
- 💰 Help DarijaScript grow and be more awesome: [Buy me a coffee](https://www.buymeacoffee.com/rabraghib).

---

<p align="center">Made with ❤️ by <a href="https://rabraghib.com/">Rabyâ Raghib</a></p>
