#!/usr/bin/env node

import { Command } from 'commander';
import { validateSpec } from './commands/validate';
import { runTests } from './commands/test';

const program = new Command();

program
  .name('openapi-test')
  .description('CLI tool for testing APIs against OpenAPI specifications')
  .version('1.0.0');

program
  .command('validate')
  .description('Validate an OpenAPI specification file')
  .argument('<file>', 'Path to the OpenAPI spec file')
  .action(async (file: string) => {
    try {
      await validateSpec(file);
      console.log('OpenAPI spec is valid.');
    } catch (error) {
      console.error('Validation failed:', error instanceof Error ? error.message : String(error));
      process.exit(1);
    }
  });

program
  .command('test')
  .description('Run API tests against an OpenAPI spec')
  .argument('<spec>', 'Path to the OpenAPI spec file')
  .argument('<baseUrl>', 'Base URL of the API to test')
  .option('-e, --export <file>', 'Export results to JSON file')
  .option('-v, --verbose', 'Show verbose output with request/response details')
  .option('-t, --timeout <ms>', 'Request timeout in milliseconds (default: 10000)', '10000')
  .option('--auth-bearer <token>', 'Bearer token authentication')
  .option('--auth-api-key <key>', 'API key authentication (use with --auth-header or --auth-query)')
  .option('--auth-header <name>', 'Header name for API key (default: X-API-Key)', 'X-API-Key')
  .option('--auth-query <name>', 'Query parameter name for API key')
  .option('--auth-basic <user:pass>', 'Basic authentication (username:password)')
  .option('-H, --header <header>', 'Custom header (Name: Value), repeatable', (value, previous: string[] = []) => [...previous, value], [])
  .option('-m, --methods <methods>', 'Filter by HTTP methods (comma-separated, e.g., GET,POST)')
  .option('-q, --quiet', 'Quiet mode - only show errors and final exit code')
  .option('-p, --paths <pattern>', 'Filter by path pattern (supports * wildcard, e.g., /users/*)')
  .action(async (spec: string, baseUrl: string, options: {
    export?: string;
    verbose?: boolean;
    timeout?: string;
    authBearer?: string;
    authApiKey?: string;
    authHeader?: string;
    authQuery?: string;
    authBasic?: string;
    header?: string[];
    methods?: string;
    quiet?: boolean;
    paths?: string;
  }) => {
    try {
      await runTests(spec, baseUrl, options);
      console.log('All tests passed.');
    } catch (error) {
      console.error('Tests failed:', error instanceof Error ? error.message : String(error));
      process.exit(1);
    }
  });

program.parse();