# Change: Update Auth Service Swagger Documentation to Align with Refactored Handler Structure

## Why

The auth-service has undergone architectural refactoring with the new handler structure (auth, pin, oauth, user handlers), but the Swagger documentation needs to be reviewed and potentially updated to ensure it accurately reflects the current API endpoints, request/response models, and is properly organized according to the new handler structure. This will ensure API documentation remains accurate and useful for developers.

## What Changes

- Review and update Swagger documentation to match current endpoint implementations
- Ensure all endpoints are properly documented with correct handler groupings
- Verify request/response schemas match actual implementation
- Update tags and organization to align with new handler structure (AuthHandler, PinHandler, OAuthHandler, UserHandler)
- Add missing error responses or update existing ones based on current implementation
- Ensure authentication requirements are correctly documented
- Update examples to reflect actual current behavior

## Impact

- Affected specs: `auth-service` API documentation
- Affected code: `auth-service/static/swagger.json`
- Documentation improvements:
  - Accurate API documentation reflecting current implementation
  - Proper handler grouping and organization
  - Complete request/response schemas
  - Correct error handling documentation
  - Updated examples and authentication requirements
- Developer experience: Improved API discovery and integration clarity
- No breaking changes - only documentation updates