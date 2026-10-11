'use server';

export interface ICodingExample {
  name: string;
  key: string;
  code: string;
}

export async function getExamples(): Promise<ICodingExample[]> {
  return [
    {
      name: 'Input / Output',
      key: 'interactive-input',
      code: '9ayed name = dakhel("Chno smiytek? ");\n9ayed birthYear = rdo3adad(dakhel("Ina 3am tzaditi a " + name + "? "));\n\ngolih("");\ngolih("Smiytek " + name + " w 3omrek " + (2024-birthYear) + " sana");\n',
    },
    {
      name: 'Hello, World!',
      key: 'hello-world',
      code: '\ngolih("Salam 3likom!!");\n',
    },
    {
      name: 'FizzBuzz',
      key: 'fizz-buzz',
      code: '# FizzBuzz\n\nfizzBuzz(100);\n\nfonksyon fizzBuzz(n) {\n    9ayed i = 0;\n    ma7dBa9i(i < n) {\n        i = i + 1;\n        9ayed result = "";\n        ilakan (i % 3 == 0) {\n            result = result + "Fizz";\n        }\n        ilakan (i % 5 == 0) {\n            result = result + "Buzz";\n        }\n        ilamakanch (result == "") {\n            golih(result);\n        } sinn {\n            golih(i);\n        }\n    }\n}\n',
    },
    {
      name: 'Numerical Integration',
      key: 'numerical-integration',
      code: 'golih("Calculating the integral of a function using different methods:");\ngolih("f(x) = x^2");\ngolih("a = 0, b = 1");\ngolih("n = 1000");\ngolih("");\ngolih("Rectangles method: " + methodRectangles(0, 1, 1000));\ngolih("Trapezes method: " + methodTrapezes(0, 1, 1000));\ngolih("Simpson method: " + methodSimpson(0, 1, 1000));\n\nfonksyon f(x) {\n  rjje3 (x*x);\n}\n\nfonksyon methodRectangles(a, b, n) {\n  9ayed h = (b-a)/n;\n  9ayed sum = 0;\n  9ayed i = 0;\n  ma7dBa9i(i < n) {\n    9ayed xi = a + (i*h);\n    sum = sum + (f(xi));\n    i = i + 1;\n  }\n  rjje3 (sum*h);\n}\n\nfonksyon methodTrapezes(a, b, n) {\n  9ayed h = (b-a)/n;\n  9ayed sum = (f(a) + f(b))/2;\n  9ayed i = 1;\n  ma7dBa9i(i < n) {\n    9ayed xi = a + (i*h);\n    sum = sum + f(xi);\n    i = i + 1;\n  }\n  rjje3 (sum*h);\n}\n\nfonksyon methodSimpson(a, b, n) {\n  9ayed h = (b-a)/n;\n  9ayed sum = f(a) + f(b);\n  9ayed i = 1;\n  ma7dBa9i(i < n) {\n    9ayed xi = a + (i*h);\n    sum = sum + 2*f(xi);\n    i = i + 1;\n  }\n  i = 0;\n  ma7dBa9i(i < n) {\n    9ayed xi = a + ((i+1/2)*h);\n    sum = sum + 4*f(xi);\n    i = i + 1;\n  }\n  rjje3 (sum*h/6);\n}\n',
    },
    {
      name: 'Non Linear Equation',
      key: 'non-linear-equation',
      code: 'fonksyon f(x) { rjje3 x*x-7; }\nfonksyon df(x) { rjje3 2*x; }\nfonksyon g(x) { rjje3 (3*x+7/x)/4; }\n\n9ayed acceptableErr = 0.0001;\n\ngolih("Calculating the root of f(x) = x^2 - 7 using different methods:");\ngolih("1. Dicotomic method: " + dicothomy(2, 3, acceptableErr));\ngolih("2. Newton\'s method: " + newton(2, acceptableErr));\ngolih("3. Point fixe method: " + pointFixe(2, acceptableErr));\n\nfonksyon dicothomy(a,b,e) {\n  ma7dBa9i (((b-a)/2) > e) {\n    9ayed mid = (a+b)/2;\n    ilakan (f(mid) == 0) {\n      rjje3 mid;\n    } sinn {\n      ilakan (f(a) * f(mid) < 0) {\n        b = mid;\n      } sinn {\n        a = mid;\n      }\n    }\n  }\n  rjje3 (a+b)/2;\n}\n\nfonksyon newton(x0, e) {\n  9ayed x1 = x0 - (f(x0)/df(x0));\n  ma7dBa9i (abs(x1-x0) > e) {\n    x0 = x1;\n    x1 = x0 - (f(x0)/df(x0));\n  }\n  rjje3 x1;\n}\n\nfonksyon pointFixe(x0, e) {\n  9ayed x1 = g(x0);\n  ma7dBa9i (abs(x1-x0) > e) {\n    x0 = x1;\n    x1 = g(x0);\n  }\n  rjje3 x1;\n}\n',
    },
    {
      name: "Djikstra's Algorithm",
      key: 'djikstra',
      code: '\n9ayed graph = l(\n  l(0, 1, 5, 0, 0, 0),\n  l(1, 0, 0, 4, 2, 0),\n  l(5, 0, 0, 0, 3, 7),\n  l(0, 4, 0, 0, 5, 0),\n  l(0, 2, 3, 5, 0, 1),\n  l(0, 0, 7, 0, 1, 0)\n);\n\n9ayed result = djikstra(graph);\ngolihA9ssarPaths(result);\n\nfonksyon golihA9ssarPaths(result) {\n    9ayed i = 0;\n    ma7dBa9i(i < len(result)) {\n        golih("MIN[" + toLetter(i) + "]: " + ara(result, i, 0));\n        9ayed path = ""; \n        9ayed prev = ara(result, i, 1);\n        ma7dBa9i(prev != -1 && prev != 0) {\n            path = " -> " + toLetter(prev) + path;\n            prev = ara(result, prev, 1);\n        }\n        path = toLetter(prev) + path + " -> " + toLetter(i);\n        golih("Path: " + path);\n        golih("");\n        i = i + 1;\n    }\n}\n\n\nfonksyon djikstra(graph) {\n  9ayed start = 0;\n  9ayed visited = dirRow(len(graph), ghalt);\n  9ayed distance = dirRow(len(graph), l(infini(), -1));\n  atih(distance, start, l(0, 0));\n  \n  9ayed i = 0;\n  9ayed continue = s7i7;\n  \n  ma7dBa9i(i < len(graph) && continue) {\n    9ayed minNode = -1;\n    9ayed j = 0;\n    ma7dBa9i(j < len(graph)) {\n      ilamakanch (ara(visited, j)) {\n        ilakan (minNode == -1) {\n          minNode = j;\n        } sinn {\n          ilakan (ara(distance, j, 0) < ara(distance, minNode, 0)) {\n            minNode = j;\n          }\n        }\n      }\n      j = j + 1;\n    }\n    \n    ilakan (minNode == -1) {\n      continue = ghalt;\n    } sinn {\n      atih(visited, minNode, s7i7);\n      \n      9ayed neighbor = 0;\n      ma7dBa9i(neighbor < len(graph)) {\n        ilakan (ara(graph, minNode, neighbor) != 0) {\n          9ayed newDist = (ara(distance, minNode, 0)) + (ara(graph, minNode, neighbor));\n          \n          ilakan (newDist < ara(distance, neighbor, 0)) {\n            atih(distance, neighbor, l(newDist, minNode));\n          }\n        }\n        neighbor = neighbor + 1;\n      }\n      \n      i = i + 1;\n    }\n  }\n\n  rjje3 distance;\n}\n',
    },
  ];
}

// const rootPath = path.join(__dirname, '..', '..', '..');

// const examples = [
//   {
//     name: 'Hello, World!',
//     codePath: path.join(rootPath, 'examples/hello_world.ds'),
//   },
//   {
//     name: 'FizzBuzz',
//     codePath: path.join(rootPath, 'examples/fizz_buzz.ds'),
//   },
//   {
//     name: 'Numerical Integration',
//     codePath: path.join(rootPath, 'examples/integral.ds'),
//   },
//   {
//     name: 'Non Linear Equation',
//     codePath: path.join(rootPath, 'examples/non_linear_eqq.ds'),
//   },
// ];

// export async function getExamples(): Promise<ICodingExample[]> {
//   return await Promise.all(
//     examples.map(async (example) => {
//       const codeBuffer = await readFile(example.codePath);
//       return {
//         name: example.name,
//         code: codeBuffer.toString(),
//       };
//     })
//   );
// }
