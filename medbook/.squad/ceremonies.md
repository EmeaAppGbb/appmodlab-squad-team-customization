# Ceremonies

> Team meetings that happen before or after work. Each squad configures their own.

## Design Review

| Field | Value |
|-------|-------|
| **Trigger** | auto |
| **When** | before |
| **Condition** | multi-agent task involving 2+ agents modifying shared systems |
| **Facilitator** | lead |
| **Participants** | all-relevant |
| **Time budget** | focused |
| **Enabled** | ✅ yes |

**Agenda:**
1. Review the task and requirements
2. Agree on interfaces and contracts between components
3. Identify risks and edge cases
4. Assign action items

---

## Retrospective

| Field | Value |
|-------|-------|
| **Trigger** | auto |
| **When** | after |
| **Condition** | build failure, test failure, or reviewer rejection |
| **Facilitator** | lead |
| **Participants** | all-involved |
| **Time budget** | focused |
| **Enabled** | ✅ yes |

**Agenda:**
1. What happened? (facts only)
2. Root cause analysis
3. What should change?
4. Action items for next iteration

---

## Security Review

| Field | Value |
|-------|-------|
| **Trigger** | auto |
| **When** | before |
| **Condition** | PR touches files in `internal/patient/` or `proto/` (patient data models) |
| **Facilitator** | hipaa-auditor |
| **Participants** | hipaa-auditor, eyes, hands |
| **Time budget** | focused |
| **Enabled** | ✅ yes |

**Agenda:**
1. **PHI Exposure Review** — Scan changed files for any new or modified fields that store, transmit, or display Protected Health Information. Verify no PHI leaks into logs, error messages, API responses, or analytics payloads.
2. **Encryption Check** — Validate that all PHI fields use encryption at rest and in transit. Confirm new proto fields carrying patient data are annotated for encryption and that Go struct tags enforce redaction in serialization.
3. **Audit Logging Validation** — Ensure every read, write, and delete operation on patient data models emits a compliant audit log entry with actor identity, timestamp, resource ID, and action type.

**Gate:** PR cannot merge until the hipaa-auditor agent signs off on all three checks. Failures produce a `security-review/blocked` status check.

---

## HIPAA Sprint Compliance Check

| Field | Value |
|-------|-------|
| **Trigger** | scheduled |
| **When** | after |
| **Condition** | end of sprint (biweekly, Friday) |
| **Facilitator** | hipaa-auditor |
| **Participants** | hipaa-auditor, hipaa-compliance, brain, eyes |
| **Time budget** | extended |
| **Enabled** | ✅ yes |

**Agenda:**
1. **Compliance Findings Summary** — Aggregate all Security Review findings, HIPAA checklist violations, and audit-log gaps discovered during the sprint. Categorize by severity (critical, high, medium, low).
2. **Remediation Status Tracking** — Review open remediation items from this and previous sprints. Update status (open → in-progress → resolved → verified). Escalate overdue items.
3. **Compliance Scorecard Generation** — Produce a scorecard covering: PHI exposure incidents, encryption coverage %, audit-log completeness %, open vs. closed findings, and overall compliance score (target ≥ 95%).

**Deliverables:**
- `sprint_compliance_summary.md` — narrative summary of all findings
- `remediation_tracker.yml` — machine-readable status of every open item
- `compliance_scorecard.yml` — scored metrics for the sprint

---

## Terminology Standup

| Field | Value |
|-------|-------|
| **Trigger** | scheduled |
| **When** | during |
| **Condition** | weekly (Monday) |
| **Facilitator** | terminology |
| **Participants** | terminology, brain, hands |
| **Time budget** | focused |
| **Enabled** | ✅ yes |

**Agenda:**
1. **New Medical Field Names** — Review any new or renamed fields in data models, proto definitions, and API contracts introduced during the past week. Validate naming against FHIR resource conventions and project glossary.
2. **ICD-10 Code Review** — Verify newly introduced ICD-10 diagnosis codes are valid, correctly versioned (ICD-10-CM), and mapped to the right domain entities.
3. **CPT Code Review** — Verify newly introduced CPT procedure codes are valid, properly licensed for use, and correctly associated with billing or encounter models.

**Gate:** The terminology agent flags non-standard names or invalid codes as warnings. Critical mismatches (e.g., deprecated ICD-10 codes, unlicensed CPT codes) block merge until resolved.
