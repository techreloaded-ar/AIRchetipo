const fs = require('node:fs');

const path = require('node:path');

const repoRoot = path.resolve(__dirname, '..');

const brandCamel = 'B' + 'Mad';
const brandUpper = 'BM' + 'AD';
const brandLower = 'b' + 'mad';
const brandTitle = 'B' + 'mad';

const skipFiles = new Set([path.join(repoRoot, 'CHANGELOG.md'), path.join(repoRoot, `migrate-${brandLower}-to-air.sh`)]);

const skipDirs = new Set(['.git', 'node_modules', 'dist', 'v4-backup', '.husky', '.idea', '.next', 'coverage', 'build', 'tmp', 'logs']);

const binaryExtensions = new Set([
  '.png',
  '.jpg',
  '.jpeg',
  '.gif',
  '.ico',
  '.bmp',
  '.pdf',
  '.mp4',
  '.mp3',
  '.zip',
  '.gz',
  '.bz2',
  '.7z',
  '.tar',
  '.xz',
  '.jar',
  '.exe',
  '.dll',
  '.so',
  '.bin',
  '.woff',
  '.woff2',
  '.ttf',
  '.otf',
]);

const replacements = [
  { pattern: new RegExp(`${brandUpper}-METHOD™`, 'g'), value: 'AIRchetipo' },
  { pattern: new RegExp(`${brandUpper}-METHOD`, 'g'), value: 'AIRchetipo' },
  { pattern: new RegExp(`${brandCamel}-METHOD`, 'g'), value: 'AIRchetipo' },
  { pattern: new RegExp(`${brandUpper}-CORE™`, 'g'), value: 'AIRchetipo-CORE' },
  { pattern: new RegExp(`${brandUpper}-CORE`, 'g'), value: 'AIRchetipo-CORE' },
  { pattern: new RegExp(`${brandCamel}-CORE`, 'g'), value: 'AIRchetipo-CORE' },
  { pattern: new RegExp(`${brandLower}-`, 'g'), value: 'air-' },
  { pattern: new RegExp(`/` + brandLower + `/`, 'g'), value: '/air/' },
  { pattern: new RegExp(`${brandLower}:`, 'g'), value: 'air:' },
  { pattern: new RegExp(`${brandLower}_`, 'g'), value: 'air_' },
  { pattern: new RegExp(`${brandUpper}`, 'g'), value: 'AIRCHETIPO' },
  { pattern: new RegExp(`${brandCamel}`, 'g'), value: 'AIRchetipo' },
  { pattern: new RegExp(`${brandTitle}`, 'g'), value: 'AIRchetipo' },
  { pattern: new RegExp(`${brandLower}`, 'g'), value: 'air' },
];

(async () => {
  let filesChanged = 0;
  let replacementsMade = 0;

  async function walk(dir) {
    const entries = await fs.promises.readdir(dir, { withFileTypes: true });
    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        if (skipDirs.has(entry.name)) continue;
        await walk(fullPath);
      } else if (entry.isFile()) {
        if (skipFiles.has(fullPath)) continue;
        const ext = path.extname(entry.name).toLowerCase();
        if (binaryExtensions.has(ext)) continue;
        await processFile(fullPath);
      }
    }
  }

  async function processFile(filePath) {
    let original;
    try {
      original = await fs.promises.readFile(filePath, 'utf8');
    } catch {
      return;
    }

    let updated = original;
    let fileReplacements = 0;
    for (const { pattern, value } of replacements) {
      const before = updated;
      updated = updated.replace(pattern, (match) => {
        fileReplacements += 1;
        return value;
      });
      if (updated === before) continue;
    }

    if (updated !== original) {
      await fs.promises.writeFile(filePath, updated, 'utf8');
      filesChanged += 1;
      replacementsMade += fileReplacements;
      console.log(`Updated: ${path.relative(repoRoot, filePath)} (${fileReplacements} replacements)`);
    }
  }

  await walk(repoRoot);
  console.log(`\nCompleted. Files changed: ${filesChanged}. Total replacements: ${replacementsMade}.`);
})().catch((error) => {
  console.error(error);
  process.exit(1);
});
