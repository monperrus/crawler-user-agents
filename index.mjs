import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const crawlers = require('./crawler-user-agents.json');

export default crawlers;
