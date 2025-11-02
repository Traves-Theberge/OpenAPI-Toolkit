#!/usr/bin/env node
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const commander_1 = require("commander");
const validate_1 = require("./commands/validate");
const test_1 = require("./commands/test");
const program = new commander_1.Command();
program
    .name('openapi-test')
    .description('CLI tool for testing APIs against OpenAPI specifications')
    .version('1.0.0');
program
    .command('validate')
    .description('Validate an OpenAPI specification file')
    .argument('<file>', 'Path to the OpenAPI spec file')
    .action(async (file) => {
    try {
        await (0, validate_1.validateSpec)(file);
        console.log('OpenAPI spec is valid.');
    }
    catch (error) {
        console.error('Validation failed:', error instanceof Error ? error.message : String(error));
        process.exit(1);
    }
});
program
    .command('test')
    .description('Run API tests against an OpenAPI spec')
    .argument('<spec>', 'Path to the OpenAPI spec file')
    .argument('<baseUrl>', 'Base URL of the API to test')
    .action(async (spec, baseUrl) => {
    try {
        await (0, test_1.runTests)(spec, baseUrl);
        console.log('All tests passed.');
    }
    catch (error) {
        console.error('Tests failed:', error instanceof Error ? error.message : String(error));
        process.exit(1);
    }
});
program.parse();
//# sourceMappingURL=cli.js.map