---
lab_id: squad-team-customization
title: "SQUAD Team Customization: Building Domain-Specific Development Agents"
category: Agentic Software Development
difficulty: intermediate
duration: 3-4 hours
technologies:
  - GitHub Copilot
  - SQUAD Framework
  - Go
  - YAML
learning_objectives:
  - Define custom SQUAD agents with specialized roles and skills
  - Configure agent charters with domain-specific review criteria
  - Set up custom ceremonies (standups, retrospectives, compliance reviews)
  - Create domain-specific quality gates in the SQUAD workflow
  - Test and validate custom agent behavior
prerequisites:
  - Completed "Getting Started with SQUAD" lab
  - Basic Go language familiarity
  - Understanding of YAML configuration
  - GitHub Copilot access
repository: https://github.com/EmeaAppGbb/appmodlab-squad-team-customization
organization: EmeaAppGbb
---

# SQUAD Team Customization: Building Domain-Specific Development Agents

## Overview

This lab teaches you how to customize a SQUAD team by adding custom agents, defining specialized skills, configuring ceremonies (retrospectives, standups), and adapting agent behavior for domain-specific workflows. You'll learn that SQUAD is not a rigid framework but a flexible platform that can be tailored to your organization's unique development culture, tech stack, and quality requirements.

### What You'll Build

You'll customize a SQUAD team for **MedBook**, a healthcare appointment scheduling platform built with Go microservices. You'll add three domain-specific agents:

1. **HIPAA Compliance Agent** — Reviews code for PHI handling violations
2. **Domain Terminology Agent** — Validates healthcare terminology usage
3. **Data Anonymizer Agent** — Generates/validates anonymized test data

You'll also configure custom ceremonies (weekly security reviews, sprint planning) and define domain-specific quality gates.

### Why This Matters

Real-world software development often involves domain-specific requirements:
- **Compliance:** Healthcare (HIPAA), finance (PCI-DSS), government (FedRAMP)
- **Industry Standards:** Medical terminology, financial regulations, safety-critical systems
- **Company Culture:** Specific review processes, quality gates, documentation standards

SQUAD's extensibility allows you to encode this domain knowledge into your development workflow.

## Architecture

### Business Context

MedBook is a healthcare SaaS platform that helps medical practices manage appointment scheduling. The system handles:
- Patient registration and profiles
- Healthcare provider schedules
- Appointment booking and management
- Insurance verification
- HIPAA-compliant data handling

### Technical Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (Ingress)                     │
└──────────────────┬──────────────────┬──────────────────────┘
                   │                  │                  
         ┌─────────▼────────┐  ┌─────▼──────────┐  ┌──────────────┐
         │ Patient Service  │  │ Provider Svc   │  │ Appointment  │
         │   (gRPC/HTTP)    │  │  (gRPC/HTTP)   │  │   Service    │
         └─────────┬────────┘  └─────┬──────────┘  └──────┬───────┘
                   │                  │                     │
                   └──────────────────┴─────────────────────┘
                                      │
                              ┌───────▼────────┐
                              │   PostgreSQL   │
                              │  (Patient PHI) │
                              └────────────────┘
```

### SQUAD Customization Architecture

```
Standard SQUAD Agents              Custom Healthcare Agents
┌──────────────────┐              ┌─────────────────────────┐
│  Brain (AI)      │              │  HIPAA Compliance       │
│  - Planning      │              │  - PHI exposure check   │
│  - Architecture  │              │  - Encryption validation│
└──────────────────┘              │  - Audit logging        │
                                  └─────────────────────────┘
┌──────────────────┐              
│  Eyes (Review)   │◄─────────────┐
│  - Code quality  │              │  Domain Terminology     │
│  - Best practices│              │  - ICD-10 validation    │
└──────────────────┘              │  - CPT code checks      │
                                  │  - Medical field names  │
┌──────────────────┐              └─────────────────────────┘
│  Hands (Code)    │              
│  - Implementation│              ┌─────────────────────────┐
│  - Refactoring   │              │  Data Anonymizer        │
└──────────────────┘              │  - Test data generation │
                                  │  - PHI scrubbing        │
