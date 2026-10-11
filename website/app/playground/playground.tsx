'use client';

import CodeMirror, { EditorView } from '@uiw/react-codemirror';
import { useCallback, useEffect, useRef, useState } from 'react';
import { Terminal } from '@xterm/xterm';

import '@xterm/xterm/css/xterm.css';
import { FitAddon } from '@xterm/addon-fit';
import DarijaScript, { FontSizeThemeExtension } from './codemirror/codemirror';
import { ICodingExample } from './examples';
import { githubDarkInit } from '@uiw/codemirror-theme-github';
import { useDarijaScript, RuntimeStatus } from '../LoadWasm';
import { useRouter } from 'next/navigation';

import { tags } from '@lezer/highlight';

export default function PlaygroundPage({
  examples,
  starterKey,
}: Readonly<{
  starterKey: string;
  examples: ICodingExample[];
}>) {
  const terminal = useRef<Terminal | null>(null);
  const terminalRef = useRef<HTMLDivElement>(null);
  const [code, setCode] = useState<string>(
    examples.find((example) => example.key === starterKey)?.code ?? ``
  );
  const writeOutput = useCallback((output: string) => {
    terminal.current?.write(output);
  }, []);
  const runtime = useDarijaScript(writeOutput);
  const { status, submitInput } = runtime;

  useEffect(() => {
    if (!terminalRef.current) return;
    const instance = new Terminal({ convertEol: true, lineHeight: 1.5, disableStdin: true });
    terminal.current = instance;
    instance.open(terminalRef.current);
    const fitAddon = new FitAddon();
    instance.loadAddon(fitAddon);
    fitAddon.fit();
    const resize = new ResizeObserver(() => fitAddon.fit());
    resize.observe(terminalRef.current);
    return () => {
      resize.disconnect();
      instance.dispose();
      terminal.current = null;
    };
  }, []);

  useEffect(() => {
    const instance = terminal.current;
    if (!instance) return;
    const waiting = status === 'waiting-input';
    instance.options.disableStdin = !waiting;
    instance.options.cursorBlink = waiting;
    instance.write(waiting ? '\x1b[?25h' : '\x1b[?25l');
    if (!waiting) return;

    let value = '';
    let accepting = true;
    instance.write('\x1b[s');
    instance.focus();
    const subscription = instance.onData((data) => {
      if (!accepting) return;
      if (data === '\r' || data === '\n') {
        accepting = false;
        instance.options.disableStdin = true;
        submitInput(value);
        return;
      }
      if (data === '\x7f' || data === '\b') {
        value = Array.from(value).slice(0, -1).join('');
      } else {
        // Ignore navigation/control sequences; paste remains one input line.
        if (data.includes('\x1b')) return;
        value += data.replace(/\r\n|[\r\n\t]/g, ' ').replace(/[\x00-\x1f\x7f-\x9f]/g, '');
      }
      instance.write('\x1b[u\x1b[J' + value);
    });
    return () => subscription.dispose();
  }, [status, submitInput]);

  function runCode() {
    terminal.current?.clear();
    terminal.current?.writeln('\x1b[32m$ darijascript run code.ds\x1b[0m');
    runtime.run(code);
  }

  return (
    <>
      <div className="h-[calc(100vh-4rem)] w-full grid grid-rows-2 lg:grid-cols-2 lg:grid-rows-1">
        <div className="w-full h-full border-r border-slate-500">
          <PlaygroundSection
            title={
              <CodeEditorTitle examples={examples} starterKey={starterKey} />
            }
            actions={<CodeEditorActions onRun={runCode} onStop={runtime.stop} onRetry={runtime.retry} status={runtime.status} />}
          >
            <AppCodeEditor
              value={code}
              onChange={(value) => {
                setCode(value);
              }}
            />
          </PlaygroundSection>
        </div>
        <PlaygroundSection title={<>Terminal</>}>
          {runtime.error && <p role="alert" className="px-4 text-red-400">{runtime.error}</p>}
          <div className="h-full flex flex-col bg-black">
            <div className="flex-1 min-h-0 w-full" ref={terminalRef}></div>
          </div>
        </PlaygroundSection>
      </div>
    </>
  );
}

function CodeEditorTitle({
  examples,
  starterKey,
}: Readonly<{
  starterKey: string;
  examples: ICodingExample[];
}>) {
  const router = useRouter();
  return (
    <div className="flex items-center space-x-2">
      <button
        aria-label="Back to examples"
        onClick={() => router.push('/playground')}
        className="border rounded p-1 text-center inline-flex items-center border-slate-500 hover:bg-slate-900 hover:bg-opacity-50"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          className="w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15.75 19.5 8.25 12l7.5-7.5"
          />
        </svg>
      </button>
      <select
        aria-label="Example"
        value={starterKey}
        onChange={(e) => {
          const starterKey = e.target.value;
          if (starterKey) {
            router.push(`/playground?starter=${starterKey}`);
          }
        }}
        className="bg-slate-950 text-white border border-slate-500 rounded p-1.5 text-sm min-w-52"
      >
        {examples.map((example) => (
          <option key={example.key} value={example.key}>
            {example.name}
          </option>
        ))}
      </select>
    </div>
  );
}

function CodeEditorActions({
  onRun,
  onStop,
  onRetry,
  status,
}: Readonly<{
  onRun: () => void;
  onStop: () => void;
  onRetry: () => void;
  status: RuntimeStatus;
}>) {
  return (
    <div className="flex space-x-4">
      <button
        onClick={onRun}
        disabled={status !== 'ready'}
        className="text-white bg-gradient-to-r from-blue-500 via-blue-600 to-blue-700 hover:bg-gradient-to-br disabled:opacity-50 font-medium rounded text-sm px-3 py-1 text-center"
      >
        {status === 'loading' ? 'Loading…' : 'Run'}
      </button>
      {(status === 'running' || status === 'waiting-input') && <button onClick={onStop}>Stop</button>}
      {status === 'error' && <button onClick={onRetry}>Retry</button>}
    </div>
  );
}

function AppCodeEditor({
  value,
  onChange,
}: Readonly<{
  value: string;
  onChange: (value: string) => void;
}>) {
  return (
    <CodeMirror
      value={value}
      lang="darijascript"
      className="h-full w-full"
      theme={githubDarkInit({
        styles: [
          {
            tag: tags.definitionKeyword,
            color: '#b972cf',
          },
          {
            tag: tags.operator,
            color: '#c9d1d9',
          },
          {
            tag: tags.string,
            color: '#93bd76',
          },
          {
            tag: tags.number,
            color: '#b1855d',
          },
          // {
          //   tag: [tags.function(tags.variableName)],
          //   color: 'yellow',
          // },
          // { tag: tags.definition(tags.variableName), color: 'yellow' },
        ],
      })}
      onChange={(value) => {
        onChange(value);
      }}
      data-enable-grammarly="false"
      extensions={[
        DarijaScript(),
        FontSizeThemeExtension,
        EditorView.lineWrapping,
      ]}
    />
  );
}

function PlaygroundSection({
  title,
  actions,
  children,
}: Readonly<{
  title: React.ReactNode;
  actions?: React.ReactNode;
  children: React.ReactNode;
}>) {
  return (
    <section className="h-full bg-slate-950 grid grid-rows-[auto,1fr] overflow-hidden">
      <div className="py-2 px-4 border-b border-slate-400 bg-black flex justify-between items-center h-14">
        <h2 className="text-lg font-medium">{title}</h2>
        {actions && <div className="flex justify-end space-x-4">{actions}</div>}
      </div>
      <div className="size-full overflow-auto">{children}</div>
    </section>
  );
}
