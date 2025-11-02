import axios, { AxiosResponse } from 'axios';
import * as fs from 'fs';
import * as yaml from 'js-yaml';
import * as path from 'path';

interface OpenAPISpec {
  openapi: string;
  info: { title: string };
  paths: Record<string, Record<string, any>>;
}

interface TestResult {
  method: string;
  endpoint: string;
  status: number | null;
  success: boolean;
  message: string;
  duration?: number;
  timestamp?: string;
  requestHeaders?: Record<string, string>;
  responseHeaders?: Record<string, string>;
}

interface TestOptions {
  export?: string;
  exportHtml?: string;
  exportJunit?: string;
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
  parallel?: string;
}

export async function runTests(specPath: string, baseUrl: string, options: TestOptions = {}): Promise<void> {
  // Load and parse spec
  const spec: OpenAPISpec = loadSpec(specPath);

  if (!options.quiet) {
    console.log(`\n🧪 Testing API: ${spec.info.title}`);
    console.log(`📍 Base URL: ${baseUrl}\n`);
  }

  const results: TestResult[] = [];
  let successCount = 0;
  let failureCount = 0;

  // Parse timeout option
  const timeoutMs = options.timeout ? parseInt(options.timeout, 10) : 10000;

  // Parse methods filter
  const allowedMethods = options.methods
    ? options.methods.split(',').map(m => m.trim().toUpperCase())
    : null;

  // Parse path pattern filter (supports * wildcard)
  const pathPattern = options.paths ? options.paths.trim() : null;

  // Parse parallel option
  const parallelLimit = options.parallel ? parseInt(options.parallel, 10) : 1;
  const isParallel = parallelLimit > 1;

  // Collect all test tasks
  interface TestTask {
    pathStr: string;
    method: string;
    operation: any;
  }
  const testTasks: TestTask[] = [];

  for (const [pathStr, methods] of Object.entries(spec.paths)) {
    // Skip if path filter is active and this path doesn't match
    if (pathPattern && !matchesPattern(pathStr, pathPattern)) {
      continue;
    }

    for (const [method, operation] of Object.entries(methods)) {
      if (typeof operation === 'object' && operation !== null) {
        // Skip if method filter is active and this method is not in the list
        const methodUpper = method.toUpperCase();
        if (allowedMethods && !allowedMethods.includes(methodUpper)) {
          continue;
        }
        testTasks.push({ pathStr, method: methodUpper, operation });
      }
    }
  }

  // Execute tests (parallel or sequential)
  if (isParallel) {
    // Parallel execution with concurrency limit
    const testResults = await executeTestsInParallel(testTasks, baseUrl, timeoutMs, options, parallelLimit);
    results.push(...testResults);

    // Count successes and failures
    for (const result of testResults) {
      if (result.success) {
        successCount++;
        if (!options.quiet) {
          console.log(`\x1b[32m✓\x1b[0m ${result.method.padEnd(7)} ${result.endpoint.padEnd(40)} - ${result.status} ${result.message}`);
          if (options.verbose && result.duration) {
            console.log(`  \x1b[90mDuration: ${result.duration}ms\x1b[0m`);
            if (result.responseHeaders) {
              console.log(`  \x1b[90mResponse Headers: ${JSON.stringify(result.responseHeaders)}\x1b[0m`);
            }
          }
        }
      } else {
        failureCount++;
        // Always show errors even in quiet mode
        console.log(`\x1b[31m✗\x1b[0m ${result.method.padEnd(7)} ${result.endpoint.padEnd(40)} - ${result.message}`);
      }
    }
  } else {
    // Sequential execution (original behavior)
    for (const task of testTasks) {
      const result = await testEndpoint(baseUrl, task.pathStr, task.method, task.operation, options.verbose, timeoutMs, options);
      results.push(result);

      if (result.success) {
        successCount++;
        if (!options.quiet) {
          console.log(`\x1b[32m✓\x1b[0m ${result.method.padEnd(7)} ${result.endpoint.padEnd(40)} - ${result.status} ${result.message}`);
          if (options.verbose && result.duration) {
            console.log(`  \x1b[90mDuration: ${result.duration}ms\x1b[0m`);
            if (result.responseHeaders) {
              console.log(`  \x1b[90mResponse Headers: ${JSON.stringify(result.responseHeaders)}\x1b[0m`);
            }
          }
        }
      } else {
        failureCount++;
        // Always show errors even in quiet mode
        console.log(`\x1b[31m✗\x1b[0m ${result.method.padEnd(7)} ${result.endpoint.padEnd(40)} - ${result.message}`);
      }
    }
  }

  // Summary
  if (!options.quiet) {
    console.log(`\n${'='.repeat(80)}`);
    console.log(`📊 Summary: ${successCount} passed, ${failureCount} failed, ${results.length} total`);
  }

  // Export results if requested
  if (options.export) {
    try {
      const exportData = {
        timestamp: new Date().toISOString(),
        specPath,
        baseUrl,
        totalTests: results.length,
        passed: successCount,
        failed: failureCount,
        results: results.map(r => ({
          method: r.method,
          endpoint: r.endpoint,
          status: r.status,
          success: r.success,
          message: r.message,
          duration: r.duration,
          timestamp: r.timestamp,
          requestHeaders: r.requestHeaders,
          responseHeaders: r.responseHeaders,
        })),
      };

      fs.writeFileSync(options.export, JSON.stringify(exportData, null, 2), 'utf-8');
      if (!options.quiet) {
        console.log(`\x1b[32m✓ Results exported to ${options.export}\x1b[0m`);
      }
    } catch (error) {
      console.error(`\x1b[31m✗ Failed to export results: ${error instanceof Error ? error.message : String(error)}\x1b[0m`);
    }
  }

  // Export HTML if requested
  if (options.exportHtml) {
    try {
      const html = generateHtmlReport({
        timestamp: new Date().toISOString(),
        specPath,
        baseUrl,
        totalTests: results.length,
        passed: successCount,
        failed: failureCount,
        results,
      }, spec.info.title);

      fs.writeFileSync(options.exportHtml, html, 'utf-8');
      if (!options.quiet) {
        console.log(`\x1b[32m✓ HTML report exported to ${options.exportHtml}\x1b[0m`);
      }
    } catch (error) {
      console.error(`\x1b[31m✗ Failed to export HTML: ${error instanceof Error ? error.message : String(error)}\x1b[0m`);
    }
  }

  // Export JUnit XML if requested
  if (options.exportJunit) {
    try {
      const xml = generateJunitXml({
        timestamp: new Date().toISOString(),
        specPath,
        baseUrl,
        totalTests: results.length,
        passed: successCount,
        failed: failureCount,
        results,
      }, spec.info.title);

      fs.writeFileSync(options.exportJunit, xml, 'utf-8');
      if (!options.quiet) {
        console.log(`\x1b[32m✓ JUnit XML exported to ${options.exportJunit}\x1b[0m`);
      }
    } catch (error) {
      console.error(`\x1b[31m✗ Failed to export JUnit XML: ${error instanceof Error ? error.message : String(error)}\x1b[0m`);
    }
  }

  if (failureCount > 0) {
    if (!options.quiet) {
      console.log(`\x1b[31m✗ Some tests failed\x1b[0m\n`);
    }
    process.exit(1);
  } else {
    if (!options.quiet) {
      console.log(`\x1b[32m✓ All tests passed!\x1b[0m\n`);
    }
  }
}

