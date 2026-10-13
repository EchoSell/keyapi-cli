#!/usr/bin/env node
const fs = require('node:fs');
const path = require('node:path');

const targets = [
  { goos: 'windows', goarch: 'amd64', packageName: 'keyapi-cli-win32-x64', binaryName: 'keyapi.exe' },
  { goos: 'linux', goarch: 'amd64', packageName: 'keyapi-cli-linux-x64', binaryName: 'keyapi' },
  { goos: 'linux', goarch: 'arm64', packageName: 'keyapi-cli-linux-arm64', binaryName: 'keyapi' },
  { goos: 'darwin', goarch: 'amd64', packageName: 'keyapi-cli-darwin-x64', binaryName: 'keyapi' },
  { goos: 'darwin', goarch: 'arm64', packageName: 'keyapi-cli-darwin-arm64', binaryName: 'keyapi' },
];

function preparePlatformPackages(options = {}) {
  const rootDir = options.rootDir || path.resolve(__dirname, '..');
  const version = options.version || require(path.join(rootDir, 'package.json')).version;
  const distDir = options.distDir || path.join(rootDir, 'dist');
  const packagesDir = options.packagesDir || path.join(rootDir, 'npm-packages');

  for (const target of targets) {
    const source = path.join(
      distDir,
      `keyapi_${version}_${target.goos}_${target.goarch}${target.binaryName.endsWith('.exe') ? '.exe' : ''}`,
    );
    const binDir = path.join(packagesDir, target.packageName, 'bin');
    const destination = path.join(binDir, target.binaryName);

    if (!fs.existsSync(source)) {
      throw new Error(`missing built binary: ${source}`);
    }

    fs.mkdirSync(binDir, { recursive: true });
    fs.copyFileSync(source, destination);
    fs.chmodSync(destination, 0o755);
  }
}

if (require.main === module) {
  preparePlatformPackages();
}

module.exports = {
  preparePlatformPackages,
  targets,
};
