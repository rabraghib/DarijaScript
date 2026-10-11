'use client';

import { ICodingExample } from './examples';
import CodeMirror, { EditorView } from '@uiw/react-codemirror';
import { darijaScriptTheme } from './codemirror/theme';
import DarijaScript, { FontSizeThemeExtension } from './codemirror/codemirror';
import { useEffect } from 'react';
import dynamic from 'next/dynamic';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';

const DynamicPlaygroundWithNoSSR = dynamic(() => import('./playground'), {
  ssr: false,
});

export default function PlaygroundManager({
  examples,
}: {
  examples: ICodingExample[];
}) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const starterKey = searchParams.get('starter');

  useEffect(() => {
    window.scrollTo(0, 0);
    if (starterKey !== null) {
      if (examples.every((example) => example.key !== starterKey)) {
        router.replace('/playground');
      }
    }
  }, [router, starterKey, examples]);

  if (starterKey !== null) {
    return (
      <DynamicPlaygroundWithNoSSR
        key={starterKey}
        examples={examples}
        starterKey={starterKey}
      />
    );
  }
  return (
    <div className="grid grid-rows-[auto,1fr] px-6 py-10 max-w-6xl mx-auto gap-6 items-center justify-center min-h-full w-full">
      <div className="text-center">
        <h1 className="text-4xl font-bold">Playground</h1>
        <p className="text-lg text-gray-600">Choose an example to run</p>
      </div>
      <div className="grid lg:grid-cols-2 gap-6">
        {examples.map((example) => (
          <ExampleCard key={example.key} example={example} />
        ))}
      </div>
    </div>
  );
}

function ExampleCard({
  example,
}: Readonly<{
  example: ICodingExample;
}>) {
  'use client';
  return (
    <div className="bg-gray-800 p-4 rounded-lg">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">{example.name}</h2>
        <Link
          href={`/playground?starter=${example.key}`}
          className="text-white bg-gradient-to-r from-blue-500 via-blue-600 to-blue-700 hover:bg-gradient-to-br font-medium rounded text-sm px-3 py-1 text-center"
        >
          Try it!
        </Link>
      </div>
      <CodeMirror
        value={example.code}
        lang="darijascript"
        className="h-44 w-full mt-4"
        theme={darijaScriptTheme}
        readOnly={true}
        data-enable-grammarly="false"
        extensions={[
          DarijaScript(),
          FontSizeThemeExtension,
          EditorView.lineWrapping,
        ]}
      />
    </div>
  );
}
