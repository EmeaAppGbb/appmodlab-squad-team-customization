---
name: "hipaa-compliance"
description: "Detect HIPAA violations in Go source: PHI leaks, missing audit logs, unprotected endpoints, encryption gaps, hardcoded patient data"
domain: "security, compliance, healthcare"
confidence: "high"
source: "earned (hipaa-auditor charter, hipaa-compliance checklist, medbook codebase patterns)"
---

## Context

MedBook is a Go healthcare microservices platform (patient-service, appointment-service, provider-service) built with Gin and PostgreSQL. All code touching Protected Health Information must comply with the HIPAA Privacy Rule (§164.500–534), Security Rule (§164.302–318), and Breach Notification Rule (§164.400–414). This skill codifies machine-checkable patterns to catch violations before they reach production.

**PHI fields under protection:**

| Category | Go Fields / DB Columns |
|---|---|
| Patient identifiers | `Name`, `FirstName`, `LastName`, `SSN`, `MRN`, `DOB`, `DateOfBirth` |
| Contact info | `Address`, `Phone`, `Email`, `ZipCode` |
| Medical records | `DiagnosisCode`, `ICD10Code`, `CPTCode`, `TreatmentNotes`, `LabResults` |
| Insurance/billing | `PolicyNumber`, `ClaimID`, `InsuranceID` |

## Patterns

### 1. PHI Exposure in Log Statements

**Rule:** Never log PHI fields. Use opaque identifiers (patient ID, request ID) instead.

**Scan targets** — any `.go` file under `internal/` or `cmd/`:
```
log.Println(...)     log.Printf(...)      log.Fatalf(...)
log.Info(...)        log.Debug(...)        log.Error(...)        log.Warn(...)
fmt.Printf(...)      fmt.Println(...)      fmt.Sprintf(...)  (when used for logging)
slog.Info(...)       slog.Debug(...)       slog.Error(...)       slog.Warn(...)
zap.String(...)      zap.Any(...)         (when field is PHI)
```

**Violation regex patterns:**
```
log\.(Println|Printf|Fatalf|Info|Debug|Error|Warn)\(.*\b(patient|p)\.(Name|FirstName|LastName|SSN|DOB|MRN|Address|Phone|Email)
fmt\.(Printf|Println|Sprintf)\(.*\b(patient|p)\.(Name|FirstName|LastName|SSN|DOB|MRN|Address|Phone|Email)
slog\.(Info|Debug|Error|Warn)\(.*\b(patient|p)\.(Name|FirstName|LastName|SSN|DOB|MRN)
```

**Also flag** string interpolation that embeds PHI:
```go
// ❌ VIOLATION
fmt.Sprintf("Processing patient %s %s", patient.FirstName, patient.LastName)
log.Printf("SSN lookup: %s", patient.SSN)

// ✅ CORRECT
log.Printf("Processing patient_id=%s", patient.ID)
slog.Info("Patient created", "patient_id", patient.ID, "request_id", reqID)
```

### 2. PHI Exposure in HTTP Responses and Error Messages

**Rule:** Error responses must never contain PHI. Return opaque identifiers and generic messages.

**Scan targets** — Gin handler files (`internal/*/service.go`, `cmd/*/main.go`):
```
c\.JSON\(.*patient\.(Name|FirstName|LastName|SSN|DOB|MRN|Address|Phone|Email)
gin\.H\{.*"error".*patient\.(Name|FirstName|LastName|SSN|DOB)
fmt\.Errorf\(.*patient\.(Name|FirstName|LastName|SSN|DOB|MRN)
errors\.New\(.*patient\.(Name|FirstName|LastName|SSN|DOB|MRN)
status\.Errorf\(.*patient\.(Name|FirstName|LastName|SSN|DOB|MRN)
```

**Verify SSN omission in JSON responses:**
- Struct fields holding SSN must have `json:"ssn,omitempty"` or `json:"-"` tag
- `GetPatient` and `ListPatients` handlers must not include SSN in SELECT queries
- Response DTOs should be purpose-specific, not reuse the full entity model

