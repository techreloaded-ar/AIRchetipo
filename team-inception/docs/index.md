# AIRchetipo Documentation Index

Complete map of all AIRchetipo v6 documentation with recommended reading paths.

---

## 🎯 Getting Started (Start Here!)

**New users:** Start with one of these based on your situation:

| Your Situation              | Start Here                                                      | Then Read                                                     |
| --------------------------- | --------------------------------------------------------------- | ------------------------------------------------------------- |
| **Brand new to AIRchetipo** | [Quick Start Guide](../src/modules/aim/docs/quick-start.md)     | [AIM Workflows Guide](../src/modules/aim/workflows/README.md) |
| **Upgrading from v4**       | [v4 to v6 Upgrade Guide](./v4-to-v6-upgrade.md)                 | [Quick Start Guide](../src/modules/aim/docs/quick-start.md)   |
| **Brownfield project**      | [Brownfield Guide](../src/modules/aim/docs/brownfield-guide.md) | [Quick Start Guide](../src/modules/aim/docs/quick-start.md)   |

---

## 📋 Core Documentation

### Project-Level Docs (Root)

- **[README.md](../README.md)** - Main project overview, feature summary, and module introductions
- **[CONTRIBUTING.md](../CONTRIBUTING.md)** - How to contribute, pull request guidelines, code style
- **[CHANGELOG.md](../CHANGELOG.md)** - Version history and breaking changes
- **[CLAUDE.md](../CLAUDE.md)** - Claude Code specific guidelines for this project

### Installation & Setup

- **[v4 to v6 Upgrade Guide](./v4-to-v6-upgrade.md)** - Migration path for v4 users
- **[Document Sharding Guide](./document-sharding-guide.md)** - Split large documents for 90%+ token savings
- **[Web Bundles](./USING_WEB_BUNDLES.md)** - Use AIRchetipo agents in Claude Projects, ChatGPT, or Gemini without installation
- **[Bundle Distribution Setup](./BUNDLE_DISTRIBUTION_SETUP.md)** - Maintainer guide for bundle auto-publishing

---

## 🏗️ Module Documentation

### AIRchetipo (AIM) - Software Development

The flagship module for agile AI-driven development.

- **[AIM Module README](../src/modules/aim/README.md)** - Module overview, agents, and complete documentation index
- **[AIM Documentation](../src/modules/aim/docs/)** - All AIM-specific guides and references:
  - [Quick Start Guide](../src/modules/aim/docs/quick-start.md) - Step-by-step guide to building your first project
  - [Quick Spec Flow](../src/modules/aim/docs/quick-spec-flow.md) - Rapid Level 0-1 development
  - [Scale Adaptive System](../src/modules/aim/docs/scale-adaptive-system.md) - Understanding the 5-level system
  - [Brownfield Guide](../src/modules/aim/docs/brownfield-guide.md) - Working with existing codebases
- **[AIM Workflows Guide](../src/modules/aim/workflows/README.md)** - **ESSENTIAL READING**
- **[Test Architect Guide](../src/modules/aim/testarch/README.md)** - Testing strategy and quality assurance

### AIRchetipo Builder (AIB) - Create Custom Solutions

Build your own agents, workflows, and modules.

- **[AIB Module README](../src/modules/aib/README.md)** - Module overview and capabilities
- **[Agent Creation Guide](../src/modules/aib/workflows/create-agent/README.md)** - Design custom agents

---

## 🖥️ IDE-Specific Guides

Instructions for loading agents and running workflows in your development environment.

**Popular IDEs:**

- [Claude Code](./ide-info/claude-code.md)
- [Cursor](./ide-info/cursor.md)
- [VS Code](./ide-info/windsurf.md)

**Other Supported IDEs:**

- [Augment](./ide-info/auggie.md)
- [Cline](./ide-info/cline.md)
- [Codex](./ide-info/codex.md)
- [Crush](./ide-info/crush.md)
- [Gemini](./ide-info/gemini.md)
- [GitHub Copilot](./ide-info/github-copilot.md)
- [IFlow](./ide-info/iflow.md)
- [Kilo](./ide-info/kilo.md)
- [OpenCode](./ide-info/opencode.md)
- [Qwen](./ide-info/qwen.md)
- [Roo](./ide-info/roo.md)
- [Trae](./ide-info/trae.md)

**Key concept:** Every reference to "load an agent" or "activate an agent" in the main docs links to the [ide-info](./ide-info/) directory for IDE-specific instructions.

---

## 🔧 Advanced Topics

### Installation & Bundling

