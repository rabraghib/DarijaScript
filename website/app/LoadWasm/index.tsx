'use client';

import { useCallback, useEffect, useRef, useState } from 'react';

export type RuntimeStatus = 'loading' | 'ready' | 'running' | 'error';

export function useDarijaScript(writeOutput: (output: string) => void) {
  const [status, setStatus] = useState<RuntimeStatus>('loading');
  const [error, setError] = useState('');
  const output = useRef(writeOutput);
  const worker = useRef<Worker | null>(null);
  const timeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  const currentStatus = useRef<RuntimeStatus>('loading');

  useEffect(() => {
    output.current = writeOutput;
  }, [writeOutput]);

  const updateStatus = useCallback((next: RuntimeStatus) => {
    currentStatus.current = next;
    setStatus(next);
  }, []);

  const clearTimeoutHandle = useCallback(() => {
    if (timeout.current !== null) clearTimeout(timeout.current);
    timeout.current = null;
  }, []);

  const initializeWorker = useCallback(() => {
    clearTimeoutHandle();
    worker.current?.terminate();
    const instance = new Worker('/darijascript-worker.js');
    worker.current = instance;
    const fail = (message: string) => {
      if (worker.current !== instance) return;
      clearTimeoutHandle();
      instance.terminate();
      worker.current = null;
      setError(message);
      updateStatus('error');
      output.current('Error: ' + message + '\r\n');
    };
    instance.onerror = (event) =>
      fail(event.message || 'The runtime stopped unexpectedly.');
    instance.onmessage = (event) => {
      if (worker.current !== instance) return;
      const message = event.data;
      switch (message.type) {
        case 'ready':
          updateStatus('ready');
          break;
        case 'output':
          output.current(message.output);
          break;
        case 'done':
          clearTimeoutHandle();
          if (!message.ok) output.current('Error: ' + message.error + '\r\n');
          updateStatus('ready');
          break;
        case 'fatal':
          fail(message.error);
          break;
      }
    };
  }, [clearTimeoutHandle, updateStatus]);

  const start = useCallback(() => {
    updateStatus('loading');
    setError('');
    initializeWorker();
  }, [initializeWorker, updateStatus]);

  useEffect(() => {
    initializeWorker();
    return () => {
      clearTimeoutHandle();
      worker.current?.terminate();
      worker.current = null;
    };
  }, [initializeWorker, clearTimeoutHandle]);

  const stop = useCallback(() => {
    output.current('\r\nExecution stopped.\r\n');
    start();
  }, [start]);

  const run = useCallback(
    (code: string) => {
      if (currentStatus.current !== 'ready' || !worker.current) return;
      updateStatus('running');
      worker.current.postMessage({ type: 'run', code });
      timeout.current = setTimeout(() => {
        output.current('\r\nExecution timed out after 10 seconds.\r\n');
        start();
      }, 10_000);
    },
    [start, updateStatus],
  );

  return { status, error, run, stop, retry: start };
}
