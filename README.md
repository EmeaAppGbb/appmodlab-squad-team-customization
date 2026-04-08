# ⚡ SQUAD Team Customization Lab 🎮

```
███████╗ ██████╗ ██╗   ██╗ █████╗ ██████╗     
██╔════╝██╔═══██╗██║   ██║██╔══██╗██╔══██╗    
███████╗██║   ██║██║   ██║███████║██║  ██║    
╚════██║██║▄▄ ██║██║   ██║██╔══██║██║  ██║    
███████║╚██████╔╝╚██████╔╝██║  ██║██████╔╝    
╚══════╝ ╚══▀▀═╝  ╚═════╝ ╚═╝  ╚═╝╚═════╝     
                                               
 ██████╗██╗   ██╗███████╗████████╗ ██████╗ ███╗   ███╗
██╔════╝██║   ██║██╔════╝╚══██╔══╝██╔═══██╗████╗ ████║
██║     ██║   ██║███████╗   ██║   ██║   ██║██╔████╔██║
██║     ██║   ██║╚════██║   ██║   ██║   ██║██║╚██╔╝██║
╚██████╗╚██████╔╝███████║   ██║   ╚██████╔╝██║ ╚═╝ ██║
 ╚═════╝ ╚═════╝ ╚══════╝   ╚═╝    ╚═════╝ ╚═╝     ╚═╝
```

## 🌟 QUEST OVERVIEW

**BUILD YOUR ULTIMATE DEV PARTY! 🎲⚔️**

Transform your SQUAD from a stock team into a **custom-tailored dream team** of AI agents specialized for YOUR domain! Like assembling the perfect RPG party, you'll recruit domain-specific agents, configure their skill trees, and set up team ceremonies to tackle unique challenges.

**🎯 Mission:** Customize SQUAD for a healthcare SaaS app that needs HIPAA compliance agents, medical terminology validators, and data anonymizers — agents that don't exist in the base game!

---

## 📦 WHAT'S IN THE BOX

**🏥 Demo Application: MedBook Appointment Scheduling**

A Go-based microservices backend for healthcare appointment scheduling that needs specialized agents to handle:
- ⚕️ HIPAA compliance validation
- 🏥 Healthcare terminology checking (ICD-10, CPT codes)
- 🔐 Patient data anonymization

**💎 Tech Stack:**
- 🔹 Go 1.22 + Gin web framework
- 🔹 gRPC for inter-service communication
- 🔹 PostgreSQL with pgx driver
- 🔹 Protocol Buffers
- 🔹 Kubernetes deployment
- 🔹 Base SQUAD configuration (ready to customize!)

---

## 🎮 LEVEL SELECT (BRANCH STRUCTURE)

| Branch | Description | Status |
|--------|-------------|--------|
| `legacy` 🏚️ | Standard SQUAD config — vanilla party | ⭐ START HERE |
| `step-1-custom-agent-definition` 📝 | Define agent charters & capabilities | 🎯 QUEST 1 |
| `step-2-agent-implementation` ⚙️ | Implement custom agent configs | 🎯 QUEST 2 |
| `step-3-ceremonies` 🎭 | Configure team ceremonies | 🎯 QUEST 3 |
| `step-4-skills-and-gates` 🌳 | Add skills & quality gates | 🎯 QUEST 4 |
| `step-5-integration-test` 🧪 | Test your custom party! | 🎯 QUEST 5 |
| `solution` 🏆 | Fully customized SQUAD team | ✨ VICTORY! |
| `main` 📚 | Complete lab documentation | 📖 GUIDE |

---

## 🎯 LEARNING OBJECTIVES

**SKILLS YOU'LL UNLOCK 🌟**

By completing this lab, you'll master:

✅ **CLASS SELECTION** 🧙 — Define custom SQUAD agents with specialized roles  
✅ **CHARTER WRITING** 📜 — Create domain-specific agent instructions  
✅ **CEREMONY CONFIG** 🎭 — Set up standups, retros, and compliance reviews  
✅ **SKILL TREE UNLOCKED** 🌳 — Build domain-specific quality gates  
✅ **PARTY COMPOSITION** ⚔️ — Integrate custom agents into SQUAD workflow  
✅ **BATTLE TESTING** ⚔️ — Validate your custom agents on real tasks  

---

## ⚡ PREREQUISITES

**REQUIRED XP 📊**

Before starting this quest, make sure you have:

- ✅ Completed "Getting Started with SQUAD" lab
- ✅ Basic Go language familiarity (for reading the codebase)
- ✅ Understanding of YAML configuration
- ✅ GitHub Copilot access
- ✅ Git installed and configured

**⏱️ Estimated Duration:** 3–4 hours (grab some ☕!)

