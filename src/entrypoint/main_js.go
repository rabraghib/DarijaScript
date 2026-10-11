//go:build js

package entrypoint

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall/js"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/runtime"
)

// Waiting on a Go channel yields to JavaScript so the worker can receive input.
type browserInput struct {
	lines   chan string
	pending []byte
}

func (input *browserInput) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, nil
	}
	if len(input.pending) == 0 {
		js.Global().Call("requestDarijaScriptInput")
		input.pending = []byte(<-input.lines + "\n")
	}
	n := copy(buffer, input.pending)
	input.pending = input.pending[n:]
	return n, nil
}

func Entrypoint() {
	js.Global().Set("checkDarijaScript", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 || args[0].Type() != js.TypeString {
			return "[]"
		}
		encoded, _ := json.Marshal(runtime.Check(args[0].String()))
		return string(encoded)
	}))
	js.Global().Set("formatDarijaScript", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 || args[0].Type() != js.TypeString {
			return map[string]interface{}{"ok": false, "error": "expected source text"}
		}
		formatted, err := runtime.Format(args[0].String())
		if err != nil {
			return map[string]interface{}{"ok": false, "error": err.Error()}
		}
		return map[string]interface{}{"ok": true, "source": formatted}
	}))
	var activeInput *browserInput
	js.Global().Set("submitDarijaScriptInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if activeInput != nil && len(args) == 1 && args[0].Type() == js.TypeString {
			select {
			case activeInput.lines <- args[0].String():
			default:
			}
		}
		return nil
	}))
	js.Global().Set("runDarijaScript", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 || args[0].Type() != js.TypeString {
			return js.Global().Get("Promise").Call("resolve", map[string]interface{}{
				"ok": false, "error": fmt.Sprintf("runDarijaScript() takes one string argument, %d given", len(args)),
			})
		}
		source := args[0].String()
		// Return immediately so Go timers (including n3ess) can resume in JS.
		executor := js.FuncOf(func(this js.Value, callbacks []js.Value) interface{} {
			resolve := callbacks[0]
			go func() {
				activeInput = &browserInput{lines: make(chan string, 1)}
				_, err := runtime.Run(source, interpreter.NewEvaluatorWithIO(activeInput, os.Stdout))
				activeInput = nil
				result := map[string]interface{}{"ok": err == nil, "error": ""}
				if err != nil {
					result["error"] = err.Error()
				}
				resolve.Invoke(result)
			}()
			return nil
		})
		promise := js.Global().Get("Promise").New(executor)
		executor.Release()
		return promise
	}))
	select {}
}
