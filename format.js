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
        console.error("JSON file format is wrong. Run `node validate.js --generate` to update.");
        process.exit(1);
    }
}