┌──────────────────┐              │  - Fixture validation   │
│  Mouth (Docs)    │              └─────────────────────────┘
│  - Documentation │              
└──────────────────┘              
```

## Lab Structure

### Branch Strategy

- **`legacy`** — MedBook codebase with standard SQUAD configuration
- **`step-1-custom-agent-definition`** — Define custom agent charters and capabilities
- **`step-2-agent-implementation`** — Implement custom agent configurations
- **`step-3-ceremonies`** — Define and configure custom ceremonies
- **`step-4-skills-and-gates`** — Add domain-specific skills and quality gates
- **`step-5-integration-test`** — Run customized SQUAD on a real development task
- **`solution`** — Complete implementation with all customizations

## Prerequisites Setup

Before starting the lab, ensure you have:

```bash
# Required tools
go version          # Go 1.22+
docker --version    # For local PostgreSQL
kubectl version     # For K8s manifest validation
```

Clone the repository and checkout the `legacy` branch:

```bash
git clone https://github.com/EmeaAppGbb/appmodlab-squad-team-customization
cd appmodlab-squad-team-customization
git checkout legacy
```

## Step 1: Review Base SQUAD Configuration

**Branch:** `legacy`

**Objective:** Understand the standard SQUAD setup before customization.

### 1.1 Explore the Codebase

```bash
# Review the directory structure
ls -la medbook/

# Examine the base SQUAD configuration
cat medbook/.squad/team.yml
```

The base configuration includes four standard agents:
- **Brain** — Architecture and planning
- **Eyes** — Code review and quality
- **Hands** — Implementation
- **Mouth** — Documentation

### 1.2 Review Agent Charters

```bash
# View standard agent configurations
cat medbook/.squad/agents/brain/charter.md
cat medbook/.squad/agents/eyes/charter.md
```

**Key Observation:** Standard agents have generic software engineering capabilities. They don't understand healthcare-specific requirements like HIPAA compliance or medical terminology.

### 1.3 Run Standard SQUAD

```bash
# Test the base SQUAD configuration
cd medbook
copilot squad start

# In the SQUAD session, try:
# "Review the patient service for security issues"
```

**Expected Result:** SQUAD identifies general security issues but misses healthcare-specific concerns like PHI logging, missing encryption, or improper patient data handling.

## Step 2: Define HIPAA Compliance Agent

**Branch:** `step-1-custom-agent-definition`

**Objective:** Create a custom agent charter for HIPAA compliance review.

### 2.1 Create Agent Charter

Create `medbook/.squad/agents/hipaa-compliance/charter.md`:

```markdown
# HIPAA Compliance Agent

## Role
I am the HIPAA Compliance Agent, responsible for ensuring all code changes comply with HIPAA Privacy and Security Rules when handling Protected Health Information (PHI).

## Responsibilities
- Review code for PHI exposure in logs, error messages, or responses
- Verify encryption of PHI in transit and at rest
- Ensure proper audit logging of PHI access
- Validate patient consent handling
- Check for minimum necessary data access
- Verify secure authentication and authorization

## PHI Fields to Monitor
- Patient name, SSN, MRN (Medical Record Number)
- Date of birth, address, phone number
- Email, insurance information
- Diagnosis codes, treatment notes
- Any field in the `patients` table

## Review Checklist
- [ ] No PHI in log statements
- [ ] PHI encrypted in database (at rest)
- [ ] TLS required for PHI transmission
- [ ] Audit log created for PHI access
- [ ] Role-based access control enforced
- [ ] Patient data minimized in responses
- [ ] Error messages don't leak PHI

## When to Trigger
- Any changes to patient, appointment, or provider services
- Database schema modifications
- API response modifications
- Logging or error handling changes
```

### 2.2 Create Agent Configuration

Create `medbook/.squad/agents/hipaa-compliance/config.yml`:

```yaml
agent:
  name: hipaa-compliance
  role: HIPAA Compliance Reviewer
  description: Ensures code changes comply with HIPAA regulations
  
  triggers:
    - file_pattern: "internal/patient/**"
    - file_pattern: "internal/appointment/**"
    - file_pattern: "proto/**"
    - keyword: "PHI"
    - keyword: "patient"
    - keyword: "diagnosis"
  
  capabilities:
    - hipaa_privacy_rule_validation
    - phi_exposure_detection
    - encryption_verification
    - audit_logging_validation
  
  review_priority: critical
  
  integration_points:
    - stage: pre_commit
      action: validate_phi_handling
    - stage: code_review
      action: hipaa_compliance_check
    - stage: pre_merge
      action: require_approval
