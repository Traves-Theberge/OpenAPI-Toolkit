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
    quiet?: boolean;
}
export declare function runTests(specPath: string, baseUrl: string, options?: TestOptions): Promise<void>;
export {};
//# sourceMappingURL=test.d.ts.map