```go
// ❌ VIOLATION — PHI in error response
c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Patient %s not found", patient.Name)})

// ✅ CORRECT — opaque ID only
c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found", "patient_id": id})
```

**Verify minimum necessary data in responses:**
- `ListPatients` must return summary projections, not full records
- No `SELECT *` on tables containing PHI — use explicit column lists
- Pagination must be enforced on multi-record PHI endpoints (`page_size`, `page_token`)

### 3. Audit Logging Gaps on PHI Access

**Rule:** Every handler that reads, writes, updates, or deletes PHI must emit an audit log entry before returning. Per §164.312(b), audit logs must be tamper-evident.

**Required audit log fields:**
```go
type AuditEntry struct {
    UserID       string    // Authenticated identity
    Timestamp    time.Time // UTC
    Action       string    // CREATE, READ, UPDATE, DELETE
    ResourceType string    // "patient", "appointment", "provider"
    ResourceID   string    // ID of accessed resource
    Outcome      string    // "success" or "failure"
    IPAddress    string    // Source IP from request
}
```

**Handlers that MUST have audit logging:**

| Service | Handler | Action |
|---|---|---|
| patient-service | `CreatePatient` | CREATE |
| patient-service | `GetPatient` | READ |
| patient-service | `UpdatePatient` | UPDATE |
| patient-service | `ListPatients` | READ |
| appointment-service | `CreateAppointment` | CREATE |
| appointment-service | `GetAppointment` | READ |
| appointment-service | `CancelAppointment` | UPDATE |
| appointment-service | `ListPatientAppointments` | READ |
| provider-service | `CreateProvider` | CREATE |
| provider-service | `GetProvider` | READ |
| provider-service | `ListProviders` | READ |

**Detection:** For each handler above, verify it contains an `audit.Log` or `auditLog` call. Flag any handler accessing PHI without a corresponding audit entry.

```go
// ❌ VIOLATION — no audit log
func (s *Service) GetPatient(c *gin.Context) {
    id := c.Param("id")
    patient, err := s.db.QueryRow(...)
    c.JSON(http.StatusOK, patient)
}

// ✅ CORRECT — audit log before response
func (s *Service) GetPatient(c *gin.Context) {
    id := c.Param("id")
    patient, err := s.db.QueryRow(...)
    audit.Log(c, audit.Entry{
        Action: "READ", ResourceType: "patient", ResourceID: id,
        Outcome: "success", UserID: auth.UserID(c),
    })
    c.JSON(http.StatusOK, patient)
}
```

**Also verify:**
- Failed access attempts (auth failures, not-found on PHI) produce audit entries with `Outcome: "failure"`
- Audit logs themselves do NOT contain PHI — only resource IDs
- Audit log entries are not written via `log.Println` (must use dedicated audit package/sink)

### 4. Access Control Middleware on Patient-Facing Endpoints

**Rule:** Every PHI endpoint must have authentication and authorization middleware. Per §164.312(a)(1) and §164.312(d), access to PHI requires unique user identification and verification.

**Scan targets** — router setup in `cmd/*/main.go`:

**Check that PHI routes use auth middleware:**
```go
// ❌ VIOLATION — no auth middleware
router.POST("/patients", svc.CreatePatient)
router.GET("/patients/:id", svc.GetPatient)
router.GET("/appointments/patient/:patient_id", svc.ListPatientAppointments)

// ✅ CORRECT — auth + role middleware applied
authorized := router.Group("/")
authorized.Use(auth.RequireToken(), auth.RequireRole("physician", "nurse"))
{
    authorized.POST("/patients", svc.CreatePatient)
    authorized.GET("/patients/:id", svc.GetPatient)
}
```

**Specific checks:**
1. **Authentication present:** Routes for `/patients`, `/appointments`, `/providers` must have `auth.RequireToken()` or equivalent middleware in the chain
2. **Authorization present:** PHI write operations (POST, PUT, DELETE) must have role-based middleware (`auth.RequireRole(...)`)
3. **No public PHI routes:** The only unauthenticated endpoint should be `/health`
4. **Gin middleware chain:** Verify `router.Use(...)` or `group.Use(...)` includes auth before PHI handlers
5. **Session expiration:** Token validation must enforce expiration (max 8 hours for PHI access)