```

### 2.3 Create Compliance Checklist

Create `medbook/.squad/agents/hipaa-compliance/checklist.yml`:

```yaml
hipaa_compliance_checklist:
  phi_exposure:
    - rule: "No PHI in log.Info, log.Debug, log.Error"
      severity: critical
      pattern: 'log\.(Info|Debug|Error).*patient\.(Name|SSN|DOB)'
      
    - rule: "No PHI in HTTP error responses"
      severity: critical
      pattern: 'c\.JSON.*patient\.(Name|SSN)'
      
  encryption:
    - rule: "Database connections must use SSL"
      severity: critical
      check: "connection string contains sslmode=require"
      
    - rule: "API must enforce TLS 1.2+"
      severity: critical
      check: "TLS config set to minimum TLS 1.2"
      
  audit_logging:
    - rule: "PHI access must be logged"
      severity: high
      pattern: 'GetPatient.*audit\.Log'
      
    - rule: "Audit logs must include user ID and timestamp"
      severity: high
      check: "audit log contains userID and timestamp"
      
  access_control:
    - rule: "PHI endpoints require authentication"
      severity: critical
      check: "middleware includes auth.RequireToken"
      
    - rule: "Role-based access enforced"
      severity: high
      check: "authorization middleware present"
```

**Exercise:** Review a code change that adds patient logging. What HIPAA violations would this agent catch?

## Step 3: Define Domain Terminology Agent

**Branch:** `step-2-agent-implementation`

**Objective:** Create an agent that validates healthcare terminology and coding standards.

### 3.1 Create Terminology Agent Charter

Create `medbook/.squad/agents/terminology/charter.md`:

```markdown
# Domain Terminology Agent

## Role
I ensure consistent and correct use of healthcare terminology, medical coding standards, and industry-specific field naming conventions.

## Responsibilities
- Validate ICD-10 diagnosis codes
- Verify CPT procedure codes
- Ensure medical field names follow standards
- Check medical terminology spelling and usage
- Validate healthcare data types and formats

## Healthcare Standards
- **ICD-10:** International Classification of Diseases, 10th Revision
- **CPT:** Current Procedural Terminology
- **HL7:** Health Level 7 data exchange standards
- **SNOMED CT:** Systematized Nomenclature of Medicine

## Terminology Validation Rules
- Diagnosis codes must be valid ICD-10 format (e.g., "E11.9", "I10")
- Procedure codes must be valid CPT format (5 digits)
- Patient identifiers: MRN, SSN format validation
- Date formats: ISO 8601 for healthcare data
- Field names: Use medical standard terms (e.g., "diagnosis" not "problem")

## When to Trigger
- Changes to proto definitions
- Database schema changes
- API request/response models
- Medical record processing code
```

### 3.2 Create Terminology Validation Rules

Create `medbook/.squad/agents/terminology/validation-rules.yml`:

```yaml
terminology_validation:
  icd10_codes:
    pattern: '^[A-Z][0-9]{2}(\.[0-9]{1,4})?$'
    examples:
      valid: ["E11.9", "I10", "J44.0"]
      invalid: ["E11", "999", "ABC123"]
    description: "ICD-10 codes: Letter + 2 digits + optional decimal + up to 4 digits"
    
  cpt_codes:
    pattern: '^\d{5}$'
    examples:
      valid: ["99213", "80053"]
      invalid: ["9921", "ABC12"]
    description: "CPT codes: exactly 5 digits"
    
  medical_record_number:
    pattern: '^MRN-\d{8}$'
    examples:
      valid: ["MRN-12345678"]
      invalid: ["12345678", "MRN123"]
    description: "MRN format: MRN- prefix + 8 digits"
    
  field_naming_standards:
    patient_identifier:
      preferred: ["patient_id", "mrn"]
      avoid: ["user_id", "person_id"]
      
    diagnosis:
      preferred: ["diagnosis_code", "icd10_code"]
      avoid: ["problem", "condition_id"]
      
    provider:
      preferred: ["provider_id", "npi"]
      avoid: ["doctor_id", "physician"]
      
    appointment:
      preferred: ["appointment_datetime", "scheduled_at"]
      avoid: ["appt_time", "booking_date"]