/**
 * Generate JUnit XML from test results
 */
function generateJunitXml(data: {
  timestamp: string;
  specPath: string;
  baseUrl: string;
  totalTests: number;
  passed: number;
  failed: number;
  results: TestResult[];
}, apiTitle: string): string {
  // Calculate total time in seconds
  const totalTime = data.results.reduce((sum, r) => sum + (r.duration || 0), 0) / 1000;

  // Escape XML special characters
  const escapeXml = (str: string): string => {
    return str
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&apos;');
  };

  // Generate testcase elements
  const testCases = data.results.map(r => {
    const testName = `${r.method} ${r.endpoint}`;
    const className = `${apiTitle}.${r.method}`;
    const time = r.duration ? (r.duration / 1000).toFixed(3) : '0.000';

    if (r.success) {
      return `    <testcase name="${escapeXml(testName)}" classname="${escapeXml(className)}" time="${time}"/>`;
    } else {
      const failureMessage = escapeXml(r.message);
      const failureType = r.status ? `HTTP ${r.status}` : 'Error';
      return `    <testcase name="${escapeXml(testName)}" classname="${escapeXml(className)}" time="${time}">
      <failure message="${failureMessage}" type="${failureType}">
Test: ${escapeXml(testName)}
Status: ${r.status || 'N/A'}
Message: ${failureMessage}
Endpoint: ${escapeXml(data.baseUrl + r.endpoint)}
      </failure>
    </testcase>`;
    }
  }).join('\n');

  const hostname = new URL(data.baseUrl).hostname;

  return `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="${escapeXml(apiTitle)}" tests="${data.totalTests}" failures="${data.failed}" errors="0" time="${totalTime.toFixed(3)}">
  <testsuite name="${escapeXml(apiTitle)}" tests="${data.totalTests}" failures="${data.failed}" errors="0" skipped="0" time="${totalTime.toFixed(3)}" timestamp="${data.timestamp}" hostname="${escapeXml(hostname)}">
    <properties>
      <property name="spec" value="${escapeXml(data.specPath)}"/>
      <property name="baseUrl" value="${escapeXml(data.baseUrl)}"/>
      <property name="timestamp" value="${data.timestamp}"/>
    </properties>
${testCases}
    <system-out><![CDATA[
OpenAPI CLI Test Results
API: ${apiTitle}
Base URL: ${data.baseUrl}
Spec: ${data.specPath}
Total Tests: ${data.totalTests}
Passed: ${data.passed}
Failed: ${data.failed}
    ]]></system-out>
  </testsuite>
</testsuites>`;
}

