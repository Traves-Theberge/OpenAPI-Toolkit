import * as fs from 'fs';
import * as path from 'path';
import * as yaml from 'js-yaml'; // Note: need to add js-yaml to dependencies

export async function validateSpec(filePath: string): Promise<void> {
  if (!fs.existsSync(filePath)) {
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
      throw new Error('Unsupported file format. Use .json or .yaml');
    }
  } catch (error) {
    throw new Error(`Failed to parse spec file: ${(error as Error).message}`);
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