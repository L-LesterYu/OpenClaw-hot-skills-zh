// 稳定的设备标识符，用于节点身份识别。
// 生成基于硬件的指纹，在目录变更、重启和 evolver 升级后持久存在。
// 由 getNodeId() 和 env_fingerprint 使用。
//
// 优先级链：
//   1. EVOMAP_DEVICE_ID 环境变量       （显式覆盖，推荐容器使用）
//   2. ~/.evomap/device_id 文件        （从上次运行持久化）
//   3. <project>/.evomap_device_id     （无 $HOME 的容器回退路径）
//   4. /etc/machine-id                 （Linux，操作系统安装时设置）
//   5. IOPlatformUUID                  （macOS 硬件 UUID）
//   6. Docker/OCI 容器 ID              （来自 /proc/self/cgroup 或 /proc/self/mountinfo）
//   7. hostname + MAC 地址             （基于网络的回退）
//   8. 随机 128 位十六进制              （最后手段，立即持久化）

const os = require('os');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

const DEVICE_ID_DIR = path.join(os.homedir(), '.evomap');
const DEVICE_ID_FILE = path.join(DEVICE_ID_DIR, 'device_id');
const LOCAL_DEVICE_ID_FILE = path.resolve(__dirname, '..', '..', '.evomap_device_id');

let _cachedDeviceId = null;

const DEVICE_ID_RE = /^[a-f0-9]{16,64}$/;

function isContainer() {
  try {
    if (fs.existsSync('/.dockerenv')) return true;
  } catch {}
  try {
    const cgroup = fs.readFileSync('/proc/1/cgroup', 'utf8');
    if (/docker|kubepods|containerd|cri-o|lxc|ecs/i.test(cgroup)) return true;
  } catch {}
  try {
    if (fs.existsSync('/run/.containerenv')) return true;
  } catch {}
  return false;
}

function readMachineId() {
  try {
    const mid = fs.readFileSync('/etc/machine-id', 'utf8').trim();
    if (mid && mid.length >= 16) return mid;
  } catch {}

  if (process.platform === 'darwin') {
    try {
      const { execFileSync } = require('child_process');
      const raw = execFileSync('ioreg', ['-rd1', '-c', 'IOPlatformExpertDevice'], {
        encoding: 'utf8',
        timeout: 3000,
        stdio: ['ignore', 'pipe', 'ignore'],
      });
      const match = raw.match(/"IOPlatformUUID"\s*=\s*"([^"]+)"/);
      if (match && match[1]) return match[1];
    } catch {}
  }

  return null;
}

// 从 cgroup 或 mountinfo 中提取 Docker/OCI 容器 ID。
// 容器 ID 是 64 字符十六进制，在容器生命周期内稳定。
// 在非容器主机或解析失败时返回 null。
function readContainerId() {
  // 方法 1：/proc/self/cgroup（适用于 cgroup v1 和大多数 Docker 设置）
  try {
    const cgroup = fs.readFileSync('/proc/self/cgroup', 'utf8');
    const match = cgroup.match(/[a-f0-9]{64}/);
    if (match) return match[0];
  } catch {}

  // 方法 2：/proc/self/mountinfo（适用于 cgroup v2 / containerd）
  try {
    const mountinfo = fs.readFileSync('/proc/self/mountinfo', 'utf8');
    const match = mountinfo.match(/[a-f0-9]{64}/);
    if (match) return match[0];
  } catch {}

  // 方法 3：Docker 中 hostname 默认为短容器 ID（12 个十六进制字符）
  if (isContainer()) {
    const hostname = os.hostname();
    if (/^[a-f0-9]{12,64}$/.test(hostname)) return hostname;
  }

  return null;
}

function getMacAddresses() {
  const ifaces = os.networkInterfaces();
  const macs = [];
  for (const name of Object.keys(ifaces)) {
    for (const iface of ifaces[name]) {
      if (!iface.internal && iface.mac && iface.mac !== '00:00:00:00:00:00') {
        macs.push(iface.mac);
      }
    }
  }
  macs.sort();
  return macs;
}

