/**
 * This file is used for checking and updating the format of the JSON file.
 *
 * You can check the format via `node format.js --check` and regenerate the
 * file with the correct formatting using `node format.js --generate`.
 *
 * The formatting logic uses `JSON.stringify` with 2 spaces, which will keep
 * separating commas on the same line as any closing character. This technique
 * was chosen for simplicty and to align with common default JSON formatters,
 * such as VSCode.
 */

const fs = require("fs");
const path = require("path");

const jsonFilePath = path.join(__dirname, "crawler-user-agents.json");

const original = fs.readFileSync(jsonFilePath, "utf-8");

const updated = JSON.stringify(JSON.parse(original), null, 2) + '\n';

if (process.argv[2] === "--generate") {
    fs.writeFileSync(jsonFilePath, updated);
    process.exit(0);
}

if (process.argv[2] === "--check") {
    if (updated !== original) {
        console.error("JSON file format is wrong. Run `node format.js --generate` to update.");
        console.error("Format must be 2 spaces, with newlines for objects and arrays, and separating commas on the line with the previous closing character.");
        process.exit(1);
    }
}

