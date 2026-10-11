package runtime_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/runtime"
)

func execute(source, input string) (string, error) {
	var output bytes.Buffer
	_, err := runtime.Run(source, interpreter.NewEvaluatorWithIO(strings.NewReader(input), &output))
	return output.String(), err
}

func TestLanguagePrograms(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"collection literals and assignment", "9ayed a = [1, [2, 3],]; a[1][0] = 7; 9ayed d = {\"name\": \"Fatima\",}; d[\"age\"] = 25; golih(a); golih(d[\"name\"]); golih(d[\"missing\"]);", "[1, [7, 3]]\nFatima\nwalo\n"},
		{"aliases and resizing", "9ayed a = []; 9ayed b = a; b.append(1); b.append(2); golih(a); golih(b.remove(0)); golih(a);", "[1, 2]\n1\n[2]\n"},
		{"shallow slices", "9ayed a = [[1], 2]; 9ayed b = a.slice(0, 1); b.append(3); b[0][0] = 8; golih(a); golih(b);", "[[8], 2]\n[[8], 3]\n"},
		{"dictionary methods", "9ayed d = {\"b\": 2, \"a\": walo}; golih(d.keys()); golih(d.contains(\"a\")); golih(d.contains(\"c\")); golih(d.remove(\"b\")); golih(len(d));", "[\"a\", \"b\"]\ns7i7\nghalt\n2\n1\n"},
		{"content equality", "golih([1, {\"x\": 2}] == [1, {\"x\": 2}]); golih([1] == [2]); golih(walo == walo); golih(walo != 0); golih([1].contains(1));", "s7i7\nghalt\ns7i7\ns7i7\ns7i7\n"},
		{"cyclic lists", "9ayed a = [walo]; a[0] = a; 9ayed b = [walo]; b[0] = b; golih(a); golih(a == b);", "[<cycle>]\ns7i7\n"},
		{"cyclic dictionaries", "9ayed a = {}; a[\"self\"] = a; 9ayed b = {}; b[\"self\"] = b; golih(a); golih(a == b);", "{\"self\": <cycle>}\ns7i7\n"},
		{"mixed cycles", "9ayed d = {}; 9ayed a = [d]; d[\"list\"] = a; golih(d); golih(d == d);", "{\"list\": [<cycle>]}\ns7i7\n"},
		{"for break continue", "9ayed total = 0; kol (x f [1, 2, 3, 4, 5]) { ilakan (x == 2) { kmmel; } ilakan (x == 5) { w9ef; } total = total + x; } golih(total);", "8\n"},
		{"nested loop control", "9ayed total = 0; kol (x f [1, 2]) { kol (y f [3, 4]) { total = total + 1; w9ef; } } golih(total);", "2\n"},
		{"while control", "9ayed i = 0; ma7dBa9i (s7i7) { i = i + 1; ilakan (i < 3) { kmmel; } w9ef; } golih(i);", "3\n"},
		{"snapshot traversal", "9ayed a = [1, 2]; kol (x f a) { a.append(x); } golih(a);", "[1, 2, 1, 2]\n"},
		{"dictionary traversal", "kol (key f {\"z\": 1, \"a\": 2}) { golih(key); }", "a\nz\n"},
		{"unicode iteration", "kol (ch f \"سلام\") { golih(ch); }", "س\nل\nا\nم\n"},
		{"lexical scope", "9ayed x = 1; fonksyon read() { rjje3 x; } fonksyon caller() { 9ayed x = 2; rjje3 read(); } golih(caller());", "1\n"},
		{"returned closure", "fonksyon maker(x) { fonksyon inner() { rjje3 x; } rjje3 inner; } golih(maker(7)());", "7\n"},
		{"closure mutation", "fonksyon counter() { 9ayed n = 0; rjje3 fonksyon() { n = n + 1; rjje3 n; }; } 9ayed a = counter(); 9ayed b = counter(); golih(a()); golih(a()); golih(b());", "1\n2\n1\n"},
		{"loop closure capture", "9ayed funcs = []; kol (x f [1, 2, 3]) { funcs.append(fonksyon() { rjje3 x; }); } golih(funcs[0]()); golih(funcs[2]());", "1\n3\n"},
		{"function values", "fonksyon apply(fn, x) { rjje3 fn(x); } golih(apply(abs, -2.5)); 9ayed a = []; 9ayed add = a.append; add(4); golih(a);", "2.5\n[4]\n"},
		{"hoisting and recursion", "golih(factorial(5)); fonksyon factorial(n) { ilakan (n == 0) { rjje3 1; } rjje3 n * factorial(n - 1); }", "120\n"},
		{"f remains an identifier", "fonksyon f(x) { rjje3 x * x; } golih(f(3));", "9\n"},
		{"return through loops and try", "fonksyon find() { jerreb { kol (x f [1, 2]) { rjje3 x; } } ched (err) { rjje3 9; } } golih(find());", "1\n"},
		{"empty return", "fonksyon f() { rjje3; } golih(f());", "walo\n"},
		{"caught conversion", "jerreb { rdo3adad(\"bad\"); } ched (err) { golih(err.contains(\"cannot convert\")); } golih(3);", "s7i7\n3\n"},
		{"caught throw", "jerreb { WAAA \"oops\"; } ched (err) { golih(err.contains(\"oops\")); }", "s7i7\n"},
		{"catch scope and recovery", "9ayed x = 1; jerreb { 9ayed x = 2; 1 / 0; } ched (err) { golih(x); } golih(x);", "1\n1\n"},
		{"break through try", "kol (x f [1, 2]) { jerreb { golih(x); w9ef; } ched (err) { golih(\"bad\"); } }", "1\n"},
		{"string semantics", "golih(\"1\" + \"2\"); golih(\"age: \" + 25); golih(rdo3adad(\"1\") + rdo3adad(\"2\")); golih(rdoString(ghalt));", "12\nage: 25\n3\nghalt\n"},
		{"unicode strings", "golih(len(\"سلام\")); golih(\"سلام\"[1]); golih(\"سلام\".slice(1, 3)); golih(\"سلام\".indexOf(\"ام\")); golih(\"سلام\".contains(\"لا\"));", "4\nل\nلا\n2\ns7i7\n"},
		{"string methods", "golih(\" a,b \".trim().split(\",\").join(\"-\")); golih(\"سلام\".split(\"\").join(\"|\")); golih(\"hello\".substring(1)); golih(\"hello\".indexOf(\"z\"));", "a-b\nس|ل|ا|م\nello\n-1\n"},
		{"escapes", `golih("quote: \" slash: \\ tab:\tnewline:\nend");`, "quote: \" slash: \\ tab:\tnewline:\nend\n"},
		{"empty collections", "golih([]); golih({}); golih([].slice(0)); golih(\"\".slice(0)); golih(\"\".split(\"\"));", "[]\n{}\n[]\n\n[]\n"},
		{"math", "golih(abs(-0.25)); golih(sqrt(9)); golih(pow(2, 3)); golih(floor(1.8)); golih(ceil(1.2)); golih(round(-1.5)); golih(min(3, -1, 2)); golih(max(3, -1, 2)); 9ayed n = random(); golih(n >= 0 && n < 1);", "0.25\n3\n8\n1\n2\n-2\n-1\n3\ns7i7\n"},
		{"short circuit", "golih(ghalt && missing()); golih(s7i7 || missing());", "ghalt\ns7i7\n"},
		{"assignment evaluates target once", "9ayed calls = 0; 9ayed a = [1]; fonksyon get() { calls = calls + 1; rjje3 a; } get()[0] = 5; golih(calls); golih(a);", "1\n[5]\n"},
		{"dictionary literal order", "9ayed n = 0; fonksyon next() { n = n + 1; rjje3 n; } 9ayed d = {\"a\": next(), \"b\": next(), \"a\": next()}; golih(d);", "{\"a\": 3, \"b\": 2}\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := execute(tc.source, "")
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q; want %q", got, tc.want)
			}
		})
	}
}

