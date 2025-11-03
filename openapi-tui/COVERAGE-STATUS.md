# Test Coverage Status

## Current Coverage (as of latest commit)

| Package | Coverage | Status |
|---------|----------|--------|
| internal/errors | **100.0%** | ✅ Complete |
| internal/ui | **96.7%** | ✅ Excellent |
| internal/validation | **94.3%** | ✅ Very Good |
| internal/export | **93.0%** | ✅ Very Good |
| internal/models | **90.2%** | ✅ Good |
| internal/config | **88.1%** | ⚠️ Good |
| internal/testing | **80.3%** | ⚠️ Needs Work |

## Package Details

### internal/errors (100.0%) ✅
- Complete coverage achieved
- All error enhancement functions tested
- No further work needed

### internal/ui (96.7%) ✅
- Comprehensive view testing
- Model initialization tested
- Stats formatting tested
- Remaining 3.3% likely edge cases in Bubble Tea event handlers

**To reach 100%:**
- Add more edge case tests for ViewEndpointSelector states
- Test more InitialConfigEditorModel variations
- May require mocking Bubble Tea internals

### internal/validation (94.3%) ✅
- ValidateSpec: 100%
- ValidateResponse: 79.3%
- validateStatusCode: 90%
- All other functions: 100%

**To reach 100%:**
- The empty content-type branch is hard to test because the code auto-defaults to "application/json"
- Some branches in default response handling may be library-dependent (kin-openapi)
- May need to accept 94-95% as "complete" given library behavior

### internal/export (93.0%) ✅
- ExportResults: 88.2%
- ExportResultsToFile: 86.7%
- ExportResultsToHTML: 92.7%
- ExportResultsToJUnit: 93.9%

**To reach 100%:**
- Add tests for file write errors (permission denied, disk full)
- Test HTML template edge cases
- Test JUnit XML formatting edge cases
- Most uncovered branches are error paths

### internal/models (90.2%) ✅
- Good coverage of core functionality
- History save/load tested including error paths

**To reach 100%:**
- Some OS-level errors are hard to test (os.UserHomeDir failure)
- Need tests for JSON marshal/unmarshal errors
- May require build tags or accepting 90% for OS-dependent code

### internal/config (88.1%) ⚠️
- LoadConfig, SaveConfig tested
- Error paths tested

**To reach 100%:**
- GetConfigPath has OS-dependent error paths (UserHomeDir)
- Some filesystem errors hard to trigger without root/special permissions
- May need to accept ~88% or use build tags/mocking

### internal/testing (80.3%) ⚠️
- **RunTests function: 9.8%** ← Main gap
- All other functions: 78-96%

**To reach 100%:**
- RunTests is the main integration function - needs comprehensive integration tests
- Requires real OpenAPI spec file + mock HTTP server
- Need to test all code paths: multiple endpoints, auth types, verbose mode, retries, parallel execution
- This is the most complex function to test due to its integration nature

**Recommended approach for RunTests:**
1. Create fixture OpenAPI spec in testdata/
2. Set up httptest server with multiple endpoints
3. Test sequential flow with various configurations
4. Test error paths (invalid spec, network failures)
5. Test validation integration

## Overall Assessment

**Current state:** 90.7% average coverage (weighted by package size)

**Target achievability:**
- **Realistically achievable: 95%** across the board with moderate effort
- **100% achievable for:** errors ✅, ui, validation, export
- **100% difficult for:** models, config (OS-dependent), testing (integration complexity)

## Recommended Next Steps

1. **Priority 1: Fix `internal/testing` (80.3% → 90%+)**
   - Add integration tests for RunTests function
   - This will have the biggest impact on overall coverage

2. **Priority 2: Polish high-coverage packages (94-97% → 98-100%)**
   - validation: Add empty content-type scenario tests
   - ui: Add more edge case tests
   - export: Add error path tests

3. **Priority 3: Accept practical limits**
   - OS-dependent code in config/models may stay at 88-90%
   - Document why certain lines can't be tested without mocking OS internals

## Notes

- All tests passing: ✅
- Test execution time: ~100s (mostly internal/testing integration tests)
- No flaky tests observed
- Good test organization with clear naming conventions
