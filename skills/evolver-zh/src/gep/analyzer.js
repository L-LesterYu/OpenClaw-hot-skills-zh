const fs = require('fs');
const path = require('path');

// 创新：自我纠正分析器
// 分析过去的失败以建议更好的未来变异
// 模式：元学习

function analyzeFailures() {
  const memoryPath = path.join(process.cwd(), 'MEMORY.md');
  if (!fs.existsSync(memoryPath)) return { status: 'skipped', reason: '无记忆文件' };
  
  const content = fs.readFileSync(memoryPath, 'utf8');
  const failureRegex = /\|\s*\*\*F\d+\*\*\s*\|\s*Fix\s*\|\s*(.*?)\s*\|\s*\*\*(.*?)\*\*\s*\((.*?)\)\s*\|/g;
  
  const failures = [];
  let match;
  while ((match = failureRegex.exec(content)) !== null) {
    failures.push({
      summary: match[1].trim(),
      detail: match[2].trim()
    });
  }
  
  return {
    status: 'success',
    count: failures.length,
    failures: failures.slice(0, 3) // 返回前 3 条用于 prompt 上下文
  };
}

if (require.main === module) {
  console.log(JSON.stringify(analyzeFailures(), null, 2));
}

module.exports = { analyzeFailures };
