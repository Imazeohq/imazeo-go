#!/usr/bin/env node
/**
 * Local release for the Imazeo Go SDK.
 *
 *   $env:RELEASE_MODE="true"; npm run release
 *
 * Interactive bump + GitHub Release via release-it.
 * Syncs package.json → version.go, then commit/tag/push.
 * Module publish is handled by .github/workflows/publish.yml.
 */

import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = join(dirname(fileURLToPath(import.meta.url)), '..');
const releaseIt = join(root, 'node_modules', 'release-it', 'bin', 'release-it.js');

function fail(message) {
  console.error(`release: ${message}`);
  process.exit(1);
}

function capture(command, args) {
  return spawnSync(command, args, {
    cwd: root,
    encoding: 'utf8',
    shell: false,
    env: process.env,
  });
}

function run(command, args) {
  const result = spawnSync(command, args, {
    cwd: root,
    stdio: 'inherit',
    shell: false,
    env: process.env,
  });
  if (result.error) {
    fail(`failed to start ${command}: ${result.error.message}`);
  }
  if (result.status !== 0) process.exit(result.status ?? 1);
}

function preflight() {
  const python = capture('python', ['--version']);
  if (python.error || python.status !== 0) {
    fail('python is required on PATH (used by scripts/check-release.py)');
  }

  const branch = capture('git', ['rev-parse', '--abbrev-ref', 'HEAD']);
  if (branch.status === 0 && (branch.stdout || '').trim() !== 'main') {
    fail(`release must run on main (current branch: ${(branch.stdout || '').trim()})`);
  }

  const dirty = capture('git', ['status', '--porcelain']);
  if (dirty.status === 0 && (dirty.stdout || '').trim()) {
    fail('working tree must be clean; commit or stash changes first');
  }
}

if (process.env.GITHUB_ACTIONS === 'true') {
  fail('CI publish is handled by .github/workflows/publish.yml');
}

if (process.env.RELEASE_MODE !== 'true') {
  fail('set RELEASE_MODE=true (PowerShell: $env:RELEASE_MODE="true")');
}

if (!existsSync(releaseIt)) {
  fail('release-it is not installed; run npm install');
}

preflight();
run(process.execPath, [releaseIt]);
