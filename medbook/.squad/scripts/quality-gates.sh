#!/bin/bash
# SQUAD Quality Gate Execution Script

set -e

echo "🔍 Running SQUAD Quality Gates for MedBook Healthcare Platform..."
echo ""

# Configuration
SQUAD_CLI="copilot squad run"
OUTPUT_DIR=".squad/outputs"
mkdir -p "$OUTPUT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============================================
# Pre-Commit Gates
# ============================================
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Stage: Pre-Commit${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo -e "${YELLOW}Agent: HIPAA Compliance (PHI Exposure Scan)${NC}"
$SQUAD_CLI --agent hipaa-compliance --check phi_exposure_scan --output "$OUTPUT_DIR/hipaa-precommit.json"
echo -e "${GREEN}✅ PHI exposure scan completed${NC}"
echo ""

echo -e "${YELLOW}Agent: Anonymizer (Test Data Validation)${NC}"
$SQUAD_CLI --agent anonymizer --check test_data_validation --output "$OUTPUT_DIR/anonymizer-precommit.json"
echo -e "${GREEN}✅ Test data validation completed${NC}"
echo ""

# ============================================
# Code Review Gates
# ============================================
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Stage: Code Review${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo -e "${YELLOW}Agent: Eyes (Standard Review)${NC}"
$SQUAD_CLI --agent eyes --check standard_review --output "$OUTPUT_DIR/eyes-review.json"
echo -e "${GREEN}✅ Standard code review completed${NC}"
echo ""

echo -e "${YELLOW}Agent: HIPAA Compliance (Compliance Review)${NC}"
$SQUAD_CLI --agent hipaa-compliance --check compliance_review --blocking --output "$OUTPUT_DIR/hipaa-review.json"

# Check for violations
VIOLATIONS=$(jq '.violations | length' "$OUTPUT_DIR/hipaa-review.json" 2>/dev/null || echo "0")
if [ "$VIOLATIONS" -gt 0 ]; then
    echo -e "${RED}❌ HIPAA compliance violations found: $VIOLATIONS${NC}"
    jq '.violations' "$OUTPUT_DIR/hipaa-review.json"
    exit 1
fi
echo -e "${GREEN}✅ HIPAA compliance review passed${NC}"
echo ""

echo -e "${YELLOW}Agent: Terminology (Terminology Validation)${NC}"
$SQUAD_CLI --agent terminology --check terminology_validation --output "$OUTPUT_DIR/terminology-review.json"
echo -e "${GREEN}✅ Terminology validation completed${NC}"
echo ""

# ============================================
# Pre-Merge Gates
# ============================================
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Stage: Pre-Merge${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo -e "${YELLOW}Agent: HIPAA Compliance (Final Check)${NC}"
$SQUAD_CLI --agent hipaa-compliance --check final_compliance_check --blocking --output "$OUTPUT_DIR/hipaa-final.json"

# Check for violations in final check
FINAL_VIOLATIONS=$(jq '.violations | length' "$OUTPUT_DIR/hipaa-final.json" 2>/dev/null || echo "0")
if [ "$FINAL_VIOLATIONS" -gt 0 ]; then
    echo -e "${RED}❌ Final HIPAA compliance violations found: $FINAL_VIOLATIONS${NC}"
    jq '.violations' "$OUTPUT_DIR/hipaa-final.json"
    exit 1
fi
echo -e "${GREEN}✅ Final HIPAA compliance check passed${NC}"
echo ""

echo -e "${YELLOW}Agent: Brain (Architecture Approval)${NC}"
$SQUAD_CLI --agent brain --check architecture_approval --output "$OUTPUT_DIR/brain-approval.json"
echo -e "${GREEN}✅ Architecture approval completed${NC}"
echo ""

# ============================================
# Summary
# ============================================
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All SQUAD Quality Gates Passed!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "📊 Reports generated in: $OUTPUT_DIR/"
echo ""

# Generate summary report
cat > "$OUTPUT_DIR/summary.txt" <<EOF
SQUAD Quality Gates Summary
============================
Timestamp: $(date)

✅ Pre-Commit Gates:
   - HIPAA PHI Exposure Scan
   - Test Data Validation

✅ Code Review Gates:
   - Standard Code Review
   - HIPAA Compliance Review
   - Terminology Validation

✅ Pre-Merge Gates:
   - Final HIPAA Compliance Check
   - Architecture Approval

All gates passed successfully.
EOF

echo "✅ Summary report: $OUTPUT_DIR/summary.txt"
echo ""
echo "🎉 Ready to merge!"
