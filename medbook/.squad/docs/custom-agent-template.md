# Custom Agent Creation Template

Use this template to create new custom agents for your SQUAD team.

## Agent Definition Structure

```
.squad/agents/{agent-name}/
├── charter.md           # Agent's role, responsibilities, and guidelines
├── config.yml          # Agent configuration and triggers
├── checklist.yml       # Optional: Validation checklist
└── utilities.yml       # Optional: Utilities and templates
```

---

## 1. Charter (charter.md)

The charter defines the agent's role, responsibilities, and guidelines.

```markdown
# {Agent Name} Agent Charter

## Role
I am the {Agent Name} agent, responsible for {primary responsibility}.

## Responsibilities
- {Responsibility 1}
- {Responsibility 2}
- {Responsibility 3}

## {Domain-Specific Section}
Describe domain-specific knowledge or standards this agent enforces.

### Key Standards/Rules
- {Standard 1}
- {Standard 2}

## When to Trigger
I should be consulted when code changes involve:
- {Trigger condition 1}
- {Trigger condition 2}

## Review Checklist
- [ ] {Check 1}
- [ ] {Check 2}
- [ ] {Check 3}

## Example Violations I Catch

### ❌ Bad: {Description}
```{language}
{bad_example_code}
```

### ✅ Good: {Description}
```{language}
{good_example_code}
```

## Approval Process
- **Low-risk changes:** {approval_process}
- **Medium-risk changes:** {approval_process}
- **High-risk changes:** {approval_process}
```

---

## 2. Configuration (config.yml)

The configuration file defines agent settings, triggers, and capabilities.

```yaml
agent:
  name: {agent-name}
  role: {Agent Role Title}
  description: {Brief description of agent purpose}
  
  triggers:
    # File patterns that should trigger this agent
    - file_pattern: "path/to/**/*.ext"
    - file_pattern: "another/path/**"
    
    # Keywords that should trigger this agent
    - keyword: "important_term"
    - keyword: "another_keyword"
    
    # Exclusions (files to ignore)
    exclude:
      - "**/*.generated.*"
      - "**/vendor/**"
  
  capabilities:
    - {capability_1}
    - {capability_2}
    - {capability_3}
  
  review_priority: {low|medium|high|critical}
  
  integration_points:
    - stage: {pre_commit|code_review|pre_merge}
      action: {action_name}
      timeout: {seconds}
      blocking: {true|false}  # Optional
      
  # Optional: Reference to checklist file
  checklist: checklist.yml
  
  # Optional: Output configuration
  output:
    format: json
    file: {output_file_name}.json
    include:
      - violations
      - warnings
      - recommendations
```

---

## 3. Checklist (checklist.yml) - Optional

Define validation rules and checks the agent performs.

```yaml
{agent_name}_checklist:
  {category_1}:
    - rule: "{Rule description}"
      severity: {critical|high|medium|low}
      pattern: '{regex_pattern}'  # Optional
      check: "{what_to_check}"    # Optional
      suggestion: "{fix_suggestion}"
      reason: "{why_this_matters}"
      
    - rule: "{Another rule}"
      severity: {severity_level}
      file_pattern: "**/*.ext"
      check: "{check_description}"
      
  {category_2}:
    - rule: "{Rule description}"
      severity: {severity_level}
      pattern: '{pattern}'
      suggestion: "{suggestion}"

validation_levels:
  critical:
    action: block_merge
    notification: immediate
    
  high:
    action: require_review
    notification: within_24h
    
  medium:
    action: warn
    notification: next_review
    
  low:
    action: inform
    notification: optional
```

---

## 4. Utilities (utilities.yml) - Optional

Provide utilities, templates, or helper configurations.

```yaml
{agent_name}_utilities:
  {utility_category}:
    {utility_name}:
      pattern: "{pattern_or_format}"
      example: "{example}"
      description: "{description}"
      
  validation_checks:
    - name: "{Check name}"
      pattern: '{regex_pattern}'
      action: "{action_to_take}"
      severity: {severity_level}
      message: "{user_message}"
      
  code_templates:
    {template_name}: |
      {multi-line code template}
      
  best_practices:
    - practice: "{Practice description}"
      reason: "{Why this practice matters}"
      examples:
        - "{Example 1}"
        - "{Example 2}"
```

---

## Integration into SQUAD Team

Add your custom agent to `.squad/team.yml`:

```yaml
squad_team:
  name: {team-name}
  version: 2.0.0
  
  agents:
    # ... existing agents ...
    
    # Your custom agent
    - name: {agent-name}
      enabled: true
      role: {Agent Role}
      config: agents/{agent-name}/config.yml
      charter: agents/{agent-name}/charter.md
      custom: true
```

Add quality gates for your agent:

```yaml
  quality_gates:
    code_review:
      - agent: {agent-name}
        check: {check_name}
        required: true
        blocking: {true|false}
        description: "{Check description}"
```

