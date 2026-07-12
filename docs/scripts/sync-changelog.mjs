#!/usr/bin/env node
/**
 * Sync CHANGELOG.md to docs/en/release-notes/changelog.md
 *
 * This script copies the content from the root CHANGELOG.md to the docs site,
 * with only formatting changes (title format).
 *
 * Run from the docs directory: node scripts/sync-changelog.mjs
 */

import { readFileSync, writeFileSync, mkdirSync } from "fs";
import { dirname, join } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const docsDir = join(__dirname, "..");
const rootDir = join(docsDir, "..");

const sourcePath = join(rootDir, "CHANGELOG.md");
const targetPathEn = join(docsDir, "en/release-notes/changelog.md");
const targetPathZh = join(docsDir, "zh/release-notes/changelog.md");

// Ensure target directories exist
mkdirSync(dirname(targetPathEn), { recursive: true });
mkdirSync(dirname(targetPathZh), { recursive: true });

const HEADER_EN = `# Changelog

This page documents the changes in each CompressKit release.

`;

const HEADER_ZH = `# 更新日志

此页面记录 CompressKit 每个版本的变更。

`;

// Read the source file
let content = readFileSync(sourcePath, "utf-8");

// Remove the HTML comment block at the top
content = content.replace(/<!--[\s\S]*?-->\n*/g, "");

// Remove the "# Changelog" title (we'll add our own header)
content = content.replace(/^# Changelog\n+/, "");

// Convert title format: ## [0.69] - 2025-12-29 -> ## 0.69 (2025-12-29)
// Also handle ## [Unreleased] without date
content = content.replace(
  /^## \[([^\]]+)\](?: - (\d{4}-\d{1,2}-\d{1,2}))?/gm,
  (match, version, date) => {
    if (date) {
      return `## ${version} (${date})`;
    }
    return `## ${version}`;
  }
);

// Remove subsection headers like ### Added, ### Changed, ### Fixed
content = content.replace(/^### (Added|Changed|Fixed|Improved|Tools|SDK|添加|变更|修复|改进)\n+/gm, "");

// Write the target file (English)
try {
  writeFileSync(targetPathEn, HEADER_EN + content.trim() + "\n");
  console.log(`Synced changelog to ${targetPathEn}`);
} catch (err) {
  console.warn(`Failed to write English changelog: ${err.message}`);
}

// Write the target file (Chinese - same content, just different header)
try {
  writeFileSync(targetPathZh, HEADER_ZH + content.trim() + "\n");
  console.log(`Synced changelog to ${targetPathZh}`);
} catch (err) {
  console.warn(`Failed to write Chinese changelog: ${err.message}`);
}

console.log("Changelog sync complete!");
