# MedBook Squad Customization Summary

> Comprehensive overview of all Squad customizations for the MedBook healthcare platform.
>
> **Project:** medbook · **Team:** medbook-healthcare-team · **Version:** 2.0.0

---

## Table of Contents

- [1. Agents](#1-agents)
  - [1.1 Standard Agents](#11-standard-agents)
  - [1.2 Custom Healthcare Agents](#12-custom-healthcare-agents)
- [2. Custom Copilot Skill — HIPAA Compliance](#2-custom-copilot-skill--hipaa-compliance)
- [3. Ceremonies](#3-ceremonies)
  - [3.1 Event-Driven Ceremonies](#31-event-driven-ceremonies)
  - [3.2 Scheduled Ceremonies](#32-scheduled-ceremonies)
  - [3.3 Ceremony Settings](#33-ceremony-settings)
- [4. Quality Gates & Workflows](#4-quality-gates--workflows)
  - [4.1 Quality Gates](#41-quality-gates)
  - [4.2 Workflows](#42-workflows)
  - [4.3 CI/CD Integrations](#43-cicd-integrations)

---

## 1. Agents

The MedBook team extends the standard Squad agent roster with three domain-specific agents purpose-built for healthcare compliance.

### 1.1 Standard Agents

| Agent | Role | Responsibilities |
|-------|------|------------------|
| **Brain** | Architect & Planner | System architecture, feature planning, technology decisions. Designs with HIPAA compliance, PHI encryption, audit logging, and disaster recovery from the start. |
| **Eyes** | Code Reviewer | Code review for bugs, security, quality, error handling, test coverage, and performance. Enforces Go best practices. |
| **Hands** | Implementation Engineer | Feature implementation, bug fixes, refactoring, and test writing. Writes idiomatic Go with proper dependency injection and error handling. |
| **Mouth** | Documentation Specialist | API documentation, user guides, architectural decision records, runbooks, and README files. |
| **Ralph** | Persistent Memory Agent | Maintains context across sessions. Tracks decisions and progress history. |
| **Scribe** | Documentation Historian | Documentation specialist maintaining history, decisions, and technical records across the project. |
| **HIPAA Auditor** | Security Compliance Auditor | Independent security gate enforcing HIPAA Privacy Rule, Security Rule, and Breach Notification Rule. Reviews every code change touching PHI against a comprehensive checklist. Cites specific HIPAA regulations (e.g., §164.312) in findings. |

### 1.2 Custom Healthcare Agents

These agents are unique to the MedBook project (`custom: true` in `team.yml`).

#### HIPAA Compliance Agent

| Property | Value |
|----------|-------|
| **Config** | `.squad/agents/hipaa-compliance/config.yml` |
| **Charter** | `.squad/agents/hipaa-compliance/charter.md` |
| **Priority** | Critical |
| **Blocking** | Yes (pre-merge) |

**Purpose:** Ensures all code changes comply with the HIPAA Privacy and Security Rules when handling Protected Health Information (PHI).

**Capabilities:**
- PHI exposure detection in logs, error messages, and API responses
- Encryption verification (TLS 1.2+ in transit, column-level encryption at rest)
- Audit logging validation (every PHI access must produce an audit entry)
- Access control review (auth middleware on all PHI endpoints)
- Minimum necessary data access checks (no `SELECT *` on PHI tables)

**Triggers:** Files under `internal/patient/`, `internal/appointment/`, `internal/provider/`, `proto/`, and `*_test.go`; keywords: PHI, patient, diagnosis, SSN, MRN.

#### Terminology Agent

| Property | Value |
|----------|-------|
| **Config** | `.squad/agents/terminology/config.yml` |
| **Charter** | `.squad/agents/terminology/charter.md` |
| **Priority** | High |

**Purpose:** Ensures consistent and correct use of healthcare terminology, medical coding standards, and industry-specific field naming conventions.

**Capabilities:**
- ICD-10 diagnosis code format validation (`^[A-Z][0-9]{2}(\.[0-9]{1,4})?$`)
- CPT procedure code validation (`^\d{5}$`)
- Healthcare field naming enforcement (e.g., `patient_id` not `user_id`, `provider_id` not `doctor_id`)
- HL7 data exchange standard checks
- NPI (National Provider Identifier) validation

**Standards enforced:** ICD-10-CM, CPT, HL7, NPI, FHIR resource naming conventions.

#### Data Anonymizer Agent

| Property | Value |
|----------|-------|
| **Config** | `.squad/agents/anonymizer/config.yml` |
| **Charter** | `.squad/agents/anonymizer/charter.md` |
| **Priority** | High |

**Purpose:** Ensures test fixtures and development data use properly anonymized patient information, preventing accidental PHI in non-production environments.

**Capabilities:**
- HIPAA Safe Harbor de-identification compliance (18 identifier types)
- Test data generation with anonymized patterns
- PHI detection in test files and fixtures
- Validation of anonymization conventions

**Required test data patterns:**

| Field | Pattern | Example |
|-------|---------|---------|
| Patient name | `"Test "` prefix | `"Test Alice Smith"` |
| MRN | `"MRN-T"` prefix + 7 digits | `"MRN-T1234567"` |
| SSN | `"000-00-"` prefix | `"000-00-1234"` |
| Phone | `"(555) 01"` prefix | `"(555) 0145"` |
| Email | `@example.com` domain | `"patient.test42@example.com"` |

---

## 2. Custom Copilot Skill — HIPAA Compliance

| Property | Value |
|----------|-------|
| **Location** | `.copilot/skills/hipaa-compliance/SKILL.md` |
| **Domain** | Security, compliance, healthcare |
| **Confidence** | High |
| **Source** | Earned (hipaa-auditor charter, hipaa-compliance checklist, medbook codebase patterns) |

This is a custom Copilot skill that codifies machine-checkable HIPAA compliance patterns for the MedBook Go codebase. It covers six primary detection areas:

### Detection Patterns

| # | Pattern | Description |
|---|---------|-------------|
| 1 | **PHI Exposure in Logs** | Detects PHI fields (`Name`, `SSN`, `MRN`, `DOB`, etc.) referenced in `log.*`, `fmt.Print*`, `slog.*`, and `zap.*` calls. |
| 2 | **PHI in HTTP/gRPC Responses** | Flags PHI in `c.JSON()`, `gin.H{}`, `fmt.Errorf()`, and `status.Errorf()`. Verifies SSN uses `json:"-"` tags and responses use purpose-specific DTOs. |
| 3 | **Audit Logging Gaps** | Verifies every handler accessing PHI (`CreatePatient`, `GetPatient`, `UpdatePatient`, etc.) emits an `audit.Log` entry with required fields: `user_id`, `timestamp`, `action`, `resource_type`, `resource_id`, `outcome`. |
| 4 | **Missing Access Control Middleware** | Flags PHI routes (`/patients`, `/appointments`, `/providers`) registered without `auth.RequireToken()` and `auth.RequireRole()` middleware. |
| 5 | **Encryption Gaps** | Detects `grpc.WithInsecure()`, `ListenAndServe()` without TLS, missing `sslmode=require` on Postgres connections, hardcoded keys, and plaintext Kubernetes ConfigMaps for PHI. |
| 6 | **Hardcoded Patient Data in Tests** | Enforces anonymization prefixes (`Test `, `MRN-T`, `000-00-`, `(555) 01`, `@example.com`) in test files. Flags realistic patient data patterns. |

### Anti-Patterns Detected

- Logging PHI "temporarily for debugging"
- Using `SELECT *` on PHI tables
- Returning full entity models instead of DTOs from API endpoints
- Adding PHI endpoints without auth middleware
- Using `ListenAndServe()` without TLS on PHI-serving ports
- Writing audit entries via `log.Println` instead of a dedicated audit sink
- Storing encryption keys alongside encrypted data

---

## 3. Ceremonies

MedBook defines both event-driven ceremonies (in `.squad/ceremonies.md`) and scheduled ceremonies (in `.squad/ceremonies/ceremonies.yml`).

### 3.1 Event-Driven Ceremonies

| Ceremony | Trigger | Facilitator | Participants | Key Focus |
|----------|---------|-------------|--------------|-----------|
| **Design Review** | Multi-agent task involving 2+ agents on shared systems | Lead | All relevant | Interface contracts, risk identification, action items |
| **Retrospective** | Build failure, test failure, or reviewer rejection | Lead | All involved | Root cause analysis, process improvements |
| **Security Review** | PR touches `internal/patient/` or `proto/` | HIPAA Auditor | hipaa-auditor, eyes, hands | PHI exposure review, encryption check, audit logging validation. **Blocks merge** until hipaa-auditor signs off. |
| **HIPAA Sprint Compliance Check** | End of sprint (biweekly Friday) | HIPAA Auditor | hipaa-auditor, hipaa-compliance, brain, eyes | Compliance findings summary, remediation tracking, scorecard (target ≥ 95%). Produces `sprint_compliance_summary.md`, `remediation_tracker.yml`, `compliance_scorecard.yml`. |
| **Terminology Standup** | Weekly (Monday) | Terminology | terminology, brain, hands | New medical field names (FHIR validation), ICD-10 code review, CPT code review. **Blocks merge** on critical mismatches. |

### 3.2 Scheduled Ceremonies

| Ceremony | Frequency | Day/Time | Duration | Participants | Purpose |
|----------|-----------|----------|----------|--------------|---------|
| **Daily Standup (Healthcare)** | Daily | 09:30 UTC | 15 min | All agents | Standard standup plus compliance blockers and terminology questions |
| **Sprint Planning with Compliance** | Biweekly | Monday 09:00 | 2 hours | brain, hipaa-compliance, terminology, hands (required) | Backlog review, PHI feature identification, compliance overhead estimation, checkpoint assignment |
| **Weekly HIPAA Review** | Weekly | Friday 15:00 | 30 min | hipaa-compliance, eyes, brain (required) | Review PHI handling in merged PRs, update HIPAA checklist, share compliance lessons |
| **Retrospective with Compliance** | Biweekly | Friday 14:00 | 1 hour | All agents | Standard retro plus dedicated time for HIPAA challenges and terminology standardization wins |
| **Monthly Compliance Audit** | Monthly | First Monday 10:00 | 1 hour | hipaa-compliance, brain, eyes (required) | Full codebase PHI scan, audit log completeness, encryption review, access control validation. Produces monthly report, scorecard, and remediation plan. |

### 3.3 Ceremony Settings

- **Timezone:** UTC
- **Notifications:** Slack (`#medbook-squad`), 1-hour reminder before ceremony
- **Recording:** Enabled with 90-day retention
- **Attendance:** Required; max 2 consecutive skips before notification

---

## 4. Quality Gates & Workflows

Defined in `.squad/team.yml`.

### 4.1 Quality Gates

#### Pre-Commit

| Agent | Check | Required | Description |
|-------|-------|----------|-------------|
| hipaa-compliance | `phi_exposure_scan` | ✅ | Scan for PHI exposure in code changes |
| anonymizer | `test_data_validation` | ✅ | Validate test data is properly anonymized |

#### Code Review

| Agent | Check | Required | Blocking | Description |
|-------|-------|----------|----------|-------------|
| eyes | `standard_review` | ✅ | No | Standard code quality review |
| hipaa-compliance | `compliance_review` | ✅ | **Yes** | HIPAA compliance validation |
| terminology | `terminology_validation` | ✅ | No | Healthcare terminology and coding standards |

#### Pre-Merge

| Agent | Check | Required | Blocking | Description |
|-------|-------|----------|----------|-------------|
| hipaa-compliance | `final_compliance_check` | ✅ | **Yes** | Final HIPAA compliance verification |
| brain | `architecture_approval` | ✅ | No | Architecture and design approval |

### 4.2 Workflows

#### Feature Development (Standard)

```
Planning → Implementation → Review → Documentation
```

| Stage | Agents | Description |
|-------|--------|-------------|
| Planning | brain, hipaa-compliance | Plan feature with compliance considerations |
| Implementation | hands, terminology | Implement with proper healthcare terminology |
| Review | eyes, hipaa-compliance, terminology | Code review with compliance and terminology checks |
| Documentation | mouth | Document the feature |

#### Hotfix (Fast-Track)

```
Implementation & Review (combined)
```

| Stage | Agents | Description |
|-------|--------|-------------|
| Implementation & Review | hands, eyes, hipaa-compliance | Fast-tracked implementation with mandatory compliance — HIPAA compliance is **never** skipped, even for hotfixes |

### 4.3 CI/CD Integrations

GitHub Actions is configured with the following required checks:

- `hipaa-compliance` — HIPAA compliance validation
- `terminology-validation` — Healthcare terminology standards
- `test-data-anonymization` — Test data anonymization verification

---

## Summary of Customizations

| Category | Standard | Custom (Healthcare) |
|----------|----------|---------------------|
| **Agents** | brain, eyes, hands, mouth, ralph, scribe | **hipaa-compliance**, **terminology**, **anonymizer**, hipaa-auditor |
| **Copilot Skills** | 28 standard skills | **hipaa-compliance** (custom) |
| **Ceremonies** | Design Review, Retrospective | **Security Review**, **HIPAA Sprint Compliance Check**, **Terminology Standup**, **Weekly HIPAA Review**, **Sprint Planning with Compliance**, **Daily Standup (Healthcare)**, **Retrospective with Compliance**, **Monthly Compliance Audit** |
| **Quality Gates** | Standard code review | **PHI exposure scan**, **test data anonymization**, **HIPAA compliance review (blocking)**, **terminology validation**, **final compliance check (blocking)** |
| **Workflows** | Feature development | Feature dev with compliance stages, **hotfix with mandatory HIPAA** |
| **CI/CD Checks** | Standard CI | **hipaa-compliance**, **terminology-validation**, **test-data-anonymization** |

---

*Generated from Squad configuration files in `.squad/` and `.copilot/` directories.*
