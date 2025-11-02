import * as fs from 'fs';
import * as path from 'path';
import * as yaml from 'js-yaml';

export interface Config {
  // Authentication
  authBearer?: string;
  authApiKey?: string;
  authHeader?: string;
  authQuery?: string;
  authBasic?: string;

  // Headers
  headers?: string[];

  // Request options
  timeout?: number;
  verbose?: boolean;
  quiet?: boolean;

  // Filtering
  methods?: string;
  paths?: string;

  // Execution
  parallel?: number;

  // Export
  export?: string;
  exportHtml?: string;
  exportJunit?: string;
}

/**
 * Load configuration from a file (YAML or JSON)
 */
export function loadConfig(configPath: string): Config {
  try {
    const absolutePath = path.resolve(configPath);

    if (!fs.existsSync(absolutePath)) {
      throw new Error(`Config file not found: ${absolutePath}`);
    }

    const content = fs.readFileSync(absolutePath, 'utf8');
    const ext = path.extname(absolutePath).toLowerCase();

    let config: any;
    if (ext === '.yaml' || ext === '.yml') {
      config = yaml.load(content);
    } else if (ext === '.json') {
      config = JSON.parse(content);
    } else {
      throw new Error(`Unsupported config file format: ${ext}. Use .yaml, .yml, or .json`);
    }

    return normalizeConfig(config);
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to load config: ${error.message}`);
    }
    throw error;
  }
}

/**
 * Search for a config file in the current directory and parent directories
 */
export function findConfig(): string | null {
  const configFileNames = [
    '.openapi-cli.yaml',
    '.openapi-cli.yml',
    '.openapi-cli.json',
    'openapi-cli.yaml',
    'openapi-cli.yml',
    'openapi-cli.json',
  ];

  let currentDir = process.cwd();
  const rootDir = path.parse(currentDir).root;

  while (true) {
    for (const fileName of configFileNames) {
      const configPath = path.join(currentDir, fileName);
      if (fs.existsSync(configPath)) {
        return configPath;
      }
    }

    if (currentDir === rootDir) {
      break;
    }

    currentDir = path.dirname(currentDir);
  }

  return null;
}

/**
 * Normalize config property names to match CLI option names
 */
function normalizeConfig(config: any): Config {
  const normalized: Config = {};

  // Map config keys to their normalized versions
  const keyMap: Record<string, keyof Config> = {
    'auth-bearer': 'authBearer',
    'authBearer': 'authBearer',
    'auth_bearer': 'authBearer',

    'auth-api-key': 'authApiKey',
    'authApiKey': 'authApiKey',
    'auth_api_key': 'authApiKey',

    'auth-header': 'authHeader',
    'authHeader': 'authHeader',
    'auth_header': 'authHeader',

    'auth-query': 'authQuery',
    'authQuery': 'authQuery',
    'auth_query': 'authQuery',

    'auth-basic': 'authBasic',
    'authBasic': 'authBasic',
    'auth_basic': 'authBasic',

    'headers': 'headers',
    'header': 'headers',

    'timeout': 'timeout',
    'verbose': 'verbose',
    'quiet': 'quiet',
    'methods': 'methods',
    'paths': 'paths',
    'parallel': 'parallel',

    'export': 'export',
    'export-html': 'exportHtml',
    'exportHtml': 'exportHtml',
    'export_html': 'exportHtml',

    'export-junit': 'exportJunit',
    'exportJunit': 'exportJunit',
    'export_junit': 'exportJunit',
  };

  for (const [key, value] of Object.entries(config)) {
    const normalizedKey = keyMap[key];
    if (normalizedKey) {
      (normalized as any)[normalizedKey] = value;
    }
  }

  return normalized;
}

/**
 * Merge CLI options with config file options
 * CLI options take precedence over config file options
 */
export function mergeOptions(cliOptions: any, config: Config): any {
  const merged = { ...config, ...cliOptions };

  // Special handling for headers - merge arrays
  if (config.headers && cliOptions.header && cliOptions.header.length > 0) {
    merged.header = [...(config.headers || []), ...cliOptions.header];
  } else if (config.headers) {
    merged.header = config.headers;
  }

  return merged;
}
