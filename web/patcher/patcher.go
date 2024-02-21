//go:build js
// +build js

package main

import (
	"bufio"
	"io"
	"log"
	"syscall/js"

	"github.com/pangbox/rugburn/slipstrm/embedded"
	"github.com/pangbox/rugburn/slipstrm/patcher"
)

func main() {
	c := make(chan struct{}, 0)
	js.Global().Call("patcherLoaded",
		js.FuncOf(Patch),
		js.FuncOf(UnpackOriginal),
		js.FuncOf(CheckOriginal),
		js.FuncOf(GetRugburnVersion),
		js.FuncOf(GetEmbeddedVersion),
	)
	<-c
}

// JS arguments:
// - input: Uint8Array
// - logCallback: (error: Error, line: string) => void
// JS return value:
// - Promise<Uint8Array>
func Patch(this js.Value, p []js.Value) interface{} {
	jsinput := p[0]
	logcb := p[1]

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			// Setup logging.
			logr, logw := io.Pipe()
			go func() {
				scanner := bufio.NewScanner(logr)
				for scanner.Scan() {
					logcb.Invoke(nil, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					logcb.Invoke(err.Error(), nil)
				}
			}()
			defer logw.Close()
			logger := log.New(logw, "", log.LstdFlags)

			input := make([]byte, jsinput.Get("length").Int())
			js.CopyBytesToGo(input, jsinput)

			output, err := patcher.Patch(logger, input, embedded.RugburnDLL, embedded.Version)

			jsoutput := js.Global().Get("Uint8Array").New(len(output))
			js.CopyBytesToJS(jsoutput, output)

			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
			} else {
				resolve.Invoke(js.ValueOf(jsoutput))
			}
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// JS arguments:
// - input: Uint8Array
// JS return value:
// - Promise<Uint8Array>
func UnpackOriginal(this js.Value, p []js.Value) interface{} {
	jsinput := p[0]

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			input := make([]byte, jsinput.Get("length").Int())
			js.CopyBytesToGo(input, jsinput)

			output, err := patcher.UnpackOriginal(input)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			jsoutput := js.Global().Get("Uint8Array").New(len(output))
			js.CopyBytesToJS(jsoutput, output)
			resolve.Invoke(js.ValueOf(jsoutput))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// JS arguments:
// - input: Uint8Array
// JS return value:
// - Promise<boolean>
func CheckOriginal(this js.Value, p []js.Value) interface{} {
	jsinput := p[0]

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]

		go func() {
			input := make([]byte, jsinput.Get("length").Int())
			js.CopyBytesToGo(input, jsinput)
			output := patcher.CheckOriginalData(input)
			resolve.Invoke(js.ValueOf(output))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// JS arguments:
// - input: Uint8Array
// JS return value:
// - Promise<string>
func GetRugburnVersion(this js.Value, p []js.Value) interface{} {
	jsinput := p[0]

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			input := make([]byte, jsinput.Get("length").Int())
			js.CopyBytesToGo(input, jsinput)
			output, err := patcher.GetRugburnVersion(input)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}
			resolve.Invoke(js.ValueOf(output))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// JS return value:
// - Promise<string>
func GetEmbeddedVersion(this js.Value, p []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		resolve.Invoke(js.ValueOf(embedded.Version))
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
	return js.ValueOf(embedded.Version)
}