- [IDE Injections Reference](./installers-bundlers/ide-injections.md) - How agents are installed to IDEs
- [Installers & Platforms Reference](./installers-bundlers/installers-modules-platforms-reference.md) - CLI tool and platform support
- [Web Bundler Usage](./installers-bundlers/web-bundler-usage.md) - Creating web-compatible bundles

---

## 📊 Documentation Map

```
docs/                              # Core/cross-module documentation
├── index.md (this file)
├── v4-to-v6-upgrade.md
├── document-sharding-guide.md
├── ide-info/                      # IDE setup guides
│   ├── claude-code.md
│   ├── cursor.md
│   ├── windsurf.md
│   └── [14+ other IDEs]
└── installers-bundlers/           # Installation reference
    ├── ide-injections.md
    ├── installers-modules-platforms-reference.md
    └── web-bundler-usage.md

src/modules/
├── aim/                           # AIRchetipo module
│   ├── README.md                  # Module overview & docs index
│   ├── docs/                      # AIM-specific documentation
│   │   ├── quick-start.md
│   │   ├── quick-spec-flow.md
│   │   ├── scale-adaptive-system.md
│   │   └── brownfield-guide.md
│   ├── workflows/README.md        # ESSENTIAL workflow guide
│   └── testarch/README.md         # Testing strategy
├── aib/                           # AIRchetipo Builder module
│   ├── README.md
│   └── workflows/create-agent/README.md
└── cis/                           # Creative Intelligence Suite
    └── README.md
```

---

## 🎓 Recommended Reading Paths

### Path 1: Brand New to AIRchetipo

1. [README.md](../README.md) - Understand the vision
2. [Quick Start Guide](../src/modules/aim/docs/quick-start.md) - Get hands-on
3. [AIM Module README](../src/modules/aim/README.md) - Understand agents
4. [AIM Workflows Guide](../src/modules/aim/workflows/README.md) - Master the methodology
5. [Your IDE guide](./ide-info/) - Optimize your workflow

### Path 2: Upgrading from v4

1. [v4 to v6 Upgrade Guide](./v4-to-v6-upgrade.md) - Understand what changed
2. [Quick Start Guide](../src/modules/aim/docs/quick-start.md) - Reorient yourself
3. [AIM Workflows Guide](../src/modules/aim/workflows/README.md) - Learn new v6 workflows

### Path 3: Working with Existing Codebase (Brownfield)

1. [Brownfield Guide](../src/modules/aim/docs/brownfield-guide.md) - Approach for legacy code
2. [Quick Start Guide](../src/modules/aim/docs/quick-start.md) - Follow the process
3. [AIM Workflows Guide](../src/modules/aim/workflows/README.md) - Master the methodology

### Path 4: Building Custom Solutions

1. [AIB Module README](../src/modules/aib/README.md) - Understand capabilities
2. [Agent Creation Guide](../src/modules/aib/workflows/create-agent/README.md) - Create agents
3. [AIM Workflows Guide](../src/modules/aim/workflows/README.md) - Understand workflow structure

### Path 5: Contributing to AIRchetipo

1. [CONTRIBUTING.md](../CONTRIBUTING.md) - Contribution guidelines
2. Relevant module README - Understand the area you're contributing to
3. [Code Style section in CONTRIBUTING.md](../CONTRIBUTING.md#code-style) - Follow standards

---

## 🔍 Quick Reference

**What is each module for?**

- **AIM** - AI-driven software development
- **AIB** - Create custom agents and workflows
- **CIS** - Creative thinking and brainstorming

**How do I load an agent?**
→ See [ide-info](./ide-info/) folder for your IDE

**I'm stuck, what's next?**
→ Check the [AIM Workflows Guide](../src/modules/aim/workflows/README.md) or run `workflow-status`

**I want to contribute**
→ Start with [CONTRIBUTING.md](../CONTRIBUTING.md)

---

## 📚 Important Concepts

### Fresh Chats

Each workflow should run in a fresh chat with the specified agent to avoid context limitations. This is emphasized throughout the docs because it's critical to successful workflows.

### Scale Levels

AIM adapts to project complexity (Levels 0-4). Documentation is scale-adaptive - you only need what's relevant to your project size.

### Update-Safe Customization

All agent customizations go in `{air_folder}/_cfg/agents/` and survive updates. See your IDE guide and module README for details.

---

## 🆘 Getting Help

- **Discord**: [Join the AIRchetipo Community](https://discord.gg/gk8jAdXWmj)
  - #general-dev - Technical questions
  - #bugs-issues - Bug reports
- **Issues**: [GitHub Issue Tracker](https://github.com/airchetipo-org/AIRchetipo/issues)
- **YouTube**: [AIRchetipo Code Channel](https://www.youtube.com/@AIRchetipoCode)
