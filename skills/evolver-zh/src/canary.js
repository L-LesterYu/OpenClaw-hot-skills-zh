// 金丝雀脚本：在 forked 子进程中运行以验证 index.js 能正常加载
// 不会崩溃。退出码 0 = 安全，非零 = 异常。
//
// 这是 solidify 提交进化前的最后一道安全网。
// 如果补丁导致 index.js 损坏（语法错误、缺少 require 等），
// 金丝雀会在守护进程重启使用损坏代码之前捕获到问题。
try {
  require('../index.js');
  process.exit(0);
} catch (e) {
  process.stderr.write(String(e.message || e).slice(0, 500));
  process.exit(1);
}
