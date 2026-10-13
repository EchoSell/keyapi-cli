#!/usr/bin/env node
const { spawnSync } = require('node:child_process');
const path = require('node:path');

const {
  binaryNameForPlatform,
  packageNameForPlatform,
} = require('../scripts/platform');

function resolveBinaryPath(options = {}) {
  const env = options.env || process.env;
  const platform = options.platform || process.platform;
  const arch = options.arch || process.arch;
  const resolve = options.resolve || require.resolve;

  if (env.KEYAPI_CLI_BINARY) {
    return env.KEYAPI_CLI_BINARY;
  }

  const packageName = packageNameForPlatform({ platform, arch });
  return resolve(path.join(packageName, 'bin', binaryNameForPlatform({ platform })));
}

function main() {
  const binaryPath = resolveBinaryPath();
  const result = spawnSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
  });

  if (result.error) {
    console.error(`Failed to start KeyAPI CLI: ${result.error.message}`);
    console.error('Try reinstalling keyapi-cli, or set KEYAPI_CLI_BINARY.');
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