function generateDeviceId() {
  const machineId = readMachineId();
  if (machineId) {
    return crypto.createHash('sha256').update('evomap:' + machineId).digest('hex').slice(0, 32);
  }

  // 容器 ID：在容器生命周期内稳定，但重新创建后会变化。
  // 仍优于随机值，可在单次部署中保持身份。
  const containerId = readContainerId();
  if (containerId) {
    return crypto.createHash('sha256').update('evomap:container:' + containerId).digest('hex').slice(0, 32);
  }

  const macs = getMacAddresses();
  if (macs.length > 0) {
    const raw = os.hostname() + '|' + macs.join(',');
    return crypto.createHash('sha256').update('evomap:' + raw).digest('hex').slice(0, 32);
  }

  return crypto.randomBytes(16).toString('hex');
}

function persistDeviceId(id) {
  // Try primary path (~/.evomap/device_id)
  try {
    if (!fs.existsSync(DEVICE_ID_DIR)) {
      fs.mkdirSync(DEVICE_ID_DIR, { recursive: true, mode: 0o700 });
    }
    fs.writeFileSync(DEVICE_ID_FILE, id, { encoding: 'utf8', mode: 0o600 });
    return;
  } catch {}

  // 回退：项目本地文件（适用于 $HOME 是临时的
  // 但项目目录作为卷挂载的容器）
  try {
    fs.writeFileSync(LOCAL_DEVICE_ID_FILE, id, { encoding: 'utf8', mode: 0o600 });
    return;
  } catch {}

  console.error(
    '[evolver] WARN: failed to persist device_id to ' + DEVICE_ID_FILE +
    ' or ' + LOCAL_DEVICE_ID_FILE +
    ' -- 节点身份可能在重启后变化。' +
    ' 在容器中设置 EVOMAP_DEVICE_ID 环境变量以获得稳定身份。'
  );
}

function loadPersistedDeviceId() {
  // Try primary path
  try {
    if (fs.existsSync(DEVICE_ID_FILE)) {
      const id = fs.readFileSync(DEVICE_ID_FILE, 'utf8').trim();
      if (id && DEVICE_ID_RE.test(id)) return id;
    }
  } catch {}

  // Try project-local fallback
  try {
    if (fs.existsSync(LOCAL_DEVICE_ID_FILE)) {
      const id = fs.readFileSync(LOCAL_DEVICE_ID_FILE, 'utf8').trim();
      if (id && DEVICE_ID_RE.test(id)) return id;
    }
  } catch {}

  return null;
}

function getDeviceId() {
  if (_cachedDeviceId) return _cachedDeviceId;

  // 1. Env var override (validated)
  if (process.env.EVOMAP_DEVICE_ID) {
    const envId = String(process.env.EVOMAP_DEVICE_ID).trim().toLowerCase();
    if (DEVICE_ID_RE.test(envId)) {
      _cachedDeviceId = envId;
      return _cachedDeviceId;
    }
  }

  // 2. Previously persisted (checks both ~/.evomap/ and project-local)
  const persisted = loadPersistedDeviceId();
  if (persisted) {
    _cachedDeviceId = persisted;
    return _cachedDeviceId;
  }

  // 3. Generate from hardware / container metadata and persist
  const inContainer = isContainer();
  const generated = generateDeviceId();
  persistDeviceId(generated);
  _cachedDeviceId = generated;

  if (inContainer && !process.env.EVOMAP_DEVICE_ID) {
    console.error(
      '[evolver] 注意：在容器中运行但未设置 EVOMAP_DEVICE_ID。' +
      ' 已自动生成并持久化 device_id，但为确保' +
      ' 跨重启的稳定性，请设置 EVOMAP_DEVICE_ID 环境变量' +
      ' 或在 ~/.evomap/ 挂载持久卷'
    );
  }

  return _cachedDeviceId;
}

module.exports = { getDeviceId, isContainer };
