# MedBook Healthcare Platform

A HIPAA-compliant healthcare appointment scheduling platform built with Go microservices.

## Architecture

MedBook consists of three core microservices:
- **Patient Service** (port 8081): Patient registration and management
- **Appointment Service** (port 8082): Appointment scheduling and management  
- **Provider Service** (port 8083): Healthcare provider management

## Getting Started

### Prerequisites
- Go 1.22+
- PostgreSQL 14+
- Docker (optional, for containerized deployment)

### Local Development

1. **Start PostgreSQL:**
   ```bash
   docker run -d \
     --name medbook-postgres \
     -e POSTGRES_DB=medbook \
     -e POSTGRES_USER=medbook \
     -e POSTGRES_PASSWORD=medbook \
     -p 5432:5432 \
     postgres:14
   ```

2. **Run Database Migrations:**
   ```bash
   # TODO: Add migration scripts
   ```

3. **Start Services:**
   ```bash
   # Patient Service
   cd cmd/patient-service
   go run main.go

   # Appointment Service (in another terminal)
   cd cmd/appointment-service
   go run main.go

   # Provider Service (in another terminal)
   cd cmd/provider-service
   go run main.go
   ```

## SQUAD Customization

This project uses a customized SQUAD team with healthcare-specific agents:

### Custom Agents

- **🏥 HIPAA Compliance Agent**: Ensures PHI protection and regulatory compliance
- **📚 Terminology Agent**: Validates healthcare terminology (ICD-10, CPT codes)
- **🔐 Anonymizer Agent**: Ensures test data uses anonymized patient information

### Running SQUAD Quality Gates

```bash
# Run all quality gates
.squad/scripts/quality-gates.sh

# Run individual agent
copilot squad run --agent hipaa-compliance --check compliance_review
```

## Healthcare Standards

### ICD-10 Diagnosis Codes
Format: Letter + 2 digits + optional decimal + 1-4 digits
- Example: `E11.9` (Type 2 diabetes)

### CPT Procedure Codes
Format: 5-digit numeric code
- Example: `99213` (Office visit)

### Patient Identifiers
- **MRN**: `MRN-########` (production) or `MRN-T#######` (test data)
- **NPI**: 10-digit National Provider Identifier

## HIPAA Compliance

All code changes are reviewed for HIPAA compliance:
- PHI encryption at rest and in transit
- Audit logging for PHI access
- Role-based access control
- Minimum necessary data principle

## Testing

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...
```

### Test Data

All test data must be anonymized:
- Patient names start with "Test "
- MRNs use "MRN-T" prefix
- SSNs use "000-00-" prefix
- Phone numbers in (555) 01XX range
- Email addresses use @example.com

## API Documentation

### Patient Service (Port 8081)
- `POST /patients` - Create patient
- `GET /patients/:id` - Get patient by ID
- `PUT /patients/:id` - Update patient
- `GET /patients` - List patients

### Appointment Service (Port 8082)
- `POST /appointments` - Create appointment
- `GET /appointments/:id` - Get appointment
- `PUT /appointments/:id/cancel` - Cancel appointment
- `GET /appointments/patient/:patient_id` - List patient appointments

### Provider Service (Port 8083)
- `POST /providers` - Create provider
- `GET /providers/:id` - Get provider
- `GET /providers` - List providers

## Deployment

### Docker Build
```bash
docker build -t medbook/patient-service:latest -f cmd/patient-service/Dockerfile .
docker build -t medbook/appointment-service:latest -f cmd/appointment-service/Dockerfile .
docker build -t medbook/provider-service:latest -f cmd/provider-service/Dockerfile .
```

### Kubernetes Deployment
```bash
kubectl apply -f k8s/
```

## License

MIT
