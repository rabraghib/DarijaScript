import { LRLanguage, LanguageSupport } from '@codemirror/language';
import { parserWithMetadata } from './parser';
import { EditorView, Extension } from '@uiw/react-codemirror';

export default function DarijaScript() {
  return new LanguageSupport(darijaScriptLanguage, [darijaScriptCompletion]);
}

const FontSizeTheme = EditorView.theme({
  '&': {
    fontSize: '115%',
  },
});
export const FontSizeThemeExtension: Extension = [FontSizeTheme];

const darijaScriptLanguage = LRLanguage.define({
  parser: parserWithMetadata,
  languageData: {
    commentTokens: { line: '#', block: { open: '/*', close: '*/' } },
  },
});

const darijaScriptCompletion = darijaScriptLanguage.data.of({
  autocomplete: [
    {
      label: 'golih',
      type: 'function',
      detail: 'Prints a message to the console',
    },
    {
      label: 'dakhel',
      type: 'function',
      detail: 'Reads a number from the console',
    },
    {
      label: 'abs',
      type: 'function',
      detail: 'Returns the absolute value of a number',
    },
    {
      label: 'rdo3adad',
      type: 'function',
      detail: 'Reads a number from the console',
    },
    {
      label: 'rdoBooleen',
      type: 'function',
      detail: 'Reads a boolean from the console',
    },
    {
      label: 'rdoString',
      type: 'function',
      detail: 'Reads a string from the console',
    },
    {
      label: 'ilakan',
      type: 'keyword',
    },
    {
      label: 'ilamakanch',
      type: 'keyword',
    },
    {
      label: 'sinn',
      type: 'keyword',
    },
    {
      label: 'ma7dBa9i',
      type: 'keyword',
    },
    {
      label: 'rjje3',
      type: 'keyword',
    },
    {
      label: 'WAA',
      type: 'keyword',
    },
    {
      label: '9ayed',
      type: 'keyword',
    },
    {
      label: 'fonksyon',
      type: 'keyword',
    },
  ],
});
