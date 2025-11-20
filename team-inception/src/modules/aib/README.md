# AIB - AIRchetipo Builder Module

Specialized tools and workflows for creating, customizing, and extending AIRchetipo components including agents, workflows, and complete modules.

## Table of Contents

- [Module Structure](#module-structure)
- [Core Workflows](#core-workflows)
- [Agent Types](#agent-types)
- [Quick Start](#quick-start)
- [Best Practices](#best-practices)

## Module Structure

### 🤖 Agents

**AIRchetipo Builder** - Master builder agent orchestrating all creation workflows with deep knowledge of AIRchetipo architecture and conventions.

### 📋 Workflows

Comprehensive suite for building and maintaining AIRchetipo components.

## Core Workflows

### Creation Workflows

**[create-agent](./workflows/create-agent/README.md)** - Build AIRchetipo agents

- Interactive persona development
- Command structure design
- YAML source compilation to .md

**[create-workflow](./workflows/create-workflow/README.md)** - Design workflows

- Structured multi-step processes
- Configuration validation
- Web bundle support

**[create-module](./workflows/create-module/README.md)** - Build complete modules

- Full module infrastructure
- Agent and workflow integration
- Installation automation

**[module-brief](./workflows/module-brief/README.md)** - Strategic planning

- Module blueprint creation
- Vision and architecture
- Comprehensive analysis

### Editing Workflows

**[edit-agent](./workflows/edit-agent/README.md)** - Modify existing agents

- Persona refinement
- Command updates
- Best practice compliance

**[edit-workflow](./workflows/edit-workflow/README.md)** - Update workflows

- Structure maintenance
- Configuration updates
- Documentation sync

**[edit-module](./workflows/edit-module/README.md)** - Module enhancement

- Component modifications
- Dependency management
- Version control

### Maintenance Workflows

**[convert-legacy](./workflows/convert-legacy/README.md)** - Migration tool

- v4 to v6 conversion
- Structure compliance
- Convention updates

**[audit-workflow](./workflows/audit-workflow/README.md)** - Quality validation

- Structure verification
- Config standards check
- Bloat detection
- Web bundle completeness

**[redoc](./workflows/redoc/README.md)** - Auto-documentation

- Reverse-tree approach
- Technical writer quality
- Convention compliance

## Agent Types

AIB creates three agent architectures:

### Full Module Agent

- Complete persona and role definition
- Command structure with fuzzy matching
- Workflow integration
- Module-specific capabilities

### Hybrid Agent

- Shared core capabilities
- Module-specific extensions
- Cross-module compatibility

### Standalone Agent

- Independent operation
- Minimal dependencies
- Specialized single purpose

## Quick Start

1. **Load AIRchetipo Builder agent** in your IDE
2. **Choose creation type:**
   ```
   *create-agent     # New agent
   *create-workflow  # New workflow
   *create-module    # Complete module
   ```
3. **Follow interactive prompts**

### Example: Creating an Agent

```
User: I need a code review agent
Builder: *create-agent

[Interactive session begins]
- Brainstorming phase (optional)
- Persona development
- Command structure
- Integration points
```

## Use Cases

### Custom Development Teams

Build specialized agents for:

- Domain expertise (legal, medical, finance)
- Company processes
- Tool integrations
- Automation tasks

### Workflow Extensions

Create workflows for:

- Compliance requirements
- Quality gates
- Deployment pipelines
- Custom methodologies

### Complete Solutions

Package modules for:

- Industry verticals
- Technology stacks
- Business processes
- Educational frameworks

## Best Practices

1. **Study existing patterns** - Review AIM/CIS implementations
2. **Follow conventions** - Use established structures
3. **Document thoroughly** - Clear instructions essential
4. **Test iteratively** - Validate during creation
5. **Consider reusability** - Build modular components

## Integration

AIB components integrate with:

- **AIRchetipo Core** - Framework foundation
- **AIM** - Extend development capabilities
- **CIS** - Leverage creative workflows
- **Custom Modules** - Your domain solutions

## Related Documentation

- **[Agent Creation Guide](./workflows/create-agent/README.md)** - Detailed instructions
- **[Module Structure](./workflows/create-module/module-structure.md)** - Architecture patterns
- **[AIM Module](../aim/README.md)** - Reference implementation
- **[Core Framework](../../core/README.md)** - Foundation concepts

---

AIB empowers you to extend AIRchetipo for your specific needs while maintaining framework consistency and power.