func TestLanguageErrors(t *testing.T) {
	cases := []struct{ source, want string }{
		{"[1][-1];", "expected an integer"},
		{"[1][1];", "expected an integer"},
		{"[1][0.5];", "expected an integer"},
		{"[1][s7i7];", "number"},
		{"[1][\"0\"];", "number"},
		{"{\"a\": 1}[0];", "keys must be strings"},
		{"{1: 2};", "keys must be strings"},
		{"9ayed a = \"hi\"; a[0] = \"x\";", "indexed assignment"},
		{"[].remove(0);", "expected an integer"},
		{"[1].slice(1, 0);", "slice end"},
		{"[1].slice(0, 2);", "expected an integer"},
		{"[1].join(\",\");", "expected a string"},
		{"{}.contains(2);", "expected a string"},
		{"\"abc\".split(1);", "expected a string"},
		{"\"abc\".unknown();", "unknown method"},
		{"[].append();", "wrong number"},
		{"\"1\" - 1;", "number"},
		{"sqrt(-1);", "numeric domain"},
		{"pow(10, 1000);", "numeric domain"},
		{"min();", "wrong number"},
		{"random(1);", "wrong number"},
		{"rdo3adad(\"NaN\");", "finite number"},
		{"w9ef;", "only valid inside a loop"},
		{"kmmel;", "only valid inside a loop"},
		{"rjje3 1;", "only valid inside a function"},
		{"kol (x f []) { fonksyon bad() { w9ef; } }", "only valid inside a loop"},
		{"kol (x f 1) {}", "requires a list"},
		{"fonksyon f(x, x) {}", "duplicate parameter"},
		{"9ayed f = 1; f();", "cannot call"},
		{"fonksyon f(x) {} f();", "wrong number"},
		{"9ayed a = 1; 9ayed a = 2;", "already declared"},
		{"l(1);", "identifier not found"},
		{"c(1);", "identifier not found"},
		{"r(1);", "identifier not found"},
		{"ara([], 0);", "identifier not found"},
		{"atih([], 0, 1);", "identifier not found"},
		{"dirRow(1, 0);", "identifier not found"},
		{"dirCol(1, 0);", "identifier not found"},
		{`"bad\q";`, "unknown string escape"},
		{"[1, 2;", "Expected token"},
		{"jerreb {} ched {}", "Expected token"},
		{"fonksyon recurse() { recurse(); } recurse();", "maximum evaluation nesting"},
	}
	for _, tc := range cases {
		t.Run(tc.source, func(t *testing.T) {
			_, err := execute(tc.source, "")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v; want error containing %q", err, tc.want)
			}
		})
	}
}