---

## 🚀 QUICK START

### 🎬 GAME START!

```bash
# 🌀 CLONE THE QUEST
git clone https://github.com/EmeaAppGbb/appmodlab-squad-team-customization.git
cd appmodlab-squad-team-customization

# 🏚️ START WITH LEGACY (vanilla SQUAD)
git checkout legacy

# 📖 OPEN THE COMPLETE GUIDE
# Check out APPMODLAB.md for full step-by-step instructions!
```

---

## 🎨 WHAT YOU'LL BUILD

### 🦸 YOUR CUSTOM PARTY ROSTER

**⚕️ HIPAA Compliance Agent**
- 🛡️ **Class:** Compliance Guardian
- 🎯 **Special Ability:** Scans code for PHI exposure, encryption gaps, audit logging
- 💬 **Battle Cry:** "No PHI left unprotected!"

**🏥 Domain Terminology Agent**
- 📚 **Class:** Medical Lexicon Master
- 🎯 **Special Ability:** Validates ICD-10/CPT codes, healthcare terminology
- 💬 **Battle Cry:** "Proper terminology saves lives!"

**🔐 Data Anonymizer Agent**
- 🎭 **Class:** Privacy Rogue
- 🎯 **Special Ability:** Generates anonymized test data, validates fixtures
- 💬 **Battle Cry:** "Real patients stay private!"

### 🎭 CUSTOM CEREMONIES CONFIGURED

- 📅 **Weekly HIPAA Review** — Security-focused retrospective
- 🗓️ **Sprint Planning with Compliance** — Factor in regulatory requirements
- 🔄 **Daily Healthcare Standup** — Domain-specific blockers

### 🌳 SKILL TREES UNLOCKED

- 🔹 Go-specific linting rules
- 🔹 HIPAA compliance checklists
- 🔹 Healthcare terminology validation
- 🔹 Test data anonymization patterns
- 🔹 Domain-specific quality gates

---

## 📁 REPOSITORY STRUCTURE

```
appmodlab-squad-team-customization/
├── 📖 README.md                    # You are here! 🎮
├── 📚 APPMODLAB.md                 # Complete step-by-step guide
├── 🏥 medbook/                     # Healthcare SaaS codebase
│   ├── cmd/                        # Microservices entry points
│   │   ├── appointment-service/    # Appointment scheduling
│   │   ├── provider-service/       # Healthcare provider mgmt
│   │   └── patient-service/        # Patient records
│   ├── internal/                   # Domain logic
│   │   ├── appointment/
│   │   ├── provider/
│   │   └── patient/
│   ├── proto/                      # Protocol Buffer definitions
│   ├── k8s/                        # Kubernetes manifests
│   ├── .squad/                     # ⚡ SQUAD CONFIGURATION
│   │   ├── team.yml                # Team config (customize this!)
│   │   └── agents/                 # Agent definitions
│   │       ├── brain/              # 🧠 Brain (standard)
│   │       ├── hands/              # 🙌 Hands (standard)
│   │       ├── eyes/               # 👀 Eyes (standard)
│   │       ├── mouth/              # 💬 Mouth (standard)
│   │       ├── hipaa/              # ⚕️ HIPAA (CUSTOM!)
│   │       ├── terminology/        # 🏥 Terminology (CUSTOM!)
│   │       └── anonymizer/         # 🔐 Anonymizer (CUSTOM!)
│   └── tests/                      # Test suites
└── .github/                        # CI/CD workflows
    └── workflows/
        └── squad-gates.yml         # Custom quality gates
```

---

## 🎯 KEY CONCEPTS

### 🧩 SQUAD ARCHITECTURE 101

**Base Party (Standard SQUAD):**
- 🧠 **Brain** — Architectural decisions
- 🙌 **Hands** — Code implementation
- 👀 **Eyes** — Code review
- 💬 **Mouth** — Communication & documentation

**🎨 Customization Powers:**
1. **Custom Agents** — Add domain-specific party members
2. **Agent Charters** — Define specialized review criteria
3. **Ceremonies** — Configure team rituals (standups, retros)
4. **Skills** — Build custom capabilities (linters, validators)
5. **Quality Gates** — Require custom agent approval before merge

---

## 🏆 ACCEPTANCE CRITERIA

**QUEST COMPLETION CHECKLIST ✅**

- [ ] 🏥 Base MedBook codebase runs with standard SQUAD
- [ ] ⚕️ HIPAA Compliance Agent defined with review checklist
- [ ] 🏥 Domain Terminology Agent validates healthcare terms
- [ ] 🔐 Data Anonymizer Agent generates anonymized test data
- [ ] 🎭 Custom ceremonies configured and documented
- [ ] 🛡️ Quality gates require custom agent approval
- [ ] ⚔️ Custom agents produce meaningful output on healthcare code
- [ ] 📖 APPMODLAB.md guide is complete and reproducible
- [ ] 💬 Configuration files have inline documentation

