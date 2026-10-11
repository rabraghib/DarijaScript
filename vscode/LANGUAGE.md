# DarijaScript language reference

Write UTF-8 programs in `.ds` files. Keywords are case-sensitive.

## Quick start

```ds
9ayed name = dakhel("Chno smiytek? ");
golih("Salam " + name);
```

## Syntax and values

- End statements with `;`. Blocks and named function declarations need no trailing `;`.
- Comments: `# line`, `// line`, `/_ block _/`.
- Identifiers support Unicode letters, marks, `_`, and digits after the first character; `9ayed` is a special keyword.

| Type       | Example                                |
| ---------- | -------------------------------------- |
| Number     | `42`, `3.5`, `-2`                      |
| String     | `"Salam"`, `"سلام"`                    |
| Boolean    | `s7i7`, `ghalt` (also `true`, `false`) |
| No value   | `walo`                                 |
| List       | `[1, 2, 3]`                            |
| Dictionary | `{"ism": "Fatima"}`                    |
| Function   | `fonksyon(x) { rjje3 x + 1; }`         |

- Declare with `9ayed x = 1;`, assign with `x = 2;`. Blocks have their own scope; same-scope redeclaration is an error.
- Arithmetic requires numbers. `+` concatenates if either operand is a string. Convert input with `rdo3adad()`.
- Conditions require booleans. `&&` and `||` short-circuit.
- Operators, weakest to strongest: `=`, `||`, `&&`, `== !=`, `< <= > >=`, `+ -`, `* / %`, unary `! - +`, calls/indexing/methods.
- Numbers use floating-point arithmetic. Division by zero is an error; `%` requires integers and a nonzero divisor.
- Collections compare by contents and share mutations through aliases. Strings are immutable.
- Missing dictionary keys and functions without a return value produce `walo`.

## Control flow

```ds
9ayed score = 12;
ilakan (score >= 10) { golih("Passed"); }
sinn ilakan (score >= 5) { golih("Try again"); }
sinn { golih("Practice"); }
ilamakanch (score >= 10) { golih("Not passed"); }

9ayed i = 0;
ma7dBa9i (i < 3) { i = i + 1; }

kol (item f [1, 2, 3]) {
    ilakan (item == 2) { kmmel; } # Continue.
    ilakan (item == 3) { w9ef; }  # Break.
    golih(item);
}
```

`kol` visits a snapshot of list values, sorted dictionary keys, or string code points. `w9ef` and `kmmel` affect the nearest loop.

## Functions

```ds
fonksyon add(a, b) { rjje3 a + b; }
9ayed double = fonksyon(x) { rjje3 x * 2; };
golih(add(2, double(3))); # 8
```

Named functions are available throughout their block. Functions capture their defining scope and can be returned, stored, or passed as arguments. Argument counts must match; `rjje3` is only valid inside functions.

## Collections and strings

```ds
9ayed scores = [10, 20];
scores.append(30);
scores[0] = 15;
9ayed person = {"ism": "Fatima"};
person["mdina"] = "Rabat";
golih(person["ism"]);
```

Indices start at zero; negative or out-of-range indices are errors. Dictionary keys are strings. Slices use an exclusive end and create shallow copies.

| Method                                                    | Result                                        |
| --------------------------------------------------------- | --------------------------------------------- |
| `list.append(value)`                                      | Adds an item; returns `walo`                  |
| `list.remove(index)`                                      | Removes and returns an item                   |
| `list.slice(start, end?)`                                 | New list; end defaults to length              |
| `list.contains(value)`                                    | Checks content equality                       |
| `listOfStrings.join(separator)`                           | Joins strings                                 |
| `dict.keys()`                                             | Sorted keys                                   |
| `dict.contains(key)`                                      | Whether the key exists                        |
| `dict.remove(key)`                                        | Removed value, or `walo`                      |
| `text.trim()`                                             | Removes surrounding whitespace                |
| `text.split(separator)`                                   | List of strings; `""` splits into code points |
| `text.slice(start, end?)` / `text.substring(start, end?)` | Substring                                     |
| `text.contains(part)`                                     | Whether the substring exists                  |
| `text.indexOf(part)`                                      | First index, or `-1`                          |

Strings use double quotes; escapes: `\"`, `\\`, `\n`, `\r`, `\t`. Raw multiline strings are allowed. Length, indexing, slicing, and iteration count Unicode code points: `len("سلام")` is `4`.

## Errors

```ds
jerreb {
    WAA "Something went wrong";
} ched (err) {
    golih(err);
}
```

`ched` receives a runtime error as text. Syntax errors cannot be caught; earlier mutations are not rolled back.

## Built-ins

| Function                                   | Purpose                                                  |
| ------------------------------------------ | -------------------------------------------------------- |
| `golih(value)`                             | Print a value and newline                                |
| `dakhel(prompt?)`                          | Read a text line                                         |
| `rdo3adad(value)`                          | Convert a number, boolean, or numeric string to a number |
| `rdoBooleen(value)`                        | Convert a boolean, number, or boolean string             |
| `rdoString(value)`                         | Convert to displayed text                                |
| `len(value)`                               | List size, dictionary size, or string length             |
| `abs(n)`, `sqrt(n)`, `pow(base, exponent)` | Absolute value, square root, power                       |
| `floor(n)`, `ceil(n)`, `round(n)`          | Rounding; `round` ties go away from zero                 |
| `min(n, ...)`, `max(n, ...)`               | Minimum/maximum; at least one argument                   |
| `random()`                                 | Random number in `[0, 1)`                                |
| `infini()`                                 | Positive infinity                                        |
| `rdoChar(index)`                           | A–Z for integer indices 0–25                             |
| `n3ess(milliseconds)`                      | Sleep; integer duration from 0 to 86,400,000             |

Invalid conversions and math domain errors raise runtime errors. Math helpers require finite numeric arguments.

## Tools

```sh
darijascript run program.ds
darijascript interactive
darijascript check program.ds
darijascript format --write program.ds
```

`check` and `format` also accept stdin. Formatting prints to stdout unless `--write` is used. In interactive mode, enter `exit` to quit; incomplete input continues at `...`.

In VS Code, run **DarijaScript: Open Language Reference** from the Command Palette. Highlighting and completions work immediately; diagnostics and formatting need the CLI on `PATH` or `darijascript.interpreterPath`.

## Older programs

| Removed helper                            | Replacement                              |
| ----------------------------------------- | ---------------------------------------- |
| `l(1, 2)`, `c(1, 2)`, `r(1, 2)`           | `[1, 2]`                                 |
| `ara(items, i)` / `ara(matrix, row, col)` | `items[i]` / `matrix[row][col]`          |
| `atih(items, i, value)`                   | `items[i] = value`                       |
| `atih(matrix, row, col, value)`           | `matrix[row][col] = value`               |
| `dirRow(n, value)`, `dirCol(n, value)`    | Build a list with a loop and `.append()` |

Convert `dakhel()` results before arithmetic. Numeric strings concatenate under `+`. File imports are not supported.
