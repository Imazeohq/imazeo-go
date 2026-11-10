#!/usr/bin/env node
/** Local pre-release checks. Go tests run only when the Go toolchain is available. */
import { spawnSync } from 'node:child_process';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = join(dirname(fileURLToPath(import.meta.url)), '..');

function fail(message) {
  console.error(`pre-release-check: ${message}`);
  process.exit(1);
}

function run(command, args) {
  const result = spawnSync(command, args, {
    cwd: root,
    stdio: 'inherit',
    shell: false,
    env: process.env,
  });
  if (result.error) fail(`failed to start ${command}: ${result.error.message}`);
  if (result.status !== 0) process.exit(result.status ?? 1);
}

run('python', ['scripts/check-release.py']);

const goVersion = spawnSync('go', ['version'], {
  cwd: root,
  encoding: 'utf8',
  shell: false,
  env: process.env,
});

if (goVersion.error || goVersion.status !== 0) {
  console.log('pre-release-check: skipping go test');
  process.exit(0);
}

run('go', ['test', './...']);
