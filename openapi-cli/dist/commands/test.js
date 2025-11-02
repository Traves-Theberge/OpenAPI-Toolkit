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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.runTests = runTests;
const axios_1 = __importDefault(require("axios"));
const fs = __importStar(require("fs"));
const yaml = __importStar(require("js-yaml"));
const path = __importStar(require("path"));
async function runTests(specPath, baseUrl, options = {}) {
    // Load and parse spec
    const spec = loadSpec(specPath);
    console.log(`\n🧪 Testing API: ${spec.info.title}`);
    console.log(`📍 Base URL: ${baseUrl}\n`);
    const results = [];
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
                }
                else {
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
        }
        catch (error) {
            console.error(`\x1b[31m✗ Failed to export results: ${error instanceof Error ? error.message : String(error)}\x1b[0m`);
        }
    }
    if (failureCount > 0) {
        console.log(`\x1b[31m✗ Some tests failed\x1b[0m\n`);
        process.exit(1);
    }
    else {
        console.log(`\x1b[32m✓ All tests passed!\x1b[0m\n`);
    }
}
function loadSpec(filePath) {
    if (!fs.existsSync(filePath)) {
        throw new Error(`File not found: ${filePath}`);
    }
    const ext = path.extname(filePath).toLowerCase();
    const content = fs.readFileSync(filePath, 'utf-8');
    if (ext === '.yaml' || ext === '.yml') {
        return yaml.load(content);
    }
    else if (ext === '.json') {
        return JSON.parse(content);
    }
    else {
        throw new Error('Unsupported file format. Use .json or .yaml');
    }
}
/**
 * Replace path parameters like {id} with actual values
 */
function replacePlaceholders(pathStr) {
    // Replace {id}, {userId}, etc. with "1"
    return pathStr.replace(/\{[^}]+\}/g, '1');
}
/**
 * Build query string from parameters
 */
function buildQueryParams(operation) {
    if (!operation.parameters) {
        return '';
    }
    const queryParams = operation.parameters
        .filter((p) => p.in === 'query')
        .map((p) => {
        // Use example value if available, otherwise use a default based on type
        let value = '1';
        if (p.example !== undefined) {
            value = String(p.example);
        }
        else if (p.schema?.type === 'string') {
            value = 'test';
        }
        else if (p.schema?.type === 'boolean') {
            value = 'true';
        }
        return `${p.name}=${encodeURIComponent(value)}`;
    });
    return queryParams.length > 0 ? '?' + queryParams.join('&') : '';
}
async function testEndpoint(baseUrl, pathStr, method, operation, verbose = false, timeout = 10000, authOptions = {}) {
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
        let response;
        const startTime = Date.now();
        // Build request headers
        const headers = {};
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
                response = await axios_1.default.get(url, config);
                break;
            case 'POST':
                // Use requestBody schema if available, otherwise empty object
                const postBody = operation.requestBody?.content?.['application/json']?.example || {};
                response = await axios_1.default.post(url, postBody, config);
                break;
            case 'PUT':
                const putBody = operation.requestBody?.content?.['application/json']?.example || {};
                response = await axios_1.default.put(url, putBody, config);
                break;
            case 'PATCH':
                const patchBody = operation.requestBody?.content?.['application/json']?.example || {};
                response = await axios_1.default.patch(url, patchBody, config);
                break;
            case 'DELETE':
                response = await axios_1.default.delete(url, config);
                break;
            case 'HEAD':
                response = await axios_1.default.head(url, config);
                break;
            case 'OPTIONS':
                response = await axios_1.default.options(url, config);
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
        const result = {
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
            result.responseHeaders = response.headers;
        }
        return result;
    }
    catch (error) {
        const err = error;
        let message = 'Unknown error';
        if (err.code === 'ECONNREFUSED') {
            message = 'Connection refused';
        }
        else if (err.code === 'ETIMEDOUT' || err.message?.includes('timeout')) {
            message = 'Request timeout';
        }
        else if (err.message) {
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
//# sourceMappingURL=test.js.map