import { readFile, writeFile } from 'node:fs/promises';

const entries = [
  ['Input / Output', 'interactive-input', 'io.ds'],
  ['Hello, World!', 'hello-world', 'hello_world.ds'],
  ['Collections, closures, and errors', 'collections', 'collections.ds'],
  ['FizzBuzz', 'fizz-buzz', 'fizz_buzz.ds'],
  ['Numerical Integration', 'numerical-integration', 'integral.ds'],
  ['Non Linear Equation', 'non-linear-equation', 'non_linear_eqq.ds'],
  ["Dijkstra's Algorithm", 'djikstra', 'djikstra.ds'],
  ['Calculator', 'calculator', 'calculator.ds'],
];
const examples = await Promise.all(entries.map(async ([name, key, file]) => ({
  name, key, code: (await readFile(new URL('../../examples/' + file, import.meta.url), 'utf8')).replace(/\r\n/g, '\n'),
})));
const content = JSON.stringify(examples, null, 2) + '\n';
const target = new URL('../app/playground/examples.generated.json', import.meta.url);
if (process.argv.includes('--check')) {
  if (await readFile(target, 'utf8') !== content) {
    throw new Error('Playground examples are stale. Run npm --prefix website run sync-examples.');
  }
} else {
  await writeFile(target, content);
}