```

## Step 4: Define Data Anonymizer Agent

**Objective:** Create an agent that generates and validates anonymized test data.

### 4.1 Create Anonymizer Agent Charter

Create `medbook/.squad/agents/anonymizer/charter.md`:

```markdown
# Data Anonymizer Agent

## Role
I ensure test fixtures and development data use properly anonymized patient information, preventing accidental use of real PHI in non-production environments.

## Responsibilities
- Generate realistic but fake patient data for tests
- Validate test fixtures don't contain real PHI
- Provide anonymization utilities for developers
- Review test data for compliance with de-identification rules

## Anonymization Rules (HIPAA Safe Harbor)
Remove or replace 18 types of identifiers:
1. Names (use fake names)
2. Geographic subdivisions smaller than state
3. Dates (except year) — use shifted dates
4. Telephone, fax numbers
5. Email addresses
6. SSN, MRN
7. Account numbers
8. Certificate/license numbers
9. Vehicle identifiers
10. Device IDs
11. Web URLs
12. IP addresses
13. Biometric identifiers
14. Photos
15. Any other unique identifier

## Test Data Standards
- Use faker library for consistent fake data
- Patient names: "Test [FirstName] [LastName]"
- MRN: "MRN-T" + 7 random digits
- SSN: "000-00-" + 4 random digits
- Dates: Current year - random offset
- Addresses: Use fake but realistic addresses
- Phone: (555) 0100-0199 range

## When to Trigger
- New test files created
- Test fixture modifications
- Database seed data changes
```

### 4.2 Create Anonymization Utilities

Create `medbook/.squad/agents/anonymizer/utilities.yml`:

```yaml
anonymization_utilities:
  fake_patient_generator:
    name:
      pattern: "Test {FirstName} {LastName}"
      library: "faker"
      
    mrn:
      pattern: "MRN-T{7_digits}"
      example: "MRN-T1234567"
      
    ssn:
      pattern: "000-00-{4_digits}"
      example: "000-00-1234"
      
    dob:
      method: "random_date_in_range"
      range: "1940-2020"
      format: "YYYY-MM-DD"
      
    phone:
      pattern: "(555) 01{2_digits}"
      example: "(555) 0145"
      
  validation_checks:
    - name: "No real SSNs"
      pattern: '^(?!000)(?!666)(?!9)\d{3}-(?!00)\d{2}-(?!0000)\d{4}$'
      action: "reject_if_valid_real_ssn"
      
    - name: "No real names"
      check: "name must start with 'Test '"
      
    - name: "Test MRN prefix"
      check: "mrn must start with 'MRN-T'"
```

**Exercise:** Write a test that generates 10 anonymized patient records. Verify none contain real PHI.

## Step 5: Configure Custom Ceremonies

**Branch:** `step-3-ceremonies`

**Objective:** Define custom team ceremonies for compliance and planning.

### 5.1 Create Ceremonies Configuration

Create `medbook/.squad/ceremonies/ceremonies.yml`:

```yaml
ceremonies:
  - name: weekly-hipaa-review
    type: compliance_review
    frequency: weekly
    day: friday
    duration: 30_minutes
    
    participants:
      required:
        - hipaa-compliance
        - eyes
        - brain
      optional:
        - hands
        - mouth
    
    agenda:
      - Review PHI handling in merged PRs
      - Discuss new compliance requirements
      - Update HIPAA checklist if needed
      - Share compliance findings and lessons learned
      
    artifacts:
      - compliance_report.md
      - action_items.yml
      
  - name: sprint-planning-with-compliance
    type: sprint_planning
    frequency: biweekly
    duration: 2_hours
    
    participants:
      required:
        - brain
        - hipaa-compliance
        - terminology
      optional:
        - all_agents
    
    agenda:
      - Review product backlog
      - Identify PHI-touching features
      - Estimate compliance overhead
      - Plan HIPAA reviews for high-risk features
      - Assign compliance checkpoints
      
    output:
      - sprint_plan.md
      - compliance_checkpoint_schedule.yml
      
  - name: daily-standup-healthcare
    type: standup
    frequency: daily
    duration: 15_minutes
    
    participants:
      required:
        - all_agents
    
    format:
      - What did you complete yesterday?
      - What will you work on today?
      - Any compliance blockers?
      - Any terminology questions?
      
  - name: retrospective-with-compliance
    type: retrospective
    frequency: biweekly
    duration: 1_hour
    
    participants:
      required:
        - all_agents
    
    topics:
      - What went well?
      - What could be improved?
      - HIPAA compliance challenges
      - Terminology standardization wins
      - Action items for next sprint
