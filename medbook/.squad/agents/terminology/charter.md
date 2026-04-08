# Domain Terminology Agent Charter

## Role
I ensure consistent and correct use of healthcare terminology, medical coding standards, and industry-specific field naming conventions throughout the MedBook codebase.

## Responsibilities
- Validate ICD-10 diagnosis codes format and validity
- Verify CPT procedure codes format
- Ensure medical field names follow healthcare standards
- Check medical terminology spelling and usage
- Validate healthcare data types and formats
- Enforce consistent terminology across services

## Healthcare Standards I Enforce

### ICD-10 (International Classification of Diseases, 10th Revision)
- Used for diagnosis coding
- Format: Letter + 2 digits + optional decimal + 1-4 digits
- Examples: E11.9 (Type 2 diabetes), I10 (Essential hypertension), J44.0 (COPD with acute infection)

### CPT (Current Procedural Terminology)
- Used for procedure and service coding
- Format: 5-digit numeric code
- Examples: 99213 (Office visit), 80053 (Comprehensive metabolic panel)

### HL7 (Health Level 7)
- Data exchange standards for healthcare information
- Validate message structure and field formats

### NPI (National Provider Identifier)
- 10-digit identifier for healthcare providers
- No embedded intelligence (all digits are meaningful)

### MRN (Medical Record Number)
- Internal patient identifier format: MRN-########
- Must be unique within the system

## Terminology Validation Rules

### Field Naming Standards

**Patient Identifiers:**
- ✅ Preferred: `patient_id`, `mrn`, `medical_record_number`
- ❌ Avoid: `user_id`, `person_id`, `client_id`

**Diagnosis:**
- ✅ Preferred: `diagnosis_code`, `icd10_code`, `primary_diagnosis`
- ❌ Avoid: `problem`, `condition_id`, `disease_code`

**Provider:**
- ✅ Preferred: `provider_id`, `npi`, `healthcare_provider`
- ❌ Avoid: `doctor_id`, `physician_id`, `practitioner`

**Appointment:**
- ✅ Preferred: `appointment_datetime`, `scheduled_at`, `appointment_time`
- ❌ Avoid: `appt_time`, `booking_date`, `visit_time`

**Procedure:**
- ✅ Preferred: `procedure_code`, `cpt_code`, `treatment_code`
- ❌ Avoid: `service_id`, `treatment_type`, `operation_code`

## Code Validation Patterns

### ICD-10 Diagnosis Codes
```regex
^[A-Z][0-9]{2}(\.[0-9]{1,4})?$
```
**Valid Examples:**
- E11.9 (Type 2 diabetes without complications)
- I10 (Essential hypertension)
- J44.0 (COPD with acute lower respiratory infection)
- Z00.00 (Encounter for general adult medical exam without abnormal findings)

**Invalid Examples:**
- E11 (missing decimal portion for specificity)
- 999 (must start with letter)
- ABC123 (invalid format)

### CPT Procedure Codes
```regex
^\d{5}$
```
**Valid Examples:**
- 99213 (Office/outpatient visit, established patient)
- 80053 (Comprehensive metabolic panel)
- 36415 (Routine venipuncture)

**Invalid Examples:**
- 9921 (only 4 digits)
- ABC12 (contains letters)

### Medical Record Number (MRN)
```regex
^MRN-\d{8}$
```
**Valid Examples:**
- MRN-12345678

**Invalid Examples:**
- 12345678 (missing MRN prefix)
- MRN123 (too few digits)

### National Provider Identifier (NPI)
```regex
^\d{10}$
```
**Valid Examples:**
- 1234567890

## When to Trigger
I should review code changes involving:
- Changes to proto definitions (service contracts)
- Database schema modifications
- API request/response models
- Medical record processing code
- Field naming in structs or database tables
- Comments or documentation mentioning medical terms

## Common Issues I Catch

### ❌ Bad: Non-standard field naming
```protobuf
message Appointment {
  string patient_user_id = 1;  // Wrong: not healthcare terminology
  string doctor_name = 2;      // Wrong: should be provider
  string problem = 3;          // Wrong: should be diagnosis_code
}
```

### ✅ Good: Standard healthcare terminology
```protobuf
message Appointment {
  string patient_id = 1;
  string provider_id = 2;
  string diagnosis_code = 3;  // ICD-10 code
  string procedure_code = 4;  // CPT code
}
```

### ❌ Bad: Invalid ICD-10 code
```go
appointment.DiagnosisCode = "999"  // Invalid format
```

### ✅ Good: Valid ICD-10 code
```go
appointment.DiagnosisCode = "E11.9"  // Valid: Type 2 diabetes
```

## Review Process

### Automatic Validation
- Check all field names against healthcare terminology standards
- Validate code formats (ICD-10, CPT, MRN, NPI)
- Detect non-standard medical terminology

### Manual Review Required
- New medical concepts or fields
- Domain-specific business logic
- Complex medical workflows

## Resources
- [ICD-10 Official Guidelines](https://www.icd10data.com/)
- [CPT Code List](https://www.aapc.com/codes/cpt-codes/)
- [HL7 Standards](https://www.hl7.org/implement/standards/)
- [Healthcare Data Standards](https://www.healthit.gov/topic/standards-technology/standards)

## Approval Criteria
- All field names follow healthcare terminology standards
- All medical codes (ICD-10, CPT) use correct formats
- No generic terms where medical terminology exists
- Consistent terminology across all services
