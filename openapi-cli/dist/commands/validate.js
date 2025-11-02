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
const yaml = __importStar(require("js-yaml")); // Note: need to add js-yaml to dependencies
async function validateSpec(filePath) {
    if (!fs.existsSync(filePath)) {
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
            throw new Error('Unsupported file format. Use .json or .yaml');
        }
    }
    catch (error) {
        throw new Error(`Failed to parse spec file: ${error.message}`);
    }
    // Basic validation
    if (!spec.openapi || !spec.openapi.startsWith('3.')) {
        throw new Error('Invalid OpenAPI spec: missing or invalid openapi version');
    }
    if (!spec.info || !spec.info.title) {
        throw new Error('Invalid OpenAPI spec: missing info.title');
    }
    if (!spec.paths || typeof spec.paths !== 'object') {
        throw new Error('Invalid OpenAPI spec: missing or invalid paths');
    }
    // Could add more validations here
}
//# sourceMappingURL=validate.js.map