```

### 5.2 Create Ceremony Templates

Create `medbook/.squad/ceremonies/templates/weekly-hipaa-review.md`:

```markdown
# Weekly HIPAA Compliance Review
**Date:** {{date}}
**Duration:** 30 minutes
**Participants:** {{participants}}

## Agenda

### 1. Review Merged PRs (15 min)
- List PRs merged this week that touched PHI
- HIPAA compliance check results
- Any violations found and remediated

### 2. New Compliance Requirements (5 min)
- Updates to HIPAA regulations
- New organizational policies
- Industry best practices

### 3. Checklist Updates (5 min)
- Additions to compliance checklist
- Refinements based on findings

### 4. Lessons Learned (5 min)
- What compliance issues were caught?
- What was missed?
- How can we improve detection?

## Review Summary

| PR # | Feature | PHI Impact | Compliance Status | Issues Found | Resolved |
|------|---------|------------|-------------------|--------------|----------|
| | | | | | |

## Action Items

- [ ] {{action_item}}

## Notes

{{notes}}
```

## Step 6: Add Quality Gates

**Branch:** `step-4-skills-and-gates`

**Objective:** Configure domain-specific quality gates in the SQUAD workflow.

### 6.1 Update Team Configuration

Update `medbook/.squad/team.yml`:

```yaml
squad_team:
  name: medbook-healthcare-team
  version: 2.0.0
  
  agents:
    # Standard agents
    - name: brain
      enabled: true
      
    - name: eyes
      enabled: true
      
    - name: hands
      enabled: true
      
    - name: mouth
      enabled: true
      
    # Custom healthcare agents
    - name: hipaa-compliance
      enabled: true
      config: agents/hipaa-compliance/config.yml
      charter: agents/hipaa-compliance/charter.md
      
    - name: terminology
      enabled: true
      config: agents/terminology/config.yml
      charter: agents/terminology/charter.md
      
    - name: anonymizer
      enabled: true
      config: agents/anonymizer/config.yml
      charter: agents/anonymizer/charter.md
  
  quality_gates:
    pre_commit:
      - agent: hipaa-compliance
        check: phi_exposure_scan
        required: true
        
      - agent: anonymizer
        check: test_data_validation
        required: true
        
    code_review:
      - agent: eyes
        check: standard_review
        required: true
        
      - agent: hipaa-compliance
        check: compliance_review
        required: true
        blocking: true
        
      - agent: terminology
        check: terminology_validation
        required: true
        
    pre_merge:
      - agent: hipaa-compliance
        check: final_compliance_check
        required: true
        blocking: true
        
      - agent: brain
        check: architecture_approval
        required: true
        
  workflows:
    feature_development:
      stages:
        - name: planning
          agents: [brain, hipaa-compliance]
          
        - name: implementation
          agents: [hands, terminology]
          
        - name: review
          agents: [eyes, hipaa-compliance, terminology]
          
        - name: documentation
          agents: [mouth]
          
    hotfix:
      fast_track: true
      required_agents: [hands, eyes, hipaa-compliance]
      
  ceremonies:
    config: ceremonies/ceremonies.yml