**Detection pattern — flag any direct route registration without middleware for PHI paths:**
```
router\.(GET|POST|PUT|DELETE|PATCH)\("/(patients|appointments|providers)
```
If this pattern matches without a preceding `Use(auth.*)` in the same route group, it is a violation.

### 5. Encryption Configuration for PHI at Rest and in Transit

**Rule:** All PHI must be encrypted in transit (TLS 1.2+) and at rest per §164.312(a)(2)(iv) and §164.312(e)(1).

#### In Transit

**Scan all `cmd/*/main.go` and any TLS configuration files:**

```go
// ❌ VIOLATION — no TLS config, using default HTTP
srv := &http.Server{Addr: ":8081", Handler: router}
srv.ListenAndServe()

// ✅ CORRECT — TLS enforced with minimum version
tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
srv := &http.Server{Addr: ":8081", Handler: router, TLSConfig: tlsConfig}
srv.ListenAndServeTLS(certFile, keyFile)
```

**Flag these violations:**
- `grpc.WithInsecure()` or `grpc.WithTransportCredentials(insecure.NewCredentials())` — gRPC without TLS
- `ListenAndServe()` without TLS on PHI-serving ports (8081, 8082, 8083)
- Missing `tls.Config{MinVersion: tls.VersionTLS12}` on any HTTP server serving PHI
- Weak cipher suites (RC4, DES, 3DES, NULL ciphers)

#### At Rest

**Database connection — verify SSL:**
```go
// ❌ VIOLATION
"postgres://medbook:medbook@localhost:5432/medbook"                    // no sslmode
"postgres://medbook:medbook@localhost:5432/medbook?sslmode=disable"    // SSL disabled

// ✅ CORRECT
"postgres://medbook:medbook@localhost:5432/medbook?sslmode=require"
"postgres://medbook:medbook@localhost:5432/medbook?sslmode=verify-full"
```

**Kubernetes manifests — scan `k8s/*.yaml`:**
- Verify database credentials use `secretKeyRef`, not plaintext `value:`
- Flag any PHI-adjacent values in ConfigMaps (use Secrets instead)
- Secrets containing PHI should reference encrypted-at-rest storage (`EncryptionConfiguration`)

**Encryption key management:**
- No hardcoded encryption keys, certificates, or TLS material in `.go` files
- Keys must come from environment variables, Kubernetes Secrets, or a secrets manager
- Flag patterns: `[]byte("hardcoded-key")`, `privateKey := "..."`, embedded PEM blocks

### 6. Hardcoded Patient Identifiers in Test Files

**Rule:** Test files must use clearly synthetic patient data following the project's anonymization conventions. Real-looking patient data in tests creates compliance risk and may constitute actual PHI if copied from production.

**Scan targets** — `tests/**/*_test.go` and `internal/**/*_test.go`:

**Required anonymization prefixes (per project convention):**

| Field | Required Pattern | Example |
|---|---|---|
| Patient name | `"Test "` prefix | `"Test Alice"`, `"Test Bob"` |
| MRN | `"MRN-T"` prefix | `"MRN-T001"`, `"MRN-T002"` |
| SSN | `"000-00-"` prefix | `"000-00-1234"` |
| Phone | `"(555) 01"` prefix | `"(555) 0100"` |
| Email | `@example.com` domain | `"test.alice@example.com"` |

**Flag these violations:**
```go
// ❌ VIOLATION — realistic SSN
patient := Patient{SSN: "123-45-6789", Name: "John Smith"}

// ❌ VIOLATION — real-looking MRN
mrn := "MRN-20240001"

// ❌ VIOLATION — real email domain
email := "patient@hospital.org"

// ✅ CORRECT — anonymized test data
patient := Patient{SSN: "000-00-1234", Name: "Test Alice"}
mrn := "MRN-T001"
email := "test.alice@example.com"
```

