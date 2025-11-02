import * as fs from 'fs';
import * as path from 'path';
import * as yaml from 'js-yaml';

interface ValidationError {
  path: string;
  message: string;
}

export async function validateSpec(filePath: string): Promise<void> {
  console.log(`\n📄 Validating OpenAPI specification: ${filePath}`);

  if (!fs.existsSync(filePath)) {
    console.log(`\x1b[31m✗ File not found: ${filePath}\x1b[0m\n`);
    throw new Error(`File not found: ${filePath}`);
  }

  const ext = path.extname(filePath).toLowerCase();
  let spec: any;

  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    if (ext === '.yaml' || ext === '.yml') {
      spec = yaml.load(content);
    } else if (ext === '.json') {
      spec = JSON.parse(content);
    } else {
      console.log(`\x1b[31m✗ Unsupported file format. Use .json or .yaml\x1b[0m\n`);
      throw new Error('Unsupported file format. Use .json or .yaml');
    }
  } catch (error) {
    console.log(`\x1b[31m✗ Failed to parse spec file: ${(error as Error).message}\x1b[0m\n`);
    throw new Error(`Failed to parse spec file: ${(error as Error).message}`);
  }

  // Collect all validation errors
  const errors: ValidationError[] = [];

  // Validate OpenAPI version
  if (!spec.openapi) {
    errors.push({
      path: 'openapi',
      message: 'Missing required field "openapi"',
    });
  } else if (!spec.openapi.startsWith('3.')) {
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
  } else {
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
  } else if (typeof spec.paths !== 'object') {
    errors.push({
      path: 'paths',
      message: 'Field "paths" must be an object',
    });
  } else {
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
      for (const [method, operation] of Object.entries(pathItem as any)) {
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
        const op = operation as any;
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
  } else {
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