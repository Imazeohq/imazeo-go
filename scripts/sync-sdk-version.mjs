#!/usr/bin/env node
/** Keep version.go SDKVersion in sync after release-it bumps package.json. */
import { readFileSync, writeFileSync } from 'node:fs';
import { spawnSync } from 'node:child_process';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = join(dirname(fileURLToPath(import.meta.url)), '..');
const version = JSON.parse(readFileSync(join(root, 'package.json'), 'utf8')).version;
const path = join(root, 'version.go');
const source = readFileSync(path, 'utf8');
const next = source.replace(
  /SDKVersion\s*=\s*"[^"]+"/,
  `SDKVersion = "${version}"`,
);
if (!next.includes(`SDKVersion = "${version}"`)) {
  console.error('sync-sdk-version: failed to update version.go');
  process.exit(1);
}
writeFileSync(path, next, 'utf8');
const add = spawnSync('git', ['add', 'version.go', 'package.json'], {
  cwd: root,
  stdio: 'inherit',
  shell: false,
});
if (add.status !== 0) process.exit(add.status ?? 1);
console.log(`sync-sdk-version: ${version}`);
