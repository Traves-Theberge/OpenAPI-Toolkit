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
  .action(async (spec: string, baseUrl: string, options: { export?: string; verbose?: boolean }) => {
    try {
      await runTests(spec, baseUrl, options);
      console.log('All tests passed.');
    } catch (error) {
      console.error('Tests failed:', error instanceof Error ? error.message : String(error));
      process.exit(1);
    }
  });

program.parse();