func TestSourceDiagnostic(t *testing.T) {
	source := "fonksyon divide() {\n    rjje3 1 / 0;\n}\nfonksyon caller() { divide(); }\ncaller();"
	_, err := runtime.RunNamed(source, "program.ds", interpreter.NewEvaluatorWithIO(strings.NewReader(""), &bytes.Buffer{}))
	if err == nil {
		t.Fatal("expected failure")
	}
	for _, want := range []string{"program.ds:2:13:", "rjje3 1 / 0;", "^", "called from divide at 4:", "called from caller at 5:"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("missing %q in %s", want, err)
		}
	}
}

func TestInputRecovery(t *testing.T) {
	got, err := execute("jerreb { rdo3adad(dakhel()); } ched (err) { golih(\"invalid\"); } golih(rdo3adad(dakhel()) + 1);", "bad\n2\n")
	if err != nil || got != "invalid\n3\n" {
		t.Fatalf("got %q, %v", got, err)
	}
}

func TestIncomplete(t *testing.T) {
	for _, source := range []string{"fonksyon f() {", "9ayed a = [1,", "9ayed x =", "jerreb {} ched (err) {", "\"hello", "/* comment"} {
		if !runtime.Incomplete(source) {
			t.Errorf("should be incomplete: %q", source)
		}
	}
	for _, source := range []string{"golih(1);", "w9ef;", "9ayed a = [1;"} {
		if runtime.Incomplete(source) {
			t.Errorf("should not be incomplete: %q", source)
		}
	}
}

func TestFormatting(t *testing.T) {
	source := "# keep this\n9ayed a=[1,{\"name\":\"سلام\", \"escape\":\"\\\"\\\\\\n\"}];\n/* keep block */\nkol(x f a){jerreb{golih(x);}ched(err){golih(err);}}\nfonksyon maker(x){rjje3 fonksyon(){rjje3 x;};}golih(maker(7)());"
	formatted, err := runtime.Format(source)
	if err != nil {
		t.Fatal(err)
	}
	twice, err := runtime.Format(formatted)
	if err != nil || formatted != twice {
		t.Fatalf("formatting not idempotent: %q vs %q (%v)", formatted, twice, err)
	}
	before, err := execute(source, "")
	if err != nil {
		t.Fatal(err)
	}
	after, err := execute(formatted, "")
	if err != nil || before != after {
		t.Fatalf("formatting changed behavior: %q vs %q (%v)", before, after, err)
	}
	if !strings.Contains(formatted, "# keep this") || !strings.Contains(formatted, "/* keep block */") {
		t.Fatal("lost comments")
	}
}

func TestExamples(t *testing.T) {
	paths, err := filepath.Glob("../../examples/*.ds")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("missing examples")
	}
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			input := "Rabyâ\n2000\n"
			if strings.HasSuffix(path, "calculator.ds") {
				input = "+\n1\n2\n"
			}
			if _, err := execute(string(source), input); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func FuzzParser(f *testing.F) {
	for _, source := range []string{"", "[];", "9ayed x = {\"a\": [1]};", "kol (x f []) {}", "jerreb {} ched (err) {}", "fonksyon f() { rjje3; }"} {
		f.Add(source)
	}
	f.Fuzz(func(t *testing.T, source string) { runtime.Check(source) })
}
