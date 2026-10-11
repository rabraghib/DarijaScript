import { StreamLanguage } from '@codemirror/language';
import languageData from '../../../../vscode/language-data.json';

const keywords = new Set(languageData.keywords.map(({ name }) => name));
const builtins = new Set(languageData.builtins.map(({ name }) => name));
const identifier = /[\p{L}\p{M}_][\p{L}\p{M}0-9_]*/u;
interface State {
  blockComment: boolean;
  string: boolean;
  depth: number;
  definition: 'variable' | 'function' | null;
  property: boolean;
  foreach: 'open' | 'binding' | 'separator' | null;
}

// Syntax diagnostics still come from the Go parser in the worker.
export const darijaScriptLanguage = StreamLanguage.define<State>({
  startState: () => ({
    blockComment: false,
    string: false,
    depth: 0,
    definition: null,
    property: false,
    foreach: null,
  }),
  token(stream, state) {
    if (state.blockComment) {
      if (stream.match(/.*?\*\//)) state.blockComment = false;
      else stream.skipToEnd();
      return 'comment';
    }
    if (state.string || stream.peek() === '"') {
      if (!state.string) {
        stream.next();
        state.string = true;
        state.definition = null;
        state.property = false;
      } else if (stream.peek() === '\\') {
        stream.next();
        if (stream.match(/["\\nrt]/)) return 'escape';
        stream.next();
        return 'invalid';
      }
      while (!stream.eol()) {
        if (stream.peek() === '\\') break;
        if (stream.next() === '"') {
          state.string = false;
          break;
        }
      }
      return 'string';
    }
    if (stream.eatSpace()) return null;
    if (stream.match('/*')) {
      state.blockComment = !stream.match(/.*?\*\//);
      if (state.blockComment) stream.skipToEnd();
      return 'comment';
    }
    if (stream.match('//') || stream.match('#')) {
      stream.skipToEnd();
      return 'comment';
    }

    // Unlike JavaScript's \b, this boundary includes Arabic letters and marks.
    const word =
      stream.match(/9ayed(?![\p{L}\p{M}0-9_])/u) || stream.match(identifier);
    if (Array.isArray(word)) {
      const name = word[0];
      const definition = state.definition;
      const property = state.property;
      state.definition = null;
      state.property = false;
      if (state.foreach === 'binding') {
        state.foreach = 'separator';
        return 'variableName.definition';
      }
      if (state.foreach === 'separator') {
        state.foreach = null;
        if (name === 'f') return 'keyword';
      }
      if (name === '9ayed' || name === 'fonksyon') {
        state.definition = name === '9ayed' ? 'variable' : 'function';
        return 'definitionKeyword';
      }
      if (name === 'kol') state.foreach = 'open';
      if (name === 'walo') return 'null';
      if (['s7i7', 'ghalt', 'true', 'false'].includes(name)) return 'bool';
      if (keywords.has(name) || /^WAA+$/.test(name)) return 'controlKeyword';
      if (definition === 'function') return 'variableName.function.definition';
      if (definition === 'variable') return 'variableName.definition';
      if (property) return 'propertyName.function';
      if (builtins.has(name)) return 'variableName.function.standard';
      if (stream.match(/^\s*\(/, false)) return 'variableName.function';
      return 'variableName';
    }
    state.definition = null;
    state.property = false;
    if (stream.match(/[0-9]+(?:\.[0-9]*)?/)) return 'number';
    if (stream.match(/&&|\|\||[=!<>]=?|[+*/%\-]/)) return 'operator';
    const char = stream.next();
    if (state.foreach === 'open')
      state.foreach = char === '(' ? 'binding' : null;
    else state.foreach = null;
    if (char === '{' || char === '[' || char === '(') state.depth++;
    if (char === '}' || char === ']' || char === ')')
      state.depth = Math.max(0, state.depth - 1);
    if (char === '.') {
      state.property = true;
      return 'derefOperator';
    }
    if (char === '[' || char === ']') return 'squareBracket';
    if (char === '{' || char === '}') return 'brace';
    if (char === '(' || char === ')') return 'paren';
    if (char && ',:;'.includes(char)) return 'punctuation';
    return 'invalid';
  },
  indent(state, textAfter, context) {
    const closing = /^\s*[}\])]/.test(textAfter) ? 1 : 0;
    return Math.max(0, state.depth - closing) * context.unit;
  },
  languageData: {
    commentTokens: { line: '#', block: { open: '/*', close: '*/' } },
  },
});