/**
 * Generate HTML report from test results
 */
function generateHtmlReport(data: {
  timestamp: string;
  specPath: string;
  baseUrl: string;
  totalTests: number;
  passed: number;
  failed: number;
  results: TestResult[];
}, apiTitle: string): string {
  const passRate = data.totalTests > 0 ? ((data.passed / data.totalTests) * 100).toFixed(1) : '0';
  const statusClass = data.failed > 0 ? 'failed' : 'passed';

  const resultsHtml = data.results.map(r => {
    const rowClass = r.success ? 'success' : 'failure';
    const statusBadge = r.success ? '✓' : '✗';
    const statusBadgeClass = r.success ? 'badge-success' : 'badge-failure';

    return `
      <tr class="${rowClass}">
        <td><span class="badge ${statusBadgeClass}">${statusBadge}</span></td>
        <td><span class="method method-${r.method.toLowerCase()}">${r.method}</span></td>
        <td class="endpoint">${r.endpoint}</td>
        <td class="status">${r.status || 'N/A'}</td>
        <td>${r.message}</td>
        <td>${r.duration ? r.duration + 'ms' : 'N/A'}</td>
      </tr>`;
  }).join('');

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>API Test Results - ${apiTitle}</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
      background: #f5f7fa;
      color: #2c3e50;
      padding: 20px;
      line-height: 1.6;
    }
    .container {
      max-width: 1200px;
      margin: 0 auto;
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      overflow: hidden;
    }
    header {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      padding: 30px;
    }
    h1 {
      font-size: 28px;
      margin-bottom: 10px;
    }
    .meta {
      opacity: 0.9;
      font-size: 14px;
    }
    .meta-item {
      display: inline-block;
      margin-right: 20px;
      margin-top: 5px;
    }
    .summary {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 20px;
      padding: 30px;
      background: #f8f9fa;
      border-bottom: 1px solid #e9ecef;
    }
    .stat-card {
      background: white;
      padding: 20px;
      border-radius: 8px;
      border-left: 4px solid #667eea;
      box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    }
    .stat-card.passed {
      border-left-color: #28a745;
    }
    .stat-card.failed {
      border-left-color: #dc3545;
    }
    .stat-label {
      font-size: 12px;
      text-transform: uppercase;
      color: #6c757d;
      font-weight: 600;
      letter-spacing: 0.5px;
    }
    .stat-value {
      font-size: 32px;
      font-weight: bold;
      margin-top: 5px;
    }
    .stat-card.passed .stat-value {
      color: #28a745;
    }
    .stat-card.failed .stat-value {
      color: #dc3545;
    }
    .results-section {
      padding: 30px;
    }
    h2 {
      font-size: 20px;
      margin-bottom: 20px;
      color: #2c3e50;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 14px;
    }
    thead {
      background: #f8f9fa;
      position: sticky;
      top: 0;
      z-index: 10;
    }
    th {
      text-align: left;
      padding: 12px;
      font-weight: 600;
      color: #495057;
      border-bottom: 2px solid #dee2e6;
    }
    td {
      padding: 12px;
      border-bottom: 1px solid #f1f3f5;
    }
    tr.success {
      background: #f8fff9;
    }
    tr.failure {
      background: #fff5f5;
    }
    tr:hover {
      background: #e9ecef;
    }
    .badge {
      display: inline-block;
      width: 20px;
      height: 20px;
      line-height: 20px;
      text-align: center;
      border-radius: 50%;
      font-weight: bold;
      font-size: 12px;
    }
    .badge-success {
      background: #28a745;
      color: white;
    }
    .badge-failure {
      background: #dc3545;
      color: white;
    }
    .method {
      display: inline-block;
      padding: 4px 8px;
      border-radius: 4px;
      font-weight: 600;
      font-size: 12px;
      font-family: 'Monaco', 'Courier New', monospace;
    }
    .method-get { background: #e3f2fd; color: #1976d2; }
    .method-post { background: #e8f5e9; color: #388e3c; }
    .method-put { background: #fff3e0; color: #f57c00; }
    .method-patch { background: #fce4ec; color: #c2185b; }
    .method-delete { background: #ffebee; color: #d32f2f; }
    .method-head { background: #f3e5f5; color: #7b1fa2; }
    .method-options { background: #e0f2f1; color: #00796b; }
    .endpoint {
      font-family: 'Monaco', 'Courier New', monospace;
      font-size: 13px;
    }
    .status {
      font-weight: 600;
    }
    footer {
      padding: 20px 30px;
      background: #f8f9fa;
      text-align: center;
      font-size: 12px;
      color: #6c757d;
      border-top: 1px solid #e9ecef;
    }
    @media print {
      body {
        background: white;
        padding: 0;
      }
      .container {
        box-shadow: none;
      }
      tr:hover {
        background: inherit !important;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <h1>🧪 API Test Results</h1>
      <div class="meta">
        <span class="meta-item"><strong>API:</strong> ${apiTitle}</span>
        <span class="meta-item"><strong>Base URL:</strong> ${data.baseUrl}</span>
        <span class="meta-item"><strong>Spec:</strong> ${data.specPath}</span>
        <span class="meta-item"><strong>Timestamp:</strong> ${new Date(data.timestamp).toLocaleString()}</span>
      </div>
    </header>

    <div class="summary">
      <div class="stat-card">
        <div class="stat-label">Total Tests</div>
        <div class="stat-value">${data.totalTests}</div>
      </div>
      <div class="stat-card passed">
        <div class="stat-label">Passed</div>
        <div class="stat-value">${data.passed}</div>
      </div>
      <div class="stat-card failed">
        <div class="stat-label">Failed</div>
        <div class="stat-value">${data.failed}</div>
      </div>
      <div class="stat-card ${statusClass}">
        <div class="stat-label">Success Rate</div>
        <div class="stat-value">${passRate}%</div>
      </div>
    </div>

    <div class="results-section">
      <h2>Test Results</h2>
      <table>
        <thead>
          <tr>
            <th style="width: 50px;"></th>
            <th style="width: 80px;">Method</th>
            <th>Endpoint</th>
            <th style="width: 80px;">Status</th>
            <th>Message</th>
            <th style="width: 100px;">Duration</th>
          </tr>
        </thead>
        <tbody>
          ${resultsHtml}
        </tbody>
      </table>
    </div>

    <footer>
      Generated by <strong>OpenAPI CLI Tester</strong> | ${new Date().getFullYear()}
    </footer>
  </div>
</body>
</html>`;
}

/**
 * Execute tests in parallel with concurrency limit
 */
async function executeTestsInParallel(
  tasks: Array<{ pathStr: string; method: string; operation: any }>,
  baseUrl: string,
  timeout: number,
  options: TestOptions,
  concurrencyLimit: number
): Promise<TestResult[]> {
  const results: TestResult[] = [];
  const executing: Promise<void>[] = [];

  for (const task of tasks) {
    // Create test promise
    const testPromise = testEndpoint(
      baseUrl,
      task.pathStr,
      task.method,
      task.operation,
      options.verbose || false,
      timeout,
      options
    ).then(result => {
      results.push(result);
    });

    executing.push(testPromise);

    // If we've reached the concurrency limit, wait for one to finish
    if (executing.length >= concurrencyLimit) {
      await Promise.race(executing);
      // Remove completed promises
      executing.splice(0, executing.findIndex(p => p !== testPromise) + 1);
    }
  }

  // Wait for all remaining tests to complete
  await Promise.all(executing);

  return results;
}

function loadSpec(filePath: string): OpenAPISpec {
  if (!fs.existsSync(filePath)) {
    throw new Error(`File not found: ${filePath}`);
  }

  const ext = path.extname(filePath).toLowerCase();
  const content = fs.readFileSync(filePath, 'utf-8');

  if (ext === '.yaml' || ext === '.yml') {
    return yaml.load(content) as OpenAPISpec;
  } else if (ext === '.json') {
    return JSON.parse(content);
  } else {
    throw new Error('Unsupported file format. Use .json or .yaml');
  }
}

/**
 * Generate a request body from a JSON Schema definition
 */
function generateBodyFromSchema(schema: any): any {
  if (!schema) {
    return {};
  }

  // If there's an example, use it
  if (schema.example !== undefined) {
    return schema.example;
  }

  // Generate based on type
  switch (schema.type) {
    case 'object':
      const obj: any = {};
      const properties = schema.properties || {};
      const required = schema.required || [];

      // Generate required fields first
      for (const key of required) {
        if (properties[key]) {
          obj[key] = generateValueFromSchema(properties[key]);
        }
      }

      // Optionally add some non-required fields (50% chance for each)
      for (const key in properties) {
        if (!required.includes(key) && Math.random() > 0.5) {
          obj[key] = generateValueFromSchema(properties[key]);
        }
      }

      return obj;

    case 'array':
      const items = schema.items || {};
      // Generate array with 2 items by default
      return [
        generateValueFromSchema(items),
        generateValueFromSchema(items)
      ];

    default:
      return generateValueFromSchema(schema);
  }
}

/**
 * Generate a value from a JSON Schema property definition
 */
function generateValueFromSchema(schema: any): any {
  if (!schema) {
    return null;
  }

  // Use example if available
  if (schema.example !== undefined) {
    return schema.example;
  }

  // Use default if available
  if (schema.default !== undefined) {
    return schema.default;
  }

  // Use enum first value if available
  if (schema.enum && schema.enum.length > 0) {
    return schema.enum[0];
  }

  // Generate based on type
  switch (schema.type) {
    case 'string':
      // Handle format-specific strings
      if (schema.format === 'email') {
        return 'test@example.com';
      } else if (schema.format === 'uri' || schema.format === 'url') {
        return 'https://example.com';
      } else if (schema.format === 'date') {
        return '2025-01-01';
      } else if (schema.format === 'date-time') {
        return '2025-01-01T00:00:00Z';
      } else if (schema.format === 'uuid') {
        return '123e4567-e89b-12d3-a456-426614174000';
      }
      // Regular string with min/max length consideration
      const minLength = schema.minLength || 1;
      const maxLength = schema.maxLength || 10;
      const length = Math.min(maxLength, Math.max(minLength, 5));
      return 'test'.padEnd(length, 'x');

    case 'number':
    case 'integer':
      const min = schema.minimum !== undefined ? schema.minimum : 0;
      const max = schema.maximum !== undefined ? schema.maximum : 100;
      const value = Math.floor((min + max) / 2);
      return schema.type === 'integer' ? Math.floor(value) : value;

    case 'boolean':
      return true;

    case 'array':
      const items = schema.items || {};
      // Generate array with 2 items
      return [
        generateValueFromSchema(items),
        generateValueFromSchema(items)
      ];

    case 'object':
      const obj: any = {};
      const properties = schema.properties || {};
      const required = schema.required || [];

      // Generate all required fields
      for (const key of required) {
        if (properties[key]) {
          obj[key] = generateValueFromSchema(properties[key]);
        }
      }

      // Add some optional fields
      for (const key in properties) {
        if (!required.includes(key) && Math.random() > 0.5) {
          obj[key] = generateValueFromSchema(properties[key]);
        }
      }

      return obj;

    case 'null':
      return null;

    default:
      // If oneOf/anyOf/allOf exists, use the first schema
      if (schema.oneOf && schema.oneOf.length > 0) {
        return generateValueFromSchema(schema.oneOf[0]);
      }
      if (schema.anyOf && schema.anyOf.length > 0) {
        return generateValueFromSchema(schema.anyOf[0]);
      }
      if (schema.allOf && schema.allOf.length > 0) {
        // Merge all schemas in allOf
        let merged: any = {};
        for (const subSchema of schema.allOf) {
          const generated = generateValueFromSchema(subSchema);
          merged = { ...merged, ...generated };
        }
        return merged;
      }

      return 'test';
  }
}

/**
 * Replace path parameters like {id} with actual values
 */
function replacePlaceholders(pathStr: string): string {
  // Replace {id}, {userId}, etc. with "1"
  return pathStr.replace(/\{[^}]+\}/g, '1');
}

/**
 * Build query string from parameters
 */
function buildQueryParams(operation: any): string {
  if (!operation.parameters) {
    return '';
  }

  const queryParams = operation.parameters
    .filter((p: any) => p.in === 'query')
    .map((p: any) => {
      // Use example value if available, otherwise use a default based on type
      let value = '1';
      if (p.example !== undefined) {
        value = String(p.example);
      } else if (p.schema?.type === 'string') {
        value = 'test';
      } else if (p.schema?.type === 'boolean') {
        value = 'true';
      }
      return `${p.name}=${encodeURIComponent(value)}`;
    });

  return queryParams.length > 0 ? '?' + queryParams.join('&') : '';
}

async function testEndpoint(
  baseUrl: string,
  pathStr: string,
  method: string,
  operation: any,
  verbose: boolean = false,
  timeout: number = 10000,
  authOptions: TestOptions = {}
): Promise<TestResult> {
  // Replace path placeholders like {id} with actual values
  const processedPath = replacePlaceholders(pathStr);

  // Build query parameters
  let queryString = buildQueryParams(operation);

  // Add API key to query if specified
  if (authOptions.authApiKey && authOptions.authQuery) {
    const queryParam = `${authOptions.authQuery}=${encodeURIComponent(authOptions.authApiKey)}`;
    queryString = queryString ? `${queryString}&${queryParam}` : `?${queryParam}`;
  }

  const url = `${baseUrl}${processedPath}${queryString}`;

  try {
    let response: AxiosResponse;
    const startTime = Date.now();

    // Build request headers
    const headers: Record<string, string> = {};

    // Bearer token authentication
    if (authOptions.authBearer) {
      headers['Authorization'] = `Bearer ${authOptions.authBearer}`;
    }

    // API key in header
    if (authOptions.authApiKey && !authOptions.authQuery) {
      const headerName = authOptions.authHeader || 'X-API-Key';
      headers[headerName] = authOptions.authApiKey;
    }

    // Basic authentication
    if (authOptions.authBasic) {
      const encoded = Buffer.from(authOptions.authBasic).toString('base64');
      headers['Authorization'] = `Basic ${encoded}`;
    }

    // Custom headers
    if (authOptions.header && authOptions.header.length > 0) {
      for (const headerStr of authOptions.header) {
        const colonIndex = headerStr.indexOf(':');
        if (colonIndex > 0) {
          const name = headerStr.substring(0, colonIndex).trim();
          const value = headerStr.substring(colonIndex + 1).trim();
          headers[name] = value;
        }
      }
    }

    const config = {
      timeout: timeout, // Configurable timeout
      validateStatus: () => true, // Don't throw on any status code
      headers: headers,
    };

    switch (method) {
      case 'GET':
        response = await axios.get(url, config);
        break;
      case 'POST': {
        // Generate body from schema or use example
        const contentSchema = operation.requestBody?.content?.['application/json'];
        const postBody = contentSchema?.example ||
                        (contentSchema?.schema ? generateBodyFromSchema(contentSchema.schema) : {});
        response = await axios.post(url, postBody, config);
        break;
      }
      case 'PUT': {
        const contentSchema = operation.requestBody?.content?.['application/json'];
        const putBody = contentSchema?.example ||
                       (contentSchema?.schema ? generateBodyFromSchema(contentSchema.schema) : {});
        response = await axios.put(url, putBody, config);
        break;
      }
      case 'PATCH': {
        const contentSchema = operation.requestBody?.content?.['application/json'];
        const patchBody = contentSchema?.example ||
                         (contentSchema?.schema ? generateBodyFromSchema(contentSchema.schema) : {});
        response = await axios.patch(url, patchBody, config);
        break;
      }
      case 'DELETE':
        response = await axios.delete(url, config);
        break;
      case 'HEAD':
        response = await axios.head(url, config);
        break;
      case 'OPTIONS':
        response = await axios.options(url, config);
        break;
      default:
        return {
          method,
          endpoint: processedPath,
          status: null,
          success: false,
          message: `Unsupported HTTP method`,
        };
    }

    const duration = Date.now() - startTime;
    const success = response.status >= 200 && response.status < 300;

    const result: TestResult = {
      method,
      endpoint: processedPath,
      status: response.status,
      success,
      message: success ? 'OK' : `HTTP ${response.status} ${response.statusText}`,
      duration,
      timestamp: new Date().toISOString(),
    };

    // Add headers if verbose mode is enabled
    if (verbose) {
      result.requestHeaders = {
        'User-Agent': 'openapi-cli',
        'Accept': 'application/json',
      };
      result.responseHeaders = response.headers as Record<string, string>;
    }

    return result;
  } catch (error) {
    const err = error as any;
    let message = 'Unknown error';

    if (err.code === 'ECONNREFUSED') {
      message = 'Connection refused';
    } else if (err.code === 'ETIMEDOUT' || err.message?.includes('timeout')) {
      message = 'Request timeout';
    } else if (err.message) {
      message = err.message;
    }

    return {
      method,
      endpoint: processedPath,
      status: null,
      success: false,
      message,
    };
  }
}

/**
 * Match a path against a pattern with * wildcard support
 */
function matchesPattern(path: string, pattern: string): boolean {
  // Convert pattern to regex, escaping special chars except *
  const regexPattern = pattern
    .replace(/[.+?^${}()|[\]\\]/g, '\\$&') // Escape regex special chars
    .replace(/\*/g, '.*');                  // Convert * to .*

  const regex = new RegExp(`^${regexPattern}$`);
  return regex.test(path);
}