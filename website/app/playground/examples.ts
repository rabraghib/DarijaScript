'use server';

import examples from './examples.generated.json';

export interface ICodingExample {
  name: string;
  key: string;
  code: string;
}

export async function getExamples(): Promise<ICodingExample[]> {
  return examples;
}
