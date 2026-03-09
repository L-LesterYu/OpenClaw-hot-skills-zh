#!/usr/bin/env node

/**
 * Markdown 转换器 - 使用 markitdown 将文档转换为 Markdown
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// 检查 markitdown 是否已安装
function checkMarkitdown() {
  try {
    execSync('which markitdown', { stdio: 'ignore' });
    return true;
  } catch (error) {
    console.error('错误：markitdown 未安装。');
    console.error('请使用以下命令安装：pip install markitdown');
    process.exit(1);
  }
}

// 显示帮助信息
function showHelp() {
  console.log(`
Markdown 转换器 - 将文档转换为 Markdown

用法：
  ./convert.js <文件或链接>
  ./convert.js --help

支持的格式：
  - 文档：PDF, DOCX, RTF
  - 电子表格：XLSX, XLS, CSV
  - 演示文稿：PPTX
  - 网页：HTML, URL链接
  - 数据：JSON, XML
  - 图片：PNG, JPG（带 EXIF/OCR）
  - 音频：MP3, WAV（带转录）
  - 压缩包：ZIP
  - 电子书：EPUB
  - YouTube 链接

示例：
  ./convert.js document.pdf
  ./convert.js spreadsheet.xlsx
  ./convert.js https://example.com
  ./convert.js https://youtube.com/watch?v=VIDEO_ID
`);
}

// 主函数
function main() {
  const args = process.argv.slice(2);

  if (args.length === 0 || args[0] === '--help' || args[0] === '-h') {
    showHelp();
    process.exit(0);
  }

  checkMarkitdown();

  const input = args[0];

  try {
    // 使用 markitdown 转换文件
    const command = `markitdown "${input}"`;
    const output = execSync(command, { encoding: 'utf-8', maxBuffer: 50 * 1024 * 1024 });
    console.log(output);
  } catch (error) {
    console.error('转换文件时出错：', error.message);
    if (error.stderr) {
      console.error(error.stderr.toString());
    }
    process.exit(1);
  }
}

main();
