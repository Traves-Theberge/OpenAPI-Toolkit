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
exports.validateSpec = validateSpec;
const fs = __importStar(require("fs"));
const path = __importStar(require("path"));
const yaml = __importStar(require("js-yaml"));
async function validateSpec(filePath) {
    console.log(`\n📄 Validating OpenAPI specification: ${filePath}`);
    if (!fs.existsSync(filePath)) {
        console.log(`\x1b[31m✗ File not found: ${filePath}\x1b[0m\n`);
        throw new Error(`File not found: ${filePath}`);
    }
    const ext = path.extname(filePath).toLowerCase();
    let spec;
    try {
        const content = fs.readFileSync(filePath, 'utf-8');
        if (ext === '.yaml' || ext === '.yml') {
            spec = yaml.load(content);
        }
        else if (ext === '.json') {
            spec = JSON.parse(content);
        }
        else {
            console.log(`\x1b[31m✗ Unsupported file format. Use .json or .yaml\x1b[0m\n`);
            throw new Error('Unsupported file format. Use .json or .yaml');
        }
    }
    catch (error) {
        console.log(`\x1b[31m✗ Failed to parse spec file: ${error.message}\x1b[0m\n`);
        throw new Error(`Failed to parse spec file: ${error.message}`);
    }
    // Collect all validation errors
    const errors = [];
    // Validate OpenAPI version
    if (!spec.openapi) {
        errors.push({
            path: 'openapi',
            message: 'Missing required field "openapi"',
        });
    }
    else if (typeof spec.openapi !== 'string' || !spec.openapi.startsWith('3.')) {
        errors.push({
            path: 'openapi',
            message: `Unsupported OpenAPI version: ${spec.openapi}. Only OpenAPI 3.x is supported`,
        });
    }
    // Validate info object
    if (!spec.info) {
        errors.push({
            path: 'info',
            message: 'Missing required "info" object',
        });
    }
    else {
        if (!spec.info.title) {
            errors.push({
                path: 'info.title',
                message: 'Missing required field "info.title"',
            });
        }
        if (!spec.info.version) {
            errors.push({
                path: 'info.version',
                message: 'Missing required field "info.version"',
            });
        }
    }
    // Validate paths
    if (!spec.paths) {
        errors.push({
            path: 'paths',
            message: 'Missing required "paths" object',
        });
    }
    else if (typeof spec.paths !== 'object') {
        errors.push({
            path: 'paths',
            message: 'Field "paths" must be an object',
        });
    }
    else {
        // Validate each path
        const pathCount = Object.keys(spec.paths).length;
        let operationCount = 0;
        for (const [pathName, pathItem] of Object.entries(spec.paths)) {
            if (!pathName.startsWith('/')) {
                errors.push({
                    path: `paths.${pathName}`,
                    message: 'Path must start with "/"',
                });
            }
            if (typeof pathItem !== 'object' || pathItem === null) {
                errors.push({
                    path: `paths.${pathName}`,
                    message: 'Path item must be an object',
                });
                continue;
            }
            // Validate operations
            const validMethods = ['get', 'post', 'put', 'delete', 'patch', 'head', 'options', 'trace'];
            for (const [method, operation] of Object.entries(pathItem)) {
                if (!validMethods.includes(method.toLowerCase())) {
                    continue; // Skip non-operation fields like parameters, summary, etc.
                }
                operationCount++;
                if (typeof operation !== 'object' || operation === null) {
                    errors.push({
                        path: `paths.${pathName}.${method}`,
                        message: 'Operation must be an object',
                    });
                    continue;
                }
                // Validate responses
                const op = operation;
                if (!op.responses) {
                    errors.push({
                        path: `paths.${pathName}.${method}.responses`,
                        message: 'Operation must have a "responses" object',
                    });
                }
            }
        }
        if (pathCount === 0) {
            errors.push({
                path: 'paths',
                message: 'Spec must define at least one path',
            });
        }
        console.log(`\x1b[36mℹ\x1b[0m Found ${pathCount} paths with ${operationCount} operations`);
    }
    // Display results
    if (errors.length > 0) {
        console.log(`\n\x1b[31m✗ Validation failed with ${errors.length} error(s):\x1b[0m\n`);
        errors.forEach((err, idx) => {
            console.log(`  ${idx + 1}. \x1b[33m${err.path}\x1b[0m: ${err.message}`);
        });
        console.log('');
        throw new Error('Validation failed');
    }
    else {
        console.log(`\x1b[32m✓ Validation successful!\x1b[0m`);
        console.log(`  OpenAPI Version: ${spec.openapi}`);
        console.log(`  Title: ${spec.info.title}`);
        console.log(`  Version: ${spec.info.version || 'N/A'}`);
        if (spec.info.description) {
            console.log(`  Description: ${spec.info.description}`);
        }
        console.log('');
    }
}
//# sourceMappingURL=validate.js.map