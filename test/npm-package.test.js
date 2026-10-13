const assert = require('node:assert/strict');
const path = require('node:path');
const test = require('node:test');

const { assetNameForPlatform } = require('../scripts/platform');
const { getInstallPlan } = require('../scripts/install');
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

test('getInstallPlan builds the release URL and binary destination', () => {
  const plan = getInstallPlan({
    rootDir: '/pkg',
    version: '1.2.3',
    platform: 'win32',
    arch: 'x64',
    env: {
      KEYAPI_CLI_REPOSITORY: 'Example/keyapi-cli',
    },
  });

  assert.equal(
    plan.url,
    'https://github.com/Example/keyapi-cli/releases/download/v1.2.3/keyapi_1.2.3_windows_amd64.exe',
  );
  assert.equal(plan.binaryPath, path.join('/pkg', 'npm', 'native', 'keyapi.exe'));
  assert.equal(plan.skip, false);
});