```

### 6.2 Create Quality Gate Scripts

Create `medbook/.squad/scripts/quality-gates.sh`:

```bash
#!/bin/bash
# SQUAD Quality Gate Execution Script

set -e

echo "🔍 Running SQUAD Quality Gates..."

# Pre-commit gate: PHI exposure scan
echo "Stage: Pre-Commit"
echo "Agent: HIPAA Compliance"
copilot squad run --agent hipaa-compliance --check phi_exposure_scan

echo "Agent: Anonymizer"
copilot squad run --agent anonymizer --check test_data_validation

# Code review gate
echo "Stage: Code Review"
echo "Agent: Eyes (Standard Review)"
copilot squad run --agent eyes --check standard_review

echo "Agent: HIPAA Compliance (Compliance Review)"
copilot squad run --agent hipaa-compliance --check compliance_review --blocking

echo "Agent: Terminology (Validation)"
copilot squad run --agent terminology --check terminology_validation

# Pre-merge gate
echo "Stage: Pre-Merge"
echo "Agent: HIPAA Compliance (Final Check)"
copilot squad run --agent hipaa-compliance --check final_compliance_check --blocking

echo "Agent: Brain (Architecture Approval)"
copilot squad run --agent brain --check architecture_approval

echo "✅ All quality gates passed!"
```

### 6.3 Configure CI/CD Integration

Create `.github/workflows/squad-quality-gates.yml`:

```yaml
name: SQUAD Quality Gates

on:
  pull_request:
    branches: [ main, develop ]
    paths:
      - 'medbook/**'

jobs:
  hipaa-compliance:
    name: HIPAA Compliance Check
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Copilot
        uses: github/copilot-cli-action@v1
        
      - name: Run HIPAA Compliance Agent
        run: |
          cd medbook
          copilot squad run \
            --agent hipaa-compliance \
            --check compliance_review \
            --output compliance-report.json
            
      - name: Upload Compliance Report
        uses: actions/upload-artifact@v3
        with:
          name: hipaa-compliance-report
          path: medbook/compliance-report.json
          
      - name: Check for Violations
        run: |
          VIOLATIONS=$(jq '.violations | length' medbook/compliance-report.json)
          if [ "$VIOLATIONS" -gt 0 ]; then
            echo "❌ HIPAA compliance violations found: $VIOLATIONS"
            jq '.violations' medbook/compliance-report.json
            exit 1
          fi
          echo "✅ No HIPAA compliance violations"
          
  terminology-validation:
    name: Healthcare Terminology Validation
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Copilot
        uses: github/copilot-cli-action@v1
        
      - name: Run Terminology Agent
        run: |
          cd medbook
          copilot squad run \
            --agent terminology \
            --check terminology_validation \
            --output terminology-report.json
            
  test-data-anonymization:
    name: Validate Test Data Anonymization
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Copilot
        uses: github/copilot-cli-action@v1
        
      - name: Run Anonymizer Agent
        run: |
          cd medbook
          copilot squad run \
            --agent anonymizer \
            --check test_data_validation
```

## Step 7: Integration Testing

**Branch:** `step-5-integration-test`

**Objective:** Test the customized SQUAD team on a real feature.

### 7.1 Create Test Scenario

You'll implement a new feature: **"Add appointment cancellation with reason tracking"**

This feature touches:
- Patient service (PHI data)
- Appointment service (business logic)
- Database schema (new fields)
- Proto definitions (API contracts)

**Requirements:**
- Track cancellation reason (patient, provider, emergency)
- Log cancellation for audit trail
- Notify patient via email
- Update provider schedule

### 7.2 Run SQUAD with Custom Agents

```bash
cd medbook
copilot squad start

