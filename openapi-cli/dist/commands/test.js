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
async function runTests(specPath, baseUrl) {
    // Load and parse spec
    const spec = loadSpec(specPath);
    // For each path and method, test the endpoint
    for (const [pathStr, methods] of Object.entries(spec.paths)) {
        for (const [method, operation] of Object.entries(methods)) {
            if (typeof operation === 'object' && operation !== null) {
                await testEndpoint(baseUrl, pathStr, method.toUpperCase());
            }
        }
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
async function testEndpoint(baseUrl, pathStr, method) {
    const url = `${baseUrl}${pathStr}`;
    try {
        let response;
        switch (method) {
            case 'GET':
                response = await axios_1.default.get(url);
                break;
            case 'POST':
                response = await axios_1.default.post(url, {});
                break;
            // Add more methods as needed
            default:
                console.log(`Skipping unsupported method: ${method} ${url}`);
                return;
        }
        if (response.status >= 200 && response.status < 300) {
            console.log(`✓ ${method} ${url} - ${response.status}`);
        }
        else {
            throw new Error(`Unexpected status: ${response.status}`);
        }
    }
    catch (error) {
        throw new Error(`Test failed for ${method} ${url}: ${error.message}`);
    }
}
//# sourceMappingURL=test.js.map