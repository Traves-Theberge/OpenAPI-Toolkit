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
    validateSchema?: boolean;
    retry?: string;
}
export declare function runTests(specPath: string, baseUrl: string, options?: TestOptions): Promise<void>;
export {};
//# sourceMappingURL=test.d.ts.map