---

## 🎓 LEARNING PATH

### 🗺️ RECOMMENDED QUEST ORDER

1. **📖 Read the Manual** — Review `APPMODLAB.md` for full instructions
2. **🔍 Explore Base Config** — Check out `legacy` branch, understand vanilla SQUAD
3. **📝 Define Agents** — Branch `step-1-custom-agent-definition`
4. **⚙️ Implement Config** — Branch `step-2-agent-implementation`
5. **🎭 Setup Ceremonies** — Branch `step-3-ceremonies`
6. **🌳 Add Skills & Gates** — Branch `step-4-skills-and-gates`
7. **🧪 Battle Test** — Branch `step-5-integration-test`
8. **🏆 Compare Solution** — Check `solution` branch
9. **🎨 Customize Further** — Make it your own!

---

## 💡 PRO TIPS

**⚡ POWER-UPS FOR SUCCESS**

- 🎯 **Start Small** — Add one custom agent before building all three
- 📝 **Write Good Charters** — Clear agent instructions = better results
- 🧪 **Test Early** — Validate each agent as you build it
- 🔄 **Iterate** — Refine agent behavior based on real output
- 📖 **Document Everything** — Future teams will thank you!
- 🎨 **Get Creative** — SQUAD is flexible — experiment!

---

## 🌈 SQUAD CUSTOMIZATION PHILOSOPHY

> **"SQUAD isn't a rigid framework — it's a platform for building YOUR perfect team!"** 🎮

Every organization has unique:
- 🏢 Domain requirements (healthcare, finance, gaming, etc.)
- 🔧 Tech stack preferences (Go, Python, Java, etc.)
- 📋 Quality standards (compliance, security, performance)
- 🎭 Team rituals (ceremonies, workflows, gates)

SQUAD's superpower is **flexibility** — you can tailor it to YOUR world! 🌍✨

---

## 🤝 CONTRIBUTING

Want to enhance this lab? 🚀

1. 🍴 Fork the repository
2. 🌿 Create a feature branch (`git checkout -b feature/amazing-agent`)
3. ✨ Make your changes
4. 📝 Commit with descriptive messages
5. 🚀 Push to your fork
6. 🎯 Open a Pull Request

**Ideas for contributions:**
- 🎨 More custom agent examples (FinTech, Gaming, etc.)
- 📚 Additional ceremony templates
- 🛠️ Tool integrations
- 📖 Documentation improvements
- 🐛 Bug fixes

---

## 📚 ADDITIONAL RESOURCES

**🔗 HELPFUL LINKS**

- 📖 [SQUAD Documentation](https://github.com/microsoft/squad) — Core framework
- 🏥 [HIPAA Technical Safeguards](https://www.hhs.gov/hipaa/for-professionals/security/index.html) — Healthcare compliance
- 🔹 [Go Best Practices](https://golang.org/doc/effective_go) — Language guide
- 🎯 [GitHub Copilot](https://github.com/features/copilot) — AI pair programming
- 📦 [Protocol Buffers](https://protobuf.dev/) — Service contracts

---

## 📞 SUPPORT & FEEDBACK

**🆘 NEED HELP?**

- 💬 Open an issue in this repository
- 📧 Contact the App Modernization team
- 🔍 Check `APPMODLAB.md` for detailed troubleshooting
- 🎮 Join our community discussions

---

## 📜 LICENSE

This lab is part of the App Modernization Labs collection.

---

## 🎮 READY PLAYER ONE?

```
╔══════════════════════════════════════════════╗
║  🎯 PRESS START TO BEGIN YOUR QUEST! 🎯     ║
║                                              ║
║  git checkout legacy                         ║
║  code APPMODLAB.md                           ║
║                                              ║
║  ⚡ MAY YOUR AGENTS BE EVER SPECIALIZED! ⚡  ║
╚══════════════════════════════════════════════╝
```

---

**Built with 💜 by the Azure App Modernization GBB Team**

**🌟 POWER UP YOUR SQUAD! CUSTOMIZE YOUR DESTINY! 🌟**

---

### 🏁 ACHIEVEMENT UNLOCKED

**"Custom Agent Master"** 🏆  
*You've learned to build domain-specific AI agents and create the ultimate dev team!*

```
    ⭐⭐⭐⭐⭐
   LEGEND STATUS
   ⭐⭐⭐⭐⭐
```

**NOW GO FORTH AND BUILD! 🚀✨**
