//go:build js

package entrypoint

import (
	"fmt"
	"syscall/js"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/runtime"
)

func Entrypoint() {
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
				_, err := runtime.Run(source, interpreter.NewEvaluator())
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
