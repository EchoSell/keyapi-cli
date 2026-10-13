#!/usr/bin/env node
const fs = require('node:fs');
const https = require('node:https');
const path = require('node:path');

const { assetNameForPlatform, binaryNameForPlatform } = require('./platform');

const rootDir = path.resolve(__dirname, '..');
const packageJson = require(path.join(rootDir, 'package.json'));

function download(url, destination, redirects = 0) {
  return new Promise((resolve, reject) => {
    const request = https.get(
      url,
      {
        headers: {
          'User-Agent': '@echosell/keyapi-cli installer',
        },
      },
      (response) => {
        const location = response.headers.location;
        if (
          response.statusCode >= 300 &&
          response.statusCode < 400 &&
          location
        ) {
          response.resume();
          if (redirects >= 5) {
            reject(new Error('too many redirects while downloading KeyAPI CLI'));
            return;
          }
          download(new URL(location, url).toString(), destination, redirects + 1)
            .then(resolve, reject);
          return;
        }

        if (response.statusCode !== 200) {
          response.resume();
          reject(new Error(`download failed with HTTP ${response.statusCode}`));
          return;
        }

        const file = fs.createWriteStream(destination, { mode: 0o755 });
        response.pipe(file);
        file.on('finish', () => {
          file.close(resolve);
        });
        file.on('error', reject);
      },
    );

    request.on('error', reject);
  });
}

function getInstallPlan(options = {}) {
  const env = options.env || process.env;
  const baseDir = options.rootDir || rootDir;
  const version = options.version || packageJson.version;

  if (env.KEYAPI_CLI_BINARY) {
    return { skip: true };
  }

  const assetName = assetNameForPlatform({
    version,
    platform: options.platform,
    arch: options.arch,
  });
  const repo = env.KEYAPI_CLI_REPOSITORY || 'EchoSell/keyapi-cli';
  const tag = env.KEYAPI_CLI_RELEASE_TAG || `v${version}`;
  const url =
    env.KEYAPI_CLI_DOWNLOAD_URL ||
    `https://github.com/${repo}/releases/download/${tag}/${assetName}`;

  const nativeDir = path.join(baseDir, 'npm', 'native');
  const binaryPath = path.join(
    nativeDir,
    binaryNameForPlatform({ platform: options.platform }),
  );
  const tempPath = `${binaryPath}.download`;

  return {
    binaryPath,
    nativeDir,
    skip: false,
    tempPath,
    url,
  };
}

async function install(options = {}) {
  const plan = getInstallPlan(options);
  if (plan.skip) {
    return;
  }

  fs.mkdirSync(plan.nativeDir, { recursive: true });
  if (fs.existsSync(plan.tempPath)) {
    fs.rmSync(plan.tempPath, { force: true });
  }
  if (fs.existsSync(plan.binaryPath)) {
    fs.rmSync(plan.binaryPath, { force: true });
  }

  await download(plan.url, plan.tempPath);
  fs.renameSync(plan.tempPath, plan.binaryPath);
  fs.chmodSync(plan.binaryPath, 0o755);
}

if (require.main === module) {
  install().catch((error) => {
    console.error(`Failed to install KeyAPI CLI binary: ${error.message}`);
    console.error('Set KEYAPI_CLI_BINARY to an existing keyapi binary to skip download.');
    process.exit(1);
  });
}

module.exports = {
  getInstallPlan,
  install,
};
