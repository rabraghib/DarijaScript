const vscode = require('vscode');
const { execFile } = require('node:child_process');
const path = require('node:path');
const data = require('./language-data.json');
const { extendMarkdownIt } = require('./markdown-preview');

function invoke(document, operation, args = []) {
  const executable = vscode.workspace.getConfiguration('darijascript', document.uri).get('interpreterPath', 'darijascript');
  return new Promise((resolve, reject) => {
    const child = execFile(executable, [operation, ...args], {
      cwd: document.uri.scheme === 'file' ? path.dirname(document.uri.fsPath) : undefined,
      timeout: 10_000,
      maxBuffer: 2_000_000,
      windowsHide: true,
    }, (error, stdout, stderr) => {
      if (error) reject(new Error(stderr || error.message));
      else resolve(stdout);
    });
    child.stdin.on('error', () => {}); // The child may exit before consuming invalid input.
    child.stdin.end(document.getText());
  });
}

function activate(context) {
  const selector = { language: 'darijascript' };
  const diagnostics = vscode.languages.createDiagnosticCollection('darijascript');
  const output = vscode.window.createOutputChannel('DarijaScript');
  const timers = new Map();
  let reportedError = '';
  context.subscriptions.push(diagnostics, output, { dispose() { for (const timer of timers.values()) clearTimeout(timer); timers.clear(); } });

  function schedule(document) {
    if (document.languageId !== 'darijascript') return;
    const key = document.uri.toString();
    clearTimeout(timers.get(key));
    timers.set(key, setTimeout(async () => {
      timers.delete(key);
      const version = document.version;
      try {
        const results = JSON.parse(await invoke(document, 'check', ['--json']));
        if (document.isClosed || document.version !== version) return;
        const items = results.map(item => {
          const lineIndex = Math.max(0, Math.min(document.lineCount - 1, item.line - 1));
          const text = document.lineAt(lineIndex).text;
          const column = Array.from(text).slice(0, Math.max(0, item.column - 1)).join('').length;
          const width = Array.from(text.slice(column))[0]?.length || 0;
          const range = new vscode.Range(lineIndex, column, lineIndex, Math.min(text.length, column + width));
          return new vscode.Diagnostic(range, item.message, vscode.DiagnosticSeverity.Error);
        });
        diagnostics.set(document.uri, items);
        reportedError = '';
      } catch (error) {
        if (!document.isClosed && document.version === version) diagnostics.delete(document.uri);
        if (reportedError !== error.message) {
          output.appendLine('Install the updated DarijaScript CLI or set darijascript.interpreterPath. ' + error.message);
          reportedError = error.message;
        }
      }
    }, 350));
  }

  context.subscriptions.push(
    vscode.commands.registerCommand('darijascript.openLanguageReference', () => {
      const reference = vscode.Uri.file(path.join(context.extensionPath, 'LANGUAGE.md'));
      return vscode.commands.executeCommand('markdown.showPreview', reference);
    }),
    vscode.workspace.onDidOpenTextDocument(schedule),
    vscode.workspace.onDidChangeTextDocument(event => schedule(event.document)),
    vscode.workspace.onDidChangeConfiguration(event => {
      if (event.affectsConfiguration('darijascript.interpreterPath')) vscode.workspace.textDocuments.forEach(schedule);
    }),
    vscode.workspace.onDidCloseTextDocument(document => {
      clearTimeout(timers.get(document.uri.toString()));
      timers.delete(document.uri.toString());
      diagnostics.delete(document.uri);
    }),
    vscode.languages.registerCompletionItemProvider(selector, {
      provideCompletionItems(document, position) {
        const prefix = document.lineAt(position.line).text.slice(0, position.character);
        const methods = /\.[\p{L}\p{M}\p{N}_]*$/u.test(prefix);
        const entries = methods ? data.methods : [...data.keywords, ...data.builtins];
        return entries.map(entry => {
          const kind = methods ? vscode.CompletionItemKind.Method : data.keywords.includes(entry) ? vscode.CompletionItemKind.Keyword : vscode.CompletionItemKind.Function;
          const item = new vscode.CompletionItem(entry.name, kind);
          item.detail = entry.detail;
          return item;
        });
      },
    }, '.'),
    vscode.languages.registerDocumentFormattingEditProvider(selector, {
      async provideDocumentFormattingEdits(document) {
        const version = document.version;
        const formatted = await invoke(document, 'format');
        if (document.version !== version) return [];
        const range = new vscode.Range(document.positionAt(0), document.positionAt(document.getText().length));
        return [vscode.TextEdit.replace(range, formatted)];
      },
    }),
  );
  vscode.workspace.textDocuments.forEach(schedule);
  return { extendMarkdownIt };
}

module.exports = { activate };
