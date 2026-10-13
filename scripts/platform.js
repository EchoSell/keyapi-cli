const osByPlatform = {
  win32: 'windows',
  darwin: 'darwin',
  linux: 'linux',
};

const archByArch = {
  x64: 'amd64',
  arm64: 'arm64',
};

function targetForPlatform(options = {}) {
  const platform = options.platform || process.platform;
  const arch = options.arch || process.arch;
  const os = osByPlatform[platform];
  const goarch = archByArch[arch];

  if (!os || !goarch) {
    throw new Error(`Unsupported platform: ${platform}/${arch}`);
  }

  return { os, arch: goarch };
}

function assetNameForPlatform(options = {}) {
  if (!options.version) {
    throw new Error('version is required');
  }

  const target = targetForPlatform(options);
  const suffix = target.os === 'windows' ? '.exe' : '';
  return `keyapi_${options.version}_${target.os}_${target.arch}${suffix}`;
}

function binaryNameForPlatform(options = {}) {
  const platform = options.platform || process.platform;
  return platform === 'win32' ? 'keyapi.exe' : 'keyapi';
}

module.exports = {
  assetNameForPlatform,
  binaryNameForPlatform,
  targetForPlatform,
};