# SQUAD session:
> "I need to add appointment cancellation functionality with reason tracking.
> The system should:
> 1. Allow canceling appointments with a reason (patient, provider, emergency)
> 2. Create audit logs for compliance
> 3. Send patient notifications
> 4. Update provider schedules
> 
> Please plan and implement this feature with HIPAA compliance in mind."
```

**Expected Workflow:**

1. **Brain Agent:** Creates architecture plan
   - Identifies services to modify
   - Plans database schema changes
   - Designs API contracts

2. **HIPAA Compliance Agent:** Raises compliance concerns
   - Cancellation reason is PHI (requires encryption)
   - Audit logging must include all PHI access
   - Patient notification must use secure email
   - Ensure proper access controls

3. **Terminology Agent:** Validates naming
   - Suggests `cancellation_reason_code` instead of `cancel_reason`
   - Validates field naming standards
   - Checks proto message naming

4. **Hands Agent:** Implements the feature
   - Adds database migration
   - Updates service code
   - Implements audit logging

5. **Anonymizer Agent:** Creates test data
   - Generates fake patient cancellation scenarios
   - Creates anonymized test fixtures
   - Validates no real PHI in tests

6. **Eyes Agent:** Reviews implementation
   - Code quality check
   - Best practices validation

7. **HIPAA Compliance Agent:** Final review
   - Verifies encryption
   - Checks audit logging
   - Confirms no PHI leakage

8. **Mouth Agent:** Documents the feature
   - API documentation
   - Deployment notes
   - Compliance documentation

### 7.3 Validate Agent Interactions

Check that custom agents are triggered appropriately:

```bash
# View SQUAD execution log
cat .squad/logs/session-$(date +%Y%m%d).log

# Check which agents participated
grep "Agent:" .squad/logs/session-$(date +%Y%m%d).log | sort | uniq -c

# Review compliance findings
cat .squad/outputs/hipaa-compliance-report.json
```

**Expected output:**
```json
{
  "agent": "hipaa-compliance",
  "session_id": "...",
  "checks_performed": [
    {
      "check": "phi_exposure_scan",
      "status": "passed",
      "findings": []
    },
    {
      "check": "audit_logging_validation",
      "status": "passed",
      "findings": [
        "✅ Audit log created for cancellation"
      ]
    },
    {
      "check": "encryption_verification",
      "status": "passed",
      "findings": [
        "✅ cancellation_reason stored encrypted"
      ]
    }
  ],
  "violations": [],
  "approval": "granted"
}
```

## Step 8: Refine Agent Behavior

**Objective:** Improve agent effectiveness based on real usage.

### 8.1 Review Agent Performance

After running SQUAD on several features, evaluate:

**Questions to ask:**
- Did the HIPAA agent catch all compliance issues?
- Were there false positives?
- Did the terminology agent provide helpful suggestions?
- Was the anonymizer agent's test data realistic?
- Did ceremonies happen at the right time?

### 8.2 Update Agent Charters

Based on findings, refine agent behavior. For example:

**Before:**
```yaml
# hipaa-compliance/config.yml
triggers:
  - file_pattern: "internal/patient/**"
```

**After (more precise):**
```yaml
triggers:
  - file_pattern: "internal/patient/**"
  - file_pattern: "internal/appointment/**"
  - file_pattern: "**/*_test.go"  # Review test data too
  - keyword: "PHI"
  - keyword: "sensitive"
  - keyword: "encrypt"
  exclude:
    - "**/*.pb.go"  # Don't review generated protobuf code
```

### 8.3 Add New Compliance Rules

If the agent missed an issue, add it to the checklist:

```yaml
# Add to checklist.yml
phi_exposure:
  - rule: "No PHI in metrics or telemetry"
    severity: critical
    pattern: 'metrics\.Record.*patient\.(Name|SSN|DOB)'
    reason: "PHI must not be sent to observability platforms"
```

## Step 9: Document Custom Agents

**Objective:** Create templates and guides for other teams.

### 9.1 Create Custom Agent Template

Create `medbook/.squad/docs/custom-agent-template.md`:

```markdown
# Custom Agent Creation Template

## Agent Definition

### Charter (charter.md)
```markdown
# [Agent Name]

## Role
[One sentence description of the agent's purpose]

## Responsibilities
- [Responsibility 1]
- [Responsibility 2]

## When to Trigger
- [Trigger condition 1]
- [Trigger condition 2]

## Review Checklist
- [ ] [Check 1]
- [ ] [Check 2]
```

