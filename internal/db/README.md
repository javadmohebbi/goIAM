# 📘 goIAM Database Structure

This document describes the database schema used by goIAM — a multi-tenant Identity and Access Management (IAM) system. It includes users, roles, groups, policies, and a normalized policy engine for fine-grained access control.

---

## 🏢 Organization

Represents a tenant in the IAM system. All other entities (users, groups, roles, policies) are scoped by `Organization`.

**Fields:**
- `ID`
- `Name` — display name
- `Slug` — short unique identifier (e.g. `acme-corp`)
- `Description`

**Relations:**
- Has many `Users`, `Groups`, `Roles`, and `Policies`.

---

## 👤 User

Represents an individual account in an organization. A user can inherit permissions via roles, groups, and direct policies.

**Fields:**
- `Username`, `Email` — unique within organization
- `PasswordHash`
- `FirstName`, `MiddleName`, `LastName`, `Address`
- `PhoneNumber`, `EmailVerified`, `PhoneVerified`
- `TOTPSecret`, `Requires2FA`, `IsActive`
- `OrganizationID` — foreign key to `Organization`

**Relations:**
- Belongs to one `Organization`
- Many-to-many with `Groups`, `Roles`, and `Policies`
- Has many `BackupCodes`

---

## 🧑‍🤝‍🧑 Group

Used to manage users as a unit for permission assignment.

**Fields:**
- `Name`, `Slug` — unique within org
- `Description`
- `OrganizationID`

**Relations:**
- Many-to-many with `Users`
- Many-to-many with `Policies`

---

## 🛡️ Role

Represents a set of permissions, typically by job function or access level.

**Fields:**
- `Name`, `Slug` — unique within org
- `Description`
- `OrganizationID`

**Relations:**
- Many-to-many with `Users`
- Many-to-many with `Policies`

---

## 📜 Policy

Defines access control logic in a structured format.

**Fields:**
- `Name`, `Slug` — unique within org
- `Description`
- `OrganizationID`

**Relations:**
- Has many `PolicyStatements`
- Many-to-many with `Users`, `Groups`, and `Roles`

---

## 🧾 PolicyStatement

Represents a rule like: *"Allow `user:create` on `org:1:user:*`"*.

**Fields:**
- `Effect` — "Allow" or "Deny"
- `PolicyID` — belongs to a `Policy`

**Relations:**
- Has many `PolicyActions` and `PolicyResources`

---

## 🎯 PolicyAction

Specifies actions controlled by a policy statement.

**Fields:**
- `Action` — e.g. `user:create`, `role:assign`, `*`
- `PolicyStatementID`

---

## 🧩 PolicyResource

Specifies the resources targeted by a policy statement.

**Fields:**
- `Resource` — e.g. `org:1:user:123`, `*`
- `PolicyStatementID`

---

## 🔐 BackupCode

One-time recovery codes for 2FA.

**Fields:**
- `UserID`
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
        |                        |     Group, Role       |
        |                        +-----------------------+
        |                                     |
        |                                     | many-to-many
        |                                     v
        |                            +---------------+
        |                            |    Policy     |
        |                            +---------------+
        |                                     |
        |                                     | 1
        |                                     v
        |                            +--------------------+
        |                            | PolicyStatement    |
        |                            +--------------------+
        |                                |            |
        |                                | *          | *
        |                                v            v
        |                          +-----------+   +--------------+
        |                          |  Action   |   |   Resource   |
        |                          +-----------+   +--------------+
```

---

## 🧠 How Permissions Are Evaluated

When a user attempts an action on a resource, the system calls:

```go
allowed := db.EvaluatePolicy(user, "user:create", "org:123:user:456")
```

It checks:
1. Policies assigned directly to the user
2. Policies inherited from the user's groups and roles
3. All matching `PolicyStatements` with that action and resource
4. If any statement `Deny`s access → access is blocked
5. If no `Deny`, but at least one `Allow` → access is granted

Wildcards (`*`) are supported for actions and resources.

---