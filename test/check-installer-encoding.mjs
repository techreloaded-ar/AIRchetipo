#!/usr/bin/env node

import fs from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const repoRoot = path.resolve(__dirname, "..");

const installerPaths = [path.join(repoRoot, "install.ps1")];
const utf8Bom = [0xef, 0xbb, 0xbf];

let hasFailure = false;
for (const installerPath of installerPaths) {
  const bytes = await fs.readFile(installerPath);
  const hasBom = utf8Bom.every((value, index) => bytes[index] === value);
  if (!hasBom) {
    console.log(`OK  ${path.relative(repoRoot, installerPath)} has no UTF-8 BOM`);
    continue;
  }

  hasFailure = true;
  console.error(
    `FAIL ${path.relative(
      repoRoot,
      installerPath,
    )} starts with a UTF-8 BOM. 'irm ... | iex' breaks on Windows when the script starts with BOM bytes EF BB BF.`,
  );
}

process.exit(hasFailure ? 1 : 0);
