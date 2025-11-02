#!/usr/bin/env node
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
const commander_1 = require("commander");
const validate_1 = require("./commands/validate");
const test_1 = require("./commands/test");
const config_1 = require("./config");
const chokidar = __importStar(require("chokidar"));
const path = __importStar(require("path"));
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
    .option('-r, --retry <count>', 'Retry failed requests with exponential backoff (default: 0)', '0')
    .option('-w, --watch', 'Watch spec file for changes and re-run tests')
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
        // Watch mode
        if (options.watch) {
            const specPath = path.resolve(spec);
            console.log(`\x1b[36m👁\x1b[0m  Watching ${specPath} for changes...\n`);
            console.log(`\x1b[90mPress Ctrl+C to stop\x1b[0m\n`);
            // Run tests initially
            try {
                await (0, test_1.runTests)(spec, baseUrl, mergedOptions);
                console.log('All tests passed.');
            }
            catch (error) {
                console.error('Tests failed:', error instanceof Error ? error.message : String(error));
                // Don't exit in watch mode, just show error
            }
            // Watch for changes
            const watcher = chokidar.watch(specPath, {
                persistent: true,
                ignoreInitial: true,
            });
            watcher.on('change', async () => {
                console.log(`\n\x1b[36m🔄 File changed, re-running tests...\x1b[0m\n`);
                try {
                    await (0, test_1.runTests)(spec, baseUrl, mergedOptions);
                    console.log('All tests passed.');
                }
                catch (error) {
                    console.error('Tests failed:', error instanceof Error ? error.message : String(error));
                    // Don't exit in watch mode, just show error
                }
            });
            watcher.on('error', (error) => {
                console.error('Watcher error:', error);
            });
            // Keep process alive
            process.on('SIGINT', () => {
                console.log('\n\x1b[36m👋 Stopping watch mode...\x1b[0m');
                watcher.close();
                process.exit(0);
            });
        }
        else {
            // Normal mode (run once)
            await (0, test_1.runTests)(spec, baseUrl, mergedOptions);
            console.log('All tests passed.');
        }
    }
    catch (error) {
        console.error('Tests failed:', error instanceof Error ? error.message : String(error));
        process.exit(1);
    }
});
program.parse();
//# sourceMappingURL=cli.js.map