default:
    @just --list

# Cross-compile release binaries on Linux.
build tag: wasm
    GOOS=linux GOARCH=amd64 go build -buildvcs=false -o 'bin/darijascript-bin-{{tag}}-linux' .
    cp 'bin/darijascript-bin-{{tag}}-linux' bin/darijascript-bin
    GOOS=windows GOARCH=amd64 go build -buildvcs=false -o 'bin/darijascript-bin-{{tag}}-windows.exe' .
    GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o 'bin/darijascript-bin-{{tag}}-macos' .

# Build the playground interpreter and copy Go's matching bootstrap.
wasm:
    GOOS=js GOARCH=wasm go build -buildvcs=false -o website/public/darijascript-bin-web.wasm .
    root="$(go env GOROOT)"; source="$root/lib/wasm/wasm_exec.js"; \
        if [ ! -f "$source" ]; then source="$root/misc/wasm/wasm_exec.js"; fi; \
        cp "$source" website/public/wasm_exec.js

# Build the website with a freshly built interpreter.
build-web: wasm
    npm --prefix website run build