### Configuration (config.yml)
```yaml
agent:
  name: [agent-name]
  role: [Agent Role]
  description: [Description]
  
  triggers:
    - file_pattern: "[pattern]"
    - keyword: "[keyword]"
  
  capabilities:
    - [capability_1]
    - [capability_2]
  
  review_priority: [low|medium|high|critical]
  
  integration_points:
    - stage: [pre_commit|code_review|pre_merge]
      action: [action_name]
```

## Best Practices

### When to Create Custom Agents

✅ **Good use cases:**
- Industry-specific compliance (HIPAA, PCI-DSS, SOC2)
- Domain-specific terminology or standards
- Company-specific quality requirements
- Technology stack-specific checks (e.g., Kubernetes manifest validation)

❌ **Not recommended:**
- Duplicating existing agent capabilities
- Overly narrow single-purpose checks (use linters instead)
- Personal preferences or style choices

### Agent Design Principles

1. **Single Responsibility:** Each agent should have a clear, focused purpose
2. **Composability:** Agents should work well together, not overlap
3. **Actionable Feedback:** Agent output should guide developers to solutions
4. **Domain Expertise:** Encode real domain knowledge, not just rules
5. **Evolvability:** Easy to update as requirements change

### Integration Guidelines

- **Pre-commit:** Fast checks (syntax, basic validation)
- **Code Review:** Deeper analysis (compliance, terminology, logic)
- **Pre-merge:** Final gates (approval, comprehensive checks)

## Verification

### Checklist

Confirm your SQUAD customization is complete:

- [ ] All custom agents have charters and configurations
- [ ] Quality gates are configured in team.yml
- [ ] Ceremonies are defined and scheduled
- [ ] CI/CD integration working
- [ ] Custom agents tested on real features
- [ ] Agent behavior refined based on usage
- [ ] Documentation created for other teams

### Testing Custom Agents

```bash
# Test individual agent
copilot squad run --agent hipaa-compliance --dry-run

# Test full workflow
copilot squad start --feature "test-feature"

# Validate configuration
copilot squad validate --config .squad/team.yml
```

## Common Issues and Solutions

### Issue: Custom agent not triggering

**Solution:** Check trigger patterns in config.yml
```yaml
triggers:
  - file_pattern: "internal/patient/**/*.go"  # Be specific
  - keyword: "PHI"  # Add relevant keywords
```

### Issue: Too many false positives

**Solution:** Add exclusions and refine patterns
```yaml
triggers:
  - file_pattern: "internal/patient/**"
  exclude:
    - "**/*_test.go"  # Exclude tests
    - "**/*.pb.go"    # Exclude generated code
```

### Issue: Agent feedback not actionable

**Solution:** Improve checklist with specific guidance
```yaml
- rule: "No PHI in logs"
  severity: critical
  pattern: 'log\..*patient\.(Name|SSN)'
  suggestion: "Use patient ID instead: log.Info(\"Patient action\", \"patient_id\", patient.ID)"
```

## Additional Resources

- [SQUAD Framework Documentation](https://github.com/github/squad)
- [HIPAA Compliance for Developers](https://www.hhs.gov/hipaa/for-professionals/security/index.html)
- [Healthcare Terminology Standards](https://www.icd10data.com/)
- [Go Best Practices](https://go.dev/doc/effective_go)

## Next Steps

After completing this lab:

1. **Apply to Your Codebase:** Identify domain-specific requirements in your projects
2. **Create Your Agents:** Define custom agents for your compliance/quality needs
3. **Share with Team:** Document and share your SQUAD customizations
4. **Iterate:** Continuously refine agent behavior based on real usage
5. **Expand:** Add more agents as new requirements emerge

## Conclusion

You've successfully customized a SQUAD team with domain-specific agents for healthcare development. You learned how to:

- Define custom agent charters with specialized knowledge
- Configure quality gates for compliance and terminology validation
- Set up ceremonies for healthcare-specific workflows
- Integrate custom agents into CI/CD pipelines
- Test and refine agent behavior

SQUAD's extensibility allows you to encode your organization's unique development culture, compliance requirements, and quality standards directly into your development workflow. This reduces cognitive load on developers while ensuring critical requirements are consistently met.

---

**Lab Complete!** 🎉

For questions or feedback, please open an issue in the repository.
