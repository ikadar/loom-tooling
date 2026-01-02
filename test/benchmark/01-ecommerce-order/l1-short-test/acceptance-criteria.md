# Acceptance Criteria

## AC-CUST-001 – Customer Registration

**Given** a visitor with email 'john@example.com' that is not registered
**When** they submit registration with email 'john@example.com', password 'SecurePass1', first name 'John', and last name 'Doe'
**Then** a customer account is created with status 'unverified', and a verification email is sent to 'john@example.com'

**Error Cases:**
- Email already registered → DUPLICATE_EMAIL error, suggest login or password reset
- Password less than 8 characters → INVALID_PASSWORD error with requirements message
- Password without number → INVALID_PASSWORD error with requirements message
- Invalid email format → INVALID_EMAIL error
