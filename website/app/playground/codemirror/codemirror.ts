import { LanguageSupport } from '@codemirror/language';
import { darijaScriptLanguage } from './parser';
import { EditorView, Extension } from '@uiw/react-codemirror';
import { linter } from '@codemirror/lint';
import type { CompletionContext } from '@codemirror/autocomplete';
import type { SourceDiagnostic } from '../../LoadWasm';
import languageData from '../../../../vscode/language-data.json';

export default function DarijaScript(check: (source: string) => Promise<SourceDiagnostic[]> = async () => []) {
  return new LanguageSupport(darijaScriptLanguage, [
    darijaScriptLanguage.data.of({
      autocomplete: (context: CompletionContext) => {
        const word = context.matchBefore(/[\p{L}\p{M}\p{N}_]*/u);
        if (!word || (word.from === word.to && !context.explicit)) return null;
        const method = context.state.sliceDoc(Math.max(0, word.from - 1), word.from) === '.';
        const entries = method ? languageData.methods : [...languageData.keywords, ...languageData.builtins];
        return {
          from: word.from,
          options: entries.map(({ name, detail }) => ({ label: name, detail, type: method ? 'method' : languageData.keywords.some(k => k.name === name) ? 'keyword' : 'function' })),
        };
      },
    }),
    linter(async (view) => {
      try {
        const diagnostics = await check(view.state.doc.toString());
        return diagnostics.map((diagnostic) => {
          const line = view.state.doc.line(Math.max(1, Math.min(diagnostic.line, view.state.doc.lines)));
          const prefix = Array.from(line.text).slice(0, Math.max(0, diagnostic.column - 1)).join('');
          const from = Math.min(line.to, line.from + prefix.length);
          const width = Array.from(view.state.doc.sliceString(from, line.to))[0]?.length ?? 0;
          return { from, to: Math.min(line.to, from + width), severity: 'error' as const, message: diagnostic.message };
        });
      } catch { return []; }
    }, { delay: 400 }),
  ]);
}

export const FontSizeThemeExtension: Extension = [EditorView.theme({ '&': { fontSize: '115%' } })];
