# Domain Model

## Customer

| Attribute | Type | Description |
|-----------|------|-------------|
| id | UUID | Unique identifier |
| email | Email | Customer email address |
| firstName | String | First name |
| lastName | String | Last name |
| status | CustomerStatus | Account status (unverified, verified, suspended) |
