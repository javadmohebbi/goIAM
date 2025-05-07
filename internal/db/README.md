# 📘 goIAM Database Structure

This document outlines the database schema used by goIAM, a multi-tenant Identity and Access Management (IAM) system.

## 🏢 Organization

- Represents a tenant in the multi-tenant system.
- All users, groups, roles, and policies are scoped under a single organization.

**Fields:**
- `ID`
- `Name` (unique)

## 👤 User

- Represents an individual account in an organization.
- Can be assigned to groups, roles, and policies.
- Supports 2FA and backup codes.

**Fields:**
- `Username`, `Email` (unique within organization)
- `PasswordHash`
- `FirstName`, `MiddleName`, `LastName`, `Address`
- `TOTPSecret`, `Requires2FA`
- `OrganizationID` (foreign key)

**Relations:**
- Belongs to one Organization
- Many-to-many with Groups, Roles, Policies
- Has many BackupCodes

## 🧑‍🤝‍🧑 Group

- Logical grouping of users.
- Used for assigning policies at the group level.

**Fields:**
- `Name` (unique within organization)
- `OrganizationID`

**Relations:**
- Many-to-many with Users
- Many-to-many with Policies

## 🛡️ Role

- Represents a set of responsibilities or permissions.
- Used to apply policies at a role level.

**Fields:**
- `Name` (unique within organization)
- `OrganizationID`

**Relations:**
- Many-to-many with Users
- Many-to-many with Policies

## 📜 Policy

- Defines specific access control rules.
- Can be attached to users, groups, and roles.

**Fields:**
- `Name` (unique within organization)
- `Description`
- `OrganizationID`

**Relations:**
- Many-to-many with Users, Groups, Roles

## 🔐 BackupCode

- One-time recovery code for users with 2FA enabled.

**Fields:**
- `UserID` (foreign key)
- `CodeHash`
- `Used`

---

## 🔗 Entity Relationships

```plaintext
+----------------+     1        *     +-------------+
|  Organization  |------------------>|     User     |
+----------------+                   +-------------+
        |                                     |
        |                                     | many-to-many
        |                                     v
        |                        +-----------------------+
        |                        |     Group, Role,      |
        |                        |     or Policy         |
        |                        +-----------------------+
        |
        | 1
        v
+-------------------+
|     Policy        |
+-------------------+

User ⇄ Group ⇄ Policy  
User ⇄ Role ⇄ Policy  
User ⇄ Policy  
```

---

## 📚 How It Works Together

goIAM is built to support multiple organizations (tenants), each with complete isolation of users, roles, and policies.

- Each **Organization** can register its own users.
- A **User** belongs to a single organization but can be assigned to multiple **Groups**, **Roles**, and **Policies**.
- **Groups** and **Roles** serve as collections to attach **Policies** and simplify management.
- **Policies** define what actions can be performed or accessed.
- **BackupCodes** provide fallback authentication for users with 2FA enabled.

All IAM operations—like role enforcement, access control, and audit logging—are scoped by the user's organization to ensure strong isolation and security between tenants.