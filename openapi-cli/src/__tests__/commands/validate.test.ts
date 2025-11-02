import { validateSpec } from '../../commands/validate';
import * as fs from 'fs';

jest.mock('fs');

describe('validateSpec', () => {
  const mockFs = fs as jest.Mocked<typeof fs>;

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should validate a correct OpenAPI spec', async () => {
    const specContent = `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK
`;
    mockFs.existsSync.mockReturnValue(true);
    mockFs.readFileSync.mockReturnValue(specContent);

    await expect(validateSpec('test.yaml')).resolves.toBeUndefined();
  });

  it('should throw error for missing file', async () => {
    mockFs.existsSync.mockReturnValue(false);

    await expect(validateSpec('missing.yaml')).rejects.toThrow('File not found');
  });

  it('should throw error for invalid spec', async () => {
    const invalidSpec = `
info:
  title: Test
`;
    mockFs.existsSync.mockReturnValue(true);
    mockFs.readFileSync.mockReturnValue(invalidSpec);

    await expect(validateSpec('invalid.yaml')).rejects.toThrow('Validation failed');
  });
});