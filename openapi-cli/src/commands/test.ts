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
  verbose?: boolean;
  timeout?: string;
  authBearer?: string;
  authApiKey?: string;
  authHeader?: string;
  authQuery?: string;
  authBasic?: string;
  header?: string[];
  methods?: string;
}

export async function runTests(specPath: string, baseUrl: string, options: TestOptions = {}): Promise<void> {
  // Load and parse spec
  const spec: OpenAPISpec = loadSpec(specPath);

  console.log(`\n🧪 Testing API: ${spec.info.title}`);
  console.log(`📍 Base URL: ${baseUrl}\n`);

  const results: TestResult[] = [];
  let successCount = 0;
  let failureCount = 0;

  // Parse timeout option
  const timeoutMs = options.timeout ? parseInt(options.timeout, 10) : 10000;

  // Parse methods filter
  const allowedMethods = options.methods
    ? options.methods.split(',').map(m => m.trim().toUpperCase())
    : null;

  // For each path and method, test the endpoint
  for (const [pathStr, methods] of Object.entries(spec.paths)) {
    for (const [method, operation] of Object.entries(methods)) {
      if (typeof operation === 'object' && operation !== null) {
        // Skip if method filter is active and this method is not in the list
        const methodUpper = method.toUpperCase();
        if (allowedMethods && !allowedMethods.includes(methodUpper)) {
          continue;
        }
        const result = await testEndpoint(baseUrl, pathStr, method.toUpperCase(), operation, options.verbose, timeoutMs, options);
        results.push(result);

        if (result.success) {
          successCount++;
          console.log(`\x1b[32m✓\x1b[0m ${result.method.padEnd(7)} ${result.endpoint.padEnd(40)} - ${result.status} ${result.message}`);
          if (options.verbose && result.duration) {
            console.log(`  \x1b[90mDuration: ${result.duration}ms\x1b[0m`);
            if (result.responseHeaders) {
              console.log(`  \x1b[90mResponse Headers: ${JSON.stringify(result.responseHeaders)}\x1b[0m`);
            }
          }
        } else {
          failureCount++;
          console.log(`\x1b[31m✗\x1b[0m ${result.method.padEnd(7)} ${result.endpoint.padEnd(40)} - ${result.message}`);
        }
      }
    }
  }

  // Summary
  console.log(`\n${'='.repeat(80)}`);
  console.log(`📊 Summary: ${successCount} passed, ${failureCount} failed, ${results.length} total`);

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
      console.log(`\x1b[32m✓ Results exported to ${options.export}\x1b[0m`);
    } catch (error) {
      console.error(`\x1b[31m✗ Failed to export results: ${error instanceof Error ? error.message : String(error)}\x1b[0m`);
    }
  }

  if (failureCount > 0) {
    console.log(`\x1b[31m✗ Some tests failed\x1b[0m\n`);
    process.exit(1);
  } else {
    console.log(`\x1b[32m✓ All tests passed!\x1b[0m\n`);
  }
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
      case 'POST':
        // Use requestBody schema if available, otherwise empty object
        const postBody = operation.requestBody?.content?.['application/json']?.example || {};
        response = await axios.post(url, postBody, config);
        break;
      case 'PUT':
        const putBody = operation.requestBody?.content?.['application/json']?.example || {};
        response = await axios.put(url, putBody, config);
        break;
      case 'PATCH':
        const patchBody = operation.requestBody?.content?.['application/json']?.example || {};
        response = await axios.patch(url, patchBody, config);
        break;
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