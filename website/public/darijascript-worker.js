/* The Go bootstrap and binary are generated together by just wasm. */
let pendingOutput = '';
let outputBytes = 0;
let overflow = false;
let running = false;

function flushOutput() {
  if (pendingOutput) {
    postMessage({ type: 'output', output: pendingOutput });
    pendingOutput = '';
  }
}

async function initialize() {
  importScripts('/wasm_exec.js');
  const go = new Go();
  const response = await fetch('/darijascript-bin-web.wasm');
  if (!response.ok) throw new Error('Could not load the interpreter (' + response.status + ').');
  const { instance } = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject);
  const decoder = new TextDecoder();
  globalThis.fs.writeSync = (_fd, buffer) => {
    outputBytes += buffer.length;
    if (outputBytes > 1_000_000) {
      if (!overflow) { flushOutput(); postMessage({ type: 'fatal', error: 'Output limit exceeded.' }); }
      overflow = true;
      return buffer.length;
    }
    pendingOutput += decoder.decode(buffer, { stream: true });
    flushOutput();
    return buffer.length;
  };
  go.run(instance).then(() => {
    postMessage({ type: 'fatal', error: 'The interpreter exited unexpectedly.' });
  }).catch((error) => postMessage({ type: 'fatal', error: error.message }));
  if (typeof globalThis.runDarijaScript !== 'function') throw new Error('The interpreter did not initialize.');
  postMessage({ type: 'ready' });
}

const ready = initialize().catch((error) => {
  postMessage({ type: 'fatal', error: error.message });
  return false;
});

onmessage = async ({ data }) => {
  if (data.type !== 'run' || running || await ready === false) return;
  running = true;
  outputBytes = 0;
  overflow = false;
  try {
    const result = await globalThis.runDarijaScript(data.code);
    flushOutput();
    postMessage({ type: 'done', ...result });
  } catch (error) {
    flushOutput();
    postMessage({ type: 'fatal', error: error.message });
  } finally {
    running = false;
  }
};
