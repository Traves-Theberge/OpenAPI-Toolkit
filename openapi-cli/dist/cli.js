#!/usr/bin/env node
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const commander_1 = require("commander");
const validate_1 = require("./commands/validate");
const test_1 = require("./commands/test");
const config_1 = require("./config");
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
    .option('-c, --config <file>', 'Path to config file (YAML or JSON)')
    .option('-e, --export <file>', 'Export results to JSON file')
    .option('--export-html <file>', 'Export results to HTML file')
    .option('--export-junit <file>', 'Export results to JUnit XML file')
    .option('-v, --verbose', 'Show verbose output with request/response details')
    .option('-t, --timeout <ms>', 'Request timeout in milliseconds (default: 10000)', '10000')
    .option('--auth-bearer <token>', 'Bearer token authentication')
    .option('--auth-api-key <key>', 'API key authentication (use with --auth-header or --auth-query)')
    .option('--auth-header <name>', 'Header name for API key (default: X-API-Key)', 'X-API-Key')
    .option('--auth-query <name>', 'Query parameter name for API key')
    .option('--auth-basic <user:pass>', 'Basic authentication (username:password)')
    .option('-H, --header <header>', 'Custom header (Name: Value), repeatable', (value, previous = []) => [...previous, value], [])
    .option('-m, --methods <methods>', 'Filter by HTTP methods (comma-separated, e.g., GET,POST)')
    .option('-q, --quiet', 'Quiet mode - only show errors and final exit code')
    .option('-p, --paths <pattern>', 'Filter by path pattern (supports * wildcard, e.g., /users/*)')
    .option('--parallel <limit>', 'Run tests in parallel with concurrency limit (default: 5)', '5')
    .option('--validate-schema', 'Validate response bodies against OpenAPI schemas')
    .action(async (spec, baseUrl, options) => {
    try {
        // Load config file if specified or search for one
        let config = {};
        if (options.config) {
            config = (0, config_1.loadConfig)(options.config);
        }
        else {
            const configPath = (0, config_1.findConfig)();
            if (configPath) {
                config = (0, config_1.loadConfig)(configPath);
                if (!options.quiet) {
                    console.log(`Using config file: ${configPath}`);
                }
            }
        }
        // Merge config with CLI options (CLI options take precedence)
        const mergedOptions = (0, config_1.mergeOptions)(options, config);
        await (0, test_1.runTests)(spec, baseUrl, mergedOptions);
        console.log('All tests passed.');
    }
    catch (error) {
        console.error('Tests failed:', error instanceof Error ? error.message : String(error));
        process.exit(1);
    }
});
program.parse();
//# sourceMappingURL=cli.js.map