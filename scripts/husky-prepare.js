#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const { dirname, join } = require("node:path");

function resolveHuskyPackageJson() {
  try {
    return require.resolve("husky/package.json", {
      paths: [process.cwd(), __dirname],
    });
  } catch (_error) {
    return null;
  }
}

const huskyPackageJsonPath = resolveHuskyPackageJson();

if (!huskyPackageJsonPath) {
  console.log("[husky] package not installed, skipping prepare step.");
  process.exit(0);
}

const huskyCliPath = join(dirname(huskyPackageJsonPath), "bin.js");
const result = spawnSync(process.execPath, [huskyCliPath], {
  stdio: "inherit",
});

if (result.error) {
  console.error("[husky] failed to run CLI:", result.error);
  process.exit(result.status ?? 1);
}

process.exit(result.status ?? 0);

