import { fileURLToPath } from 'node:url';

/** @type {import('next').NextConfig} */
const nextConfig = {
  // Completion metadata is shared with the sibling VS Code extension.
  turbopack: { root: fileURLToPath(new URL('..', import.meta.url)) },
};

export default nextConfig;