---

## Example: Security Scanning Agent

```markdown
# Security Scanner Agent Charter

## Role
I am the Security Scanner agent, responsible for identifying security vulnerabilities and enforcing secure coding practices.

## Responsibilities
- Scan for SQL injection vulnerabilities
- Detect hardcoded secrets and credentials
- Validate input sanitization
- Check for insecure dependencies
- Review authentication and authorization

## When to Trigger
- Changes to authentication code
- Database query modifications
- API endpoint changes
- Dependency updates

## Review Checklist
- [ ] No hardcoded secrets
- [ ] SQL queries use parameterized statements
- [ ] User input is validated and sanitized
- [ ] No vulnerable dependencies
- [ ] Proper authentication on protected endpoints
```

```yaml
# config.yml
agent:
  name: security-scanner
  role: Security Vulnerability Scanner
  description: Identifies security vulnerabilities and enforces secure coding practices
  
  triggers:
    - file_pattern: "**/*.go"
    - file_pattern: "**/Dockerfile"
    - keyword: "password"
    - keyword: "secret"
    - keyword: "token"
    - keyword: "api_key"
  
  capabilities:
    - secret_detection
    - sql_injection_prevention
    - dependency_scanning
    - authentication_review
  
  review_priority: critical
  
  integration_points:
    - stage: pre_commit
      action: quick_security_scan
      timeout: 30s
      
    - stage: code_review
      action: comprehensive_security_review
      timeout: 120s
      blocking: true
```

---

## Best Practices for Custom Agents

### 1. **Single Responsibility**
Each agent should have a clear, focused purpose. Don't create catch-all agents.

✅ Good: HIPAA Compliance Agent  
❌ Bad: General Compliance and Security and Quality Agent

### 2. **Actionable Feedback**
Agent output should guide developers to solutions, not just flag problems.

✅ Good: "PHI found in log on line 42. Use patient ID instead: log.Info(\"Action\", \"patient_id\", id)"  
❌ Bad: "PHI violation detected"

### 3. **Domain Expertise**
Encode real domain knowledge, not just basic rules.

✅ Good: Validate ICD-10 code format and suggest common codes  
❌ Bad: Check if field is empty

### 4. **Right Priority Level**
Set appropriate severity for different types of issues.

- **Critical:** Security vulnerabilities, compliance violations
- **High:** Important best practices, significant code quality issues
- **Medium:** Naming conventions, documentation gaps
- **Low:** Style preferences, nice-to-haves

### 5. **Clear Triggers**
Define specific patterns and keywords that should trigger the agent.

```yaml
triggers:
  - file_pattern: "internal/patient/**"  # Specific paths
  - keyword: "PHI"                        # Domain terms
  exclude:
    - "**/*.pb.go"                        # Generated code
```

### 6. **Composability**
Agents should work well together, not duplicate functionality.

- HIPAA Agent: Checks PHI handling
- Security Agent: Checks authentication
- Don't make HIPAA Agent check authentication too

### 7. **Evolvability**
Make it easy to update agent behavior as requirements change.

- Use YAML for rules (easy to edit)
- Include inline comments
- Document why rules exist

---

## Testing Your Custom Agent

1. **Validate configuration:**
   ```bash
   copilot squad validate --config .squad/team.yml
   ```

2. **Test agent in isolation:**
   ```bash
   copilot squad run --agent {agent-name} --dry-run
   ```

3. **Run on sample code:**
   ```bash
   copilot squad run --agent {agent-name} --path ./test-code/
   ```

4. **Check output format:**
   ```bash
   cat .squad/outputs/{agent-name}-report.json | jq '.'
   ```

---

## Common Patterns

### Pattern 1: Compliance Agent
- **Purpose:** Enforce regulatory requirements (HIPAA, PCI-DSS, GDPR)
- **Triggers:** Data handling, logging, security
- **Priority:** Critical
- **Blocking:** Yes

### Pattern 2: Terminology Agent
- **Purpose:** Enforce domain-specific naming and standards
- **Triggers:** Field names, code patterns, documentation
- **Priority:** High
- **Blocking:** No (warn and suggest)

### Pattern 3: Quality Gate Agent
- **Purpose:** Enforce quality standards before merge
- **Triggers:** All code changes
- **Priority:** Medium-High
- **Blocking:** Configurable

### Pattern 4: Generator Agent
- **Purpose:** Provide utilities and code generation
- **Triggers:** On-demand or specific patterns
- **Priority:** Low-Medium
- **Blocking:** No

---

## Resources

- [SQUAD Documentation](https://github.com/microsoft/squad)
- [YAML Syntax Guide](https://yaml.org/)
- [Regex Testing Tool](https://regex101.com/)

---

**Happy Agent Building! 🤖✨**
