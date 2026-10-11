const data = require('./language-data.json');

const keywords = new Set(data.keywords.map(entry => entry.name));
const builtins = new Set(data.builtins.map(entry => entry.name));
const literals = new Set(['s7i7', 'ghalt', 'true', 'false', 'walo']);

function highlight(code, escapeHtml) {
  // Consume strings and comments first so their contents stay literal.
  const tokens = /(?<comment>\/\/[^\r\n]*|#[^\r\n]*|\/\*[\s\S]*?(?:\*\/|$))|(?<string>"(?:\\[\s\S]|[^"\\])*(?:"|$))|(?<declaration>(?<![\p{L}\p{M}0-9_])9ayed(?![\p{L}\p{M}0-9_]))|(?<identifier>[\p{L}\p{M}_][\p{L}\p{M}0-9_]*)|(?<number>[0-9]+(?:\.[0-9]*)?)/gu;
  let html = '';
  let position = 0;
  for (const token of code.matchAll(tokens)) {
    html += escapeHtml(code.slice(position, token.index));
    const text = token[0];
    let type;
    if (token.groups.comment !== undefined) type = 'comment';
    else if (token.groups.string !== undefined) type = 'string';
    else if (token.groups.number !== undefined) type = 'number';
    else if (literals.has(text)) type = 'literal';
    else if (keywords.has(text) || /^WAA+$/.test(text)) type = 'keyword';
    else if (builtins.has(text)) type = 'built_in';
    const escaped = escapeHtml(text);
    html += type ? `<span class="hljs-${type}">${escaped}</span>` : escaped;
    position = token.index + text.length;
  }
  return html + escapeHtml(code.slice(position));
}

function extendMarkdownIt(md) {
  const previousHighlight = md.options.highlight;
  md.options.highlight = function (code, language, ...args) {
    if (/^(darijascript|ds)$/i.test(language)) {
      return highlight(code, md.utils.escapeHtml);
    }
    return previousHighlight ? previousHighlight.call(this, code, language, ...args) : '';
  };
  return md;
}

module.exports = { extendMarkdownIt };
