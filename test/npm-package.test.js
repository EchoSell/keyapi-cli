const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const test = require('node:test');

const {
  assetNameForPlatform,
  binaryNameForPlatform,
  packageNameForPlatform,
} = require('../scripts/platform');
const { preparePlatformPackages, targets } = require('../scripts/prepare-platform-packages');
const { resolveBinaryPath } = require('../npm/keyapi');

test('assetNameForPlatform returns release asset names for supported systems', () => {
  assert.equal(
    assetNameForPlatform({ version: '1.2.3', platform: 'win32', arch: 'x64' }),
    'keyapi_1.2.3_windows_amd64.exe',
  );
  assert.equal(
    assetNameForPlatform({ version: '1.2.3', platform: 'darwin', arch: 'arm64' }),
    'keyapi_1.2.3_darwin_arm64',
  );
  assert.equal(
    assetNameForPlatform({ version: '1.2.3', platform: 'linux', arch: 'x64' }),
    'keyapi_1.2.3_linux_amd64',
  );
});

test('assetNameForPlatform rejects unsupported platforms clearly', () => {
  assert.throws(
    () => assetNameForPlatform({ version: '1.2.3', platform: 'freebsd', arch: 'x64' }),
    /Unsupported platform: freebsd\/x64/,
  );
});

test('resolveBinaryPath honors KEYAPI_CLI_BINARY override', () => {
  const override = path.join('custom', process.platform === 'win32' ? 'keyapi.exe' : 'keyapi');
  assert.equal(
    resolveBinaryPath({ env: { KEYAPI_CLI_BINARY: override } }),
    override,
  );
});

test('packageNameForPlatform returns npm platform package names', () => {
  assert.equal(packageNameForPlatform({ platform: 'win32', arch: 'x64' }), 'keyapi-cli-win32-x64');
  assert.equal(packageNameForPlatform({ platform: 'darwin', arch: 'arm64' }), 'keyapi-cli-darwin-arm64');
  assert.equal(packageNameForPlatform({ platform: 'linux', arch: 'x64' }), 'keyapi-cli-linux-x64');
});

test('resolveBinaryPath resolves the binary from the platform package', () => {
  const calls = [];
  const resolved = path.join('/npm', 'node_modules', 'keyapi-cli-win32-x64', 'bin', 'keyapi.exe');

  assert.equal(
    resolveBinaryPath({
      env: {},
      platform: 'win32',
      arch: 'x64',
      resolve: (request) => {
        calls.push(request);
        return resolved;
      },
    }),
    resolved,
  );
  assert.deepEqual(calls, [path.join('keyapi-cli-win32-x64', 'bin', binaryNameForPlatform({ platform: 'win32' }))]);
});

test('preparePlatformPackages copies built binaries into platform packages', () => {
  const rootDir = fs.mkdtempSync(path.join(os.tmpdir(), 'keyapi-packages-'));
  const distDir = path.join(rootDir, 'dist');
  const packagesDir = path.join(rootDir, 'npm-packages');
  const version = '9.8.7';

  fs.mkdirSync(distDir, { recursive: true });
  for (const target of targets) {
    const sourceName = `keyapi_${version}_${target.goos}_${target.goarch}${target.binaryName.endsWith('.exe') ? '.exe' : ''}`;
    fs.writeFileSync(path.join(distDir, sourceName), `${target.packageName}\n`);
  }

  preparePlatformPackages({ rootDir, version, distDir, packagesDir });

  for (const target of targets) {
    const output = path.join(packagesDir, target.packageName, 'bin', target.binaryName);
    assert.equal(fs.readFileSync(output, 'utf8'), `${target.packageName}\n`);
  }
});
