# API Input Validator Memory

## Project Validation Stack
- Framework: Echo v5
- Validation Library: go-playground/validator/v10
- Validator configured in App.go as ValidatorWrapper
- Pattern: Use struct tags + c.Bind() + c.Validate()

## Validation Pattern
1. Define request struct with validation tags
2. Call c.Bind(&req) to parse request body
3. Call c.Validate(&req) to run validation
4. Return structured error responses with field-level details

## Reference Implementation
See handlers/Magiclink.go for complete validation pattern:
- Nested struct validation (MagicLinkMetadata)
- Comprehensive error formatting with field names
- Proper use of validator tags: required, email, url, min, max

## Identified Vulnerabilities

### handlers/Account.go - CRITICAL
- No validation tags on request struct
- Missing email format validation
- No input sanitization
- Missing length limits on name/email
- Services layer has weak validation (empty string checks only)

### services/Account.go
- Weak validation: only checks empty strings
- No email format validation before DB operations
- Hardcoded password "abc123" - security issue

### services/AccountUsers.go
- No email format validation
- No password strength validation
- Manual role validation instead of using validator tags

## Models Validation Tags Present
Models.go has validate tags on domain models but these are NOT enforced at handler layer currently.
