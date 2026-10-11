'use client';

import { useCallback, useEffect, useRef, useState } from 'react';

export type RuntimeStatus = 'loading' | 'ready' | 'running' | 'waiting-input' | 'error';

export function useDarijaScript(writeOutput: (output: string) => void) {
  const [status, setStatus] = useState<RuntimeStatus>('loading');
  const [error, setError] = useState('');
  const output = useRef(writeOutput);
  const worker = useRef<Worker | null>(null);
  const timeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  const currentStatus = useRef<RuntimeStatus>('loading');
  const remainingTime = useRef(10_000);
  const deadline = useRef(0);

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
        case 'input-request':
          remainingTime.current = Math.max(0, deadline.current - Date.now());
          clearTimeoutHandle();
          updateStatus('waiting-input');
          break;
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

  const armTimeout = useCallback(() => {
    deadline.current = Date.now() + remainingTime.current;
    timeout.current = setTimeout(() => {
      output.current('\r\nExecution timed out after 10 seconds.\r\n');
      start();
    }, remainingTime.current);
  }, [start]);

  const submitInput = useCallback((value: string) => {
    if (currentStatus.current !== 'waiting-input' || !worker.current) return;
    output.current('\r\n');
    updateStatus('running');
    armTimeout();
    worker.current.postMessage({ type: 'input', value });
  }, [armTimeout, updateStatus]);

  const run = useCallback(
    (code: string) => {
      if (currentStatus.current !== 'ready' || !worker.current) return;
      updateStatus('running');
      remainingTime.current = 10_000;
      armTimeout();
      worker.current.postMessage({ type: 'run', code });
    },
    [armTimeout, updateStatus],
  );

  return { status, error, run, stop, submitInput, retry: start };
}