**Detection patterns for test files:**
```
SSN.*"(?!000-00-)\d{3}-\d{2}-\d{4}"          # SSN not using 000-00- prefix
MRN.*"MRN-(?!T)\w+"                           # MRN not using MRN-T prefix
Name.*"(?!Test )[A-Z][a-z]+ [A-Z][a-z]+"      # Name not using Test prefix
@(?!example\.com)[a-z]+\.[a-z]{2,}            # Email not using @example.com
"\(\d{3}\) (?!01)\d{2}"                       # Phone not using (555) 01XX range
```

**Also flag:**
- Hardcoded patient IDs that look like production UUIDs (non-zero-prefixed)
- Copy-pasted JSON fixtures containing realistic patient data
- Database seed scripts with non-anonymized records

## Examples

### ✓ Correct: Complete PHI-safe handler

```go
func (s *Service) GetPatient(c *gin.Context) {
    id := c.Param("id")
    userID := auth.UserID(c) // from auth middleware

    row := s.pool.QueryRow(context.Background(),
        "SELECT id, first_name, last_name, date_of_birth FROM patients WHERE id = $1", id)

    var p Patient
    if err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.DOB); err != nil {
        audit.Log(c, audit.Entry{
            Action: "READ", ResourceType: "patient", ResourceID: id,
            Outcome: "failure", UserID: userID,
        })
        slog.Error("Patient lookup failed", "patient_id", id, "error", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
        return
    }

    audit.Log(c, audit.Entry{
        Action: "READ", ResourceType: "patient", ResourceID: id,
        Outcome: "success", UserID: userID,
    })
    c.JSON(http.StatusOK, p)
}
```

### ✗ Incorrect: Multiple HIPAA violations in one handler

```go
func (s *Service) GetPatient(c *gin.Context) {
    id := c.Param("id")
    // ❌ No auth middleware on route
    // ❌ No audit logging

    row := s.pool.QueryRow(context.Background(),
        "SELECT * FROM patients WHERE id = $1", id) // ❌ SELECT * over-fetches PHI

    var p Patient
    if err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.SSN, &p.DOB); err != nil {
        log.Printf("Failed to find patient %s %s", p.FirstName, p.LastName) // ❌ PHI in log
        c.JSON(http.StatusNotFound, gin.H{
            "error": fmt.Sprintf("Patient %s not found", p.FirstName), // ❌ PHI in response
        })
        return
    }
    c.JSON(http.StatusOK, p) // ❌ Returns SSN (no omitempty or DTO)
}
```

### ✓ Correct: Anonymized test data

```go
func TestCreatePatient(t *testing.T) {
    patient := Patient{
        FirstName: "Test Alice",
        LastName:  "Test Smith",
        SSN:       "000-00-1234",
        MRN:       "MRN-T001",
        Email:     "test.alice@example.com",
        Phone:     "(555) 0100",
    }
    // ... test logic
}
```

## Anti-Patterns

- ❌ Logging PHI "temporarily for debugging" — use patient IDs, never PHI fields
- ❌ Using `SELECT *` on PHI tables — explicitly list non-sensitive columns
- ❌ Returning full entity models from API endpoints — use purpose-specific DTOs
- ❌ Adding PHI endpoints without auth middleware ("I'll add auth later")
- ❌ Using `ListenAndServe()` without TLS on PHI-serving ports
- ❌ Hardcoding realistic patient data in tests ("it's just test data")
- ❌ Writing audit entries via `log.Println` — use a dedicated, tamper-evident audit sink
- ❌ Assuming `gin.Default()` logger middleware is sufficient for HIPAA audit logging
- ❌ Putting database credentials in Kubernetes ConfigMaps instead of Secrets
- ❌ Using `grpc.WithInsecure()` for service-to-service calls carrying PHI
- ❌ Skipping audit logs on error paths — failed access attempts must be logged too
- ❌ Storing encryption keys alongside the encrypted data in the same repository
