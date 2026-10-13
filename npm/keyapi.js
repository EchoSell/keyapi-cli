#!/usr/bin/env node
const { spawnSync } = require('node:child_process');
const path = require('node:path');

const { binaryNameForPlatform } = require('../scripts/platform');

function resolveBinaryPath(options = {}) {
  const env = options.env || process.env;
  const platform = options.platform || process.platform;

  if (env.KEYAPI_CLI_BINARY) {
    return env.KEYAPI_CLI_BINARY;
  }

  return path.join(__dirname, 'native', binaryNameForPlatform({ platform }));
}

function main() {
  const binaryPath = resolveBinaryPath();
  const result = spawnSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
  });

  if (result.error) {
    console.error(`Failed to start KeyAPI CLI: ${result.error.message}`);
    console.error('Try reinstalling the npm package, or set KEYAPI_CLI_BINARY.');
    process.exit(1);
  }

  if (result.signal) {
    process.kill(process.pid, result.signal);
    return;
  }

  process.exit(result.status ?? 1);
}

if (require.main === module) {
  main();
}

module.exports = {
  resolveBinaryPath,
};
