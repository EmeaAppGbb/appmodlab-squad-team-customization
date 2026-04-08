package patient_test

import (
	"testing"
	"time"

	"github.com/emeaappgbb/medbook/internal/patient"
)

// TestPatient_AnonymizedData validates that we're using properly anonymized test data
func TestPatient_AnonymizedData(t *testing.T) {
	testPatient := patient.Patient{
		MRN:       "MRN-T1234567", // Test MRN with T marker
		FirstName: "Test John",     // Name starts with "Test "
		LastName:  "Smith",
		DOB:       time.Date(1975, 6, 15, 0, 0, 0, 0, time.UTC),
		SSN:       "000-00-1234", // Clearly fake SSN
		Email:     "patient.test1234@example.com",
		Phone:     "(555) 0145",
		Address:   "123 Test Street, Anytown, NY 10001",
	}

	// Validate test data follows anonymization rules
	if testPatient.MRN[:6] != "MRN-T" {
		t.Errorf("Test MRN should have 'T' marker, got: %s", testPatient.MRN)
	}

	if testPatient.FirstName[:5] != "Test " {
		t.Errorf("Test patient name should start with 'Test ', got: %s", testPatient.FirstName)
	}

	if testPatient.SSN[:6] != "000-00" {
		t.Errorf("Test SSN should use 000-00- prefix, got: %s", testPatient.SSN)
	}

	if testPatient.Phone[:6] != "(555) " {
		t.Errorf("Test phone should use (555) prefix, got: %s", testPatient.Phone)
	}

	t.Logf("✅ Test data is properly anonymized")
}

// Example of what NOT to do (would be caught by Anonymizer Agent)
// func TestPatient_BadExample(t *testing.T) {
// 	badPatient := patient.Patient{
// 		MRN:       "MRN-12345678",        // ❌ No T marker
// 		FirstName: "John",                 // ❌ No Test prefix
// 		LastName:  "Doe",
// 		SSN:       "123-45-6789",          // ❌ Could be real SSN
// 		Email:     "john@realdomain.com",  // ❌ Real domain
// 		Phone:     "(555) 1234",           // ❌ Not in test range
// 	}
// }
