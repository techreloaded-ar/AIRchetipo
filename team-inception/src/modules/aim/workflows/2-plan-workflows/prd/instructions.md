# PRD Workflow - Intent-Driven Product Planning

<critical>The workflow execution engine is governed by: {project-root}/{air_folder}/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {installed_path}/workflow.yaml</critical>
<critical>This workflow uses INTENT-DRIVEN PLANNING - adapt organically to product type and context</critical>
<critical>Communicate all responses in {communication_language} and adapt deeply to {user_skill_level}</critical>
<critical>Generate all documents in {document_output_language}</critical>
<critical>LIVING DOCUMENT: Write to PRD.md continuously as you discover - never wait until the end</critical>
<critical>GUIDING PRINCIPLE: Identify what makes this product special and ensure it's reflected throughout the PRD</critical>
<critical>Input documents specified in workflow.yaml input_file_patterns - workflow engine handles fuzzy matching, whole vs sharded document discovery automatically</critical>

<workflow>

<step n="0" goal="Validate workflow readiness" tag="workflow-status">
<action>Check if {status_file} exists</action>

<action if="status file not found">Set standalone_mode = true</action>

<check if="status file found">
  <action>Load the FULL file: {status_file}</action>
  <action>Parse workflow_status section</action>
  <action>Check status of "prd" workflow</action>
  <action>Get project_track from YAML metadata</action>
  <action>Find first non-completed workflow (next expected workflow)</action>

  <check if="project_track is Quick Flow">
    <output>**Quick Flow Track - Redirecting**

Quick Flow projects use tech-spec workflow for implementation-focused planning.
PRD is for AIRchetipo and Enterprise Method tracks that need comprehensive requirements.</output>
<action>Exit and suggest tech-spec workflow</action>
</check>

  <check if="prd status is file path (already completed)">
    <output>⚠️ PRD already completed: {{prd status}}</output>
    <ask>Re-running will overwrite the existing PRD. Continue? (y/n)</ask>
    <check if="n">
      <output>Exiting. Use workflow-status to see your next step.</output>
      <action>Exit workflow</action>
    </check>
  </check>

<action>Set standalone_mode = false</action>
</check>
</step>

<step n="0.5" goal="Discover and load input documents">
<invoke-protocol name="discover_inputs" />
<note>After discovery, these content variables are available: {product_brief_content}, {research_content}, {document_project_content}</note>
</step>

<step n="1" goal="Discovery - Project, Domain, and Vision">
<action>Welcome {user_name} and begin comprehensive discovery, and then start to GATHER ALL CONTEXT:
1. Check workflow-status.yaml for project_context (if exists)
2. Review loaded content: {product_brief_content}, {research_content}, {document_project_content} (auto-loaded in Step 0.5)
3. Detect project type AND domain complexity

Load references:
{installed_path}/project-types.csv
{installed_path}/domain-complexity.csv

Through natural conversation:
"Tell me about what you want to build - what problem does it solve and for whom?"

DUAL DETECTION:
Project type signals: API, mobile, web, CLI, SDK, SaaS
Domain complexity signals: medical, finance, government, education, aerospace

SPECIAL ROUTING:
If game detected → Inform user that game development requires the AIGD module (AIRchetipo Game Development)
If complex domain detected → Offer domain research options:
A) Run domain-research workflow (thorough)
B) Quick web search (basic)
C) User provides context
D) Continue with general knowledge

IDENTIFY WHAT MAKES IT SPECIAL early with questions such as: "What excites you most about this product?", "What would make users love this?", "What's the unique value or compelling moment?"

This becomes a thread that connects throughout the PRD.</action>

<template-output>vision_alignment</template-output>
<template-output>project_classification</template-output>
<template-output>project_type</template-output>
<template-output>domain_type</template-output>
<template-output>complexity_level</template-output>
<check if="complex domain">
<template-output>domain_context_summary</template-output>
</check>
<template-output>product_differentiator</template-output>
<template-output>product_brief_path</template-output>
<template-output>domain_brief_path</template-output>
<template-output>research_documents</template-output>
</step>

<step n="1.5" goal="Vision Articulation">
<action>Build a clear, compelling vision statement

Through conversation, extract:

"Where do you see this product in 3-5 years?"
"What change do you want to see in the world because of this product?"
"What's the ultimate impact you're aiming for?"

INTENT: Create a north star that guides all decisions

Structure the vision:

1. **Vision Statement** - One compelling sentence capturing the future state
2. **Strategic Objectives** - 3-5 key objectives that move toward the vision
3. **Long-term Impact** - The transformational change this product will create

Connect vision to what makes it special:
The vision should amplify the unique value proposition discovered in Step 1.</action>

<template-output>vision_statement</template-output>
<template-output>strategic_objectives</template-output>
<template-output>long_term_impact</template-output>
</step>

<step n="1.6" goal="Business Model Canvas (Simplified)">
<action>Create a simplified business model overview

Through natural conversation, map out the key business elements:

SIMPLIFIED CANVAS - Focus on 5 core areas:

1. **Value Proposition** - What unique value does this create?
2. **Customer Segments** - Who are the target customers? (will detail in Personas next)
3. **Revenue Model** - How will this generate revenue/value?
4. **Key Resources** - What critical resources are needed?
5. **Cost Structure** - What are the main cost drivers?

For B2C products: Focus on user value and engagement model
For B2B products: Focus on business value and pricing tiers
For internal tools: Focus on efficiency gains and resource optimization
For open source: Focus on community value and sustainability

Keep it high-level - this is strategic context, not a full business plan.</action>

<template-output>business_model_canvas</template-output>
</step>

<step n="1.7" goal="Persona Definition">
<action>Define the 2 primary user personas

For each persona, create a rich profile through conversation:

"Who are the main types of users for this product?"
"What are their goals, frustrations, and contexts?"

For each of 2 personas, capture:

- **Name and Role** (fictional but representative)
- **Background** - Context, experience level, environment
- **Goals** - What they want to achieve
- **Pain Points** - Current frustrations and challenges
- **Behaviors** - How they work, what tools they use
- **Motivations** - What drives them
- **Tech Savviness** - Their comfort level with technology

INTENT: Create empathy anchors for all design decisions

These personas should reflect the customer segments from the Business Model Canvas.</action>

<template-output>persona_1_name</template-output>
<template-output>persona_1_profile</template-output>
<template-output>persona_2_name</template-output>
<template-output>persona_2_profile</template-output>
</step>

<step n="1.8" goal="Customer Journey Mapping">
<action>Map the journey for each persona

For each persona, map their journey with the product:

JOURNEY STRUCTURE:

1. **Awareness** - How do they discover the product?
2. **Consideration** - What makes them consider using it?
3. **First Use** - Their initial experience (critical!)
4. **Regular Use** - How they engage day-to-day
5. **Advocacy** - What makes them recommend it?

For each stage, identify:

- Key touchpoints
- Emotions (frustrations, delights)
- Opportunities for value delivery
- Potential drop-off points

INTENT: Understand the complete user experience from discovery to advocacy

Connect journey stages to functional requirements:
Each stage should inform what capabilities the product needs.</action>

<template-output>persona_1_journey</template-output>
<template-output>persona_2_journey</template-output>
</step>

<step n="2" goal="Success Definition">
<action>Define what winning looks like for THIS specific product

INTENT: Meaningful success criteria, not generic metrics

Adapt to context:

- Consumer: User love, engagement, retention
- B2B: ROI, efficiency, adoption
- Developer tools: Developer experience, community
- Regulated: Compliance, safety, validation

Make it specific:

- NOT: "10,000 users"
- BUT: "100 power users who rely on it daily"

- NOT: "99.9% uptime"
- BUT: "Zero data loss during critical operations"

Connect to what makes the product special:

- "Success means users experience [key value moment] and achieve [desired outcome]"</action>

<template-output>success_criteria</template-output>
<check if="business focus">
<template-output>business_metrics</template-output>
</check>
</step>

<step n="3" goal="Scope Definition">
<action>Smart scope negotiation - find the sweet spot

The Scoping Game:

1. "What must work for this to be useful?" → MVP
2. "What makes it competitive?" → Growth
3. "What's the dream version?" → Vision

Challenge scope creep conversationally:

- "Could that wait until after launch?"
- "Is that essential for proving the concept?"

For complex domains:

- Include compliance minimums in MVP
- Note regulatory gates between phases</action>

<template-output>mvp_scope</template-output>
<template-output>growth_features</template-output>
<template-output>vision_features</template-output>
</step>

<step n="4" goal="Domain-Specific Exploration" optional="true">
<action>Only if complex domain detected or domain-brief exists

Synthesize domain requirements that will shape everything:

- Regulatory requirements
- Compliance needs
- Industry standards
- Safety/risk factors
- Required validations
- Special expertise needed

These inform:

- What features are mandatory
- What NFRs are critical
- How to sequence development
- What validation is required</action>

<check if="complex domain">
  <template-output>domain_considerations</template-output>
</check>
</step>

<step n="5" goal="Innovation Discovery" optional="true">
<action>Identify truly novel patterns if applicable

Listen for innovation signals:

- "Nothing like this exists"
- "We're rethinking how [X] works"
- "Combining [A] with [B] for the first time"

Explore deeply:

- What makes it unique?
- What assumption are you challenging?
- How do we validate it?
- What's the fallback?

<WebSearch if="novel">{concept} innovations {date}</WebSearch></action>

<check if="innovation detected">
  <template-output>innovation_patterns</template-output>
  <template-output>validation_approach</template-output>
</check>
</step>

<step n="6" goal="Project-Specific Deep Dive">
<action>Based on detected project type, dive deep into specific needs

Load project type requirements from CSV and expand naturally.

FOR API/BACKEND:

- Map out endpoints, methods, parameters
- Define authentication and authorization
- Specify error codes and rate limits
- Document data schemas

FOR MOBILE:

- Platform requirements (iOS/Android/both)
- Device features needed
- Offline capabilities
- Store compliance

FOR SAAS B2B:

- Multi-tenant architecture
- Permission models
- Subscription tiers
- Critical integrations

[Continue for other types...]

Always connect requirements to product value:
"How does [requirement] support the product's core value proposition?"</action>

<template-output>project_type_requirements</template-output>

<!-- Dynamic sections based on project type -->
<check if="API/Backend project">
  <template-output>endpoint_specification</template-output>
  <template-output>authentication_model</template-output>
</check>

<check if="Mobile project">
  <template-output>platform_requirements</template-output>
  <template-output>device_features</template-output>
</check>

<check if="SaaS B2B project">
  <template-output>tenant_model</template-output>
  <template-output>permission_matrix</template-output>
</check>
</step>

<step n="7" goal="UX Principles" if="project has UI or UX">
  <action>Only if product has a UI

Light touch on UX - not full design:

- Visual personality
- Key interaction patterns
- Critical user flows

"How should this feel to use?"
"What's the vibe - professional, playful, minimal?"

Connect UX to product vision:
"The UI should reinforce [core value proposition] through [design approach]"</action>

  <check if="has UI">
    <template-output>ux_principles</template-output>
    <template-output>key_interactions</template-output>
  </check>
</step>

<step n="8" goal="Functional Requirements Synthesis">
<critical>This section is THE CAPABILITY CONTRACT for all downstream work</critical>
<critical>UX designers will ONLY design what's listed here</critical>
<critical>Architects will ONLY support what's listed here</critical>
<critical>Epic breakdown will ONLY implement what's listed here</critical>
<critical>If a capability is missing from FRs, it will NOT exist in the final product</critical>

<action>Before writing FRs, understand their PURPOSE and USAGE:

**Purpose:**
FRs define WHAT capabilities the product must have. They are the complete inventory
of user-facing and system capabilities that deliver the product vision.

**How They Will Be Used:**

1. UX Designer reads FRs → designs interactions for each capability
2. Architect reads FRs → designs systems to support each capability
3. PM reads FRs → creates epics and stories to implement each capability
4. Dev Agent reads assembled context → implements stories based on FRs

**Critical Property - COMPLETENESS:**
Every capability discussed in vision, scope, domain requirements, and project-specific
sections MUST be represented as an FR. Missing FRs = missing capabilities.

**Critical Property - ALTITUDE:**
FRs state WHAT capability exists and WHO it serves, NOT HOW it's implemented or
specific UI/UX details. Those come later from UX and Architecture.
</action>

<action>Transform everything discovered into comprehensive functional requirements:

**Coverage - Pull from EVERYWHERE:**

- Core features from MVP scope → FRs
- Growth features → FRs (marked as post-MVP if needed)
- Domain-mandated features → FRs
- Project-type specific needs → FRs
- Innovation requirements → FRs
- Anti-patterns (explicitly NOT doing) → Note in FR section if needed

**Organization - Group by CAPABILITY AREA:**
Don't organize by technology or layer. Group by what users/system can DO:

- ✅ "User Management" (not "Authentication System")
- ✅ "Content Discovery" (not "Search Algorithm")
- ✅ "Team Collaboration" (not "WebSocket Infrastructure")

**Format - Flat, Numbered List:**
Each FR is one clear capability statement:

- FR#: [Actor] can [capability] [context/constraint if needed]
- Number sequentially (FR1, FR2, FR3...)
- Aim for 20-50 FRs for typical projects (fewer for simple, more for complex)

**Altitude Check:**
Each FR should answer "WHAT capability exists?" NOT "HOW is it implemented?"

- ✅ "Users can customize appearance settings"
- ❌ "Users can toggle light/dark theme with 3 font size options stored in LocalStorage"

The second example belongs in Epic Breakdown, not PRD.
</action>

<example>
**Well-written FRs at the correct altitude:**

**User Account & Access:**

- FR1: Users can create accounts with email or social authentication
- FR2: Users can log in securely and maintain sessions across devices
- FR3: Users can reset passwords via email verification
- FR4: Users can update profile information and preferences
- FR5: Administrators can manage user roles and permissions

**Content Management:**

- FR6: Users can create, edit, and delete content items
- FR7: Users can organize content with tags and categories
- FR8: Users can search content by keyword, tag, or date range
- FR9: Users can export content in multiple formats

**Data Ownership (local-first products):**

- FR10: All user data stored locally on user's device
- FR11: Users can export complete data at any time
- FR12: Users can import previously exported data
- FR13: System monitors storage usage and warns before limits

**Collaboration:**

- FR14: Users can share content with specific users or teams
- FR15: Users can comment on shared content
- FR16: Users can track content change history
- FR17: Users receive notifications for relevant updates

**Notice:**
✅ Each FR is a testable capability
✅ Each FR is implementation-agnostic (could be built many ways)
✅ Each FR specifies WHO and WHAT, not HOW
✅ No UI details, no performance numbers, no technology choices
✅ Comprehensive coverage of capability areas
</example>

<action>Generate the complete FR list by systematically extracting capabilities:

1. MVP scope → extract all capabilities → write as FRs
2. Growth features → extract capabilities → write as FRs (note if post-MVP)
3. Domain requirements → extract mandatory capabilities → write as FRs
4. Project-type specifics → extract type-specific capabilities → write as FRs
5. Innovation patterns → extract novel capabilities → write as FRs

Organize FRs by logical capability groups (5-8 groups typically).
Number sequentially across all groups (FR1, FR2... FR47).
</action>

<action>SELF-VALIDATION - Before finalizing, ask yourself:

**Completeness Check:**

1. "Did I cover EVERY capability mentioned in the MVP scope section?"
2. "Did I include domain-specific requirements as FRs?"
3. "Did I cover the project-type specific needs (API/Mobile/SaaS/etc)?"
4. "Could a UX designer read ONLY the FRs and know what to design?"
5. "Could an Architect read ONLY the FRs and know what to support?"
6. "Are there any user actions or system behaviors we discussed that have no FR?"

**Altitude Check:**

1. "Am I stating capabilities (WHAT) or implementation (HOW)?"
2. "Am I listing acceptance criteria or UI specifics?" (Remove if yes)
3. "Could this FR be implemented 5 different ways?" (Good - means it's not prescriptive)

**Quality Check:**

1. "Is each FR clear enough that someone could test whether it exists?"
2. "Is each FR independent (not dependent on reading other FRs to understand)?"
3. "Did I avoid vague terms like 'good', 'fast', 'easy'?" (Use NFRs for quality attributes)

COMPLETENESS GATE: Review your FR list against the entire PRD written so far.
Did you miss anything? Add it now before proceeding.
</action>

<template-output>functional_requirements_complete</template-output>
</step>

<step n="9" goal="Non-Functional Requirements Discovery">
<action>Only document NFRs that matter for THIS product

Performance: Only if user-facing impact
Security: Only if handling sensitive data
Scale: Only if growth expected
Accessibility: Only if broad audience
Integration: Only if connecting systems

For each NFR:

- Why it matters for THIS product
- Specific measurable criteria
- Domain-driven requirements

Skip categories that don't apply!</action>

<!-- Only output sections that were discussed -->
<check if="performance matters">
  <template-output>performance_requirements</template-output>
</check>
<check if="security matters">
  <template-output>security_requirements</template-output>
</check>
<check if="scale matters">
  <template-output>scalability_requirements</template-output>
</check>
<check if="accessibility matters">
  <template-output>accessibility_requirements</template-output>
</check>
<check if="integration matters">
  <template-output>integration_requirements</template-output>
</check>
</step>

<step n="9.5" goal="High-Level Architecture">
<action>Define the technical architecture at a high level

INTENT: Provide architectural direction without over-specifying implementation details

Through conversation, determine:

"What kind of architecture are you envisioning?"
"Are there technology preferences or constraints?"
"What's the expected scale and complexity?"

ARCHITECTURE AREAS TO COVER:

1. **High-Level Architecture**
   - Overall system architecture (monolith, microservices, serverless, etc.)
   - Key components and their relationships
   - Data flow and integration patterns
   - Architectural diagram description (textual)

2. **Technology Stack**
   - Primary programming language(s)
   - Frontend framework (if applicable)
   - Backend framework
   - API style (REST, GraphQL, gRPC, etc.)

3. **Database and Persistence**
   - Database type (relational, document, graph, time-series, etc.)
   - Specific database technology (PostgreSQL, MongoDB, etc.)
   - Data modeling approach
   - Caching strategy if relevant

4. **Frameworks and Libraries**
   - Core frameworks for each layer
   - Key libraries for critical functionality
   - Development and build tools

5. **Infrastructure** (optional)
   - Cloud provider preference (AWS, Azure, GCP, on-premise)
   - Hosting approach (containers, serverless, VMs)
   - CI/CD considerations

ADAPTATION RULES:

- For MVPs: Favor proven, simple technologies
- For scale: Consider distributed architecture early
- For specific domains: Respect domain-specific tech (e.g., real-time needs WebSocket)
- For regulated industries: Consider compliance requirements in tech choices

Connect architecture to requirements:

- FRs drive what components are needed
- NFRs drive how components are architected
- Project type influences architectural patterns

Keep it high-level: Detailed technical decisions belong in the Architecture workflow.</action>

<template-output>high_level_architecture</template-output>
<template-output>technology_stack</template-output>
<template-output>database_architecture</template-output>
<template-output>frameworks_and_libraries</template-output>
<check if="infrastructure discussed">
<template-output>infrastructure_overview</template-output>
</check>
</step>

<step n="10" goal="Epic Breakdown">
<action>Transform functional requirements into implementable epics and high-level stories

CRITICAL: Epics are the bridge between strategic requirements and tactical implementation

EPIC BREAKDOWN PROCESS:

1. **Review all Functional Requirements** from Step 8
   - Group related FRs into logical capability domains
   - Identify natural implementation boundaries

2. **Create Epics** - Each epic should:
   - Represent a significant capability area
   - Deliver standalone value when completed
   - Be implementable in 2-4 weeks (typical)
   - Include 5-15 user stories

3. **Define High-Level User Stories** for each epic:
   - Use standard format: "As a [persona], I want [capability], so that [benefit]"
   - Keep stories at high level (detailed AC comes later)
   - Ensure stories cover all FRs
   - Identify dependencies between stories

4. **Prioritization**:
   - Mark epics as MVP, Growth, or Vision (based on scope from Step 3)
   - Order epics by dependency and value delivery
   - Highlight critical path epics

EPIC STRUCTURE:

For each epic, provide:

- Epic title and goal
- Which FRs it implements
- 5-15 high-level user stories
- Dependencies on other epics
- Estimated complexity (Small/Medium/Large)
- Priority (MVP/Growth/Vision)

STORY ALTITUDE:
Stories should be high-level - detailed acceptance criteria come during UX and Architecture workflows.

Example good story: "As a user, I can create and save custom reports"
Example too detailed: "As a user, I can click the 'New Report' button, select from 5 chart types, configure X/Y axes with dropdowns, and save to LocalStorage with auto-generated UUID"

COVERAGE VALIDATION:
Before completing, verify every FR from Step 8 is covered by at least one story.
Create an FR coverage map showing which epic/story implements each FR.</action>

<template-output>epics_overview</template-output>
<template-output>epic_details</template-output>
</step>

<step n="10.5" goal="Roadmap Planning">
<action>Create a high-level roadmap that sequences epic delivery over time

INTENT: Provide a strategic timeline without over-committing to specific dates

ROADMAP STRUCTURE - Organize by PHASES, not specific dates:

Think in phases/quarters rather than weeks:

- **Phase 1 / Q1**: Foundation & MVP
- **Phase 2 / Q2**: Growth Features
- **Phase 3 / Q3**: Scale & Optimization
- **Phase 4 / Q4**: Vision Features

For each phase, define:

1. **Goals** - What this phase achieves
2. **Epics Included** - Which epics from Step 10 are delivered
3. **Key Milestones** - Major achievements or releases
4. **Success Metrics** - How to measure phase completion
5. **Dependencies** - What must be complete from previous phases

PRIORITIZATION PRINCIPLES:

- Foundation first: Core capabilities that others depend on
- Value early: Deliver user value as soon as possible
- Risk early: Tackle unknowns and innovations early
- Dependencies: Respect technical and functional dependencies
- MVP completeness: Ensure Phase 1 delivers complete MVP from Step 3

ADAPTATION RULES:

- Startup MVP: 1-2 phases max, focus on proving concept
- Enterprise product: 4-6 phases, include pilot and rollout
- Internal tool: 2-3 phases, include training and adoption
- Complex domain: Add validation/compliance phases

Connect to scope:

- MVP scope → Phase 1
- Growth features → Phase 2-3
- Vision features → Phase 3-4

Keep it high-level and flexible:
"This roadmap provides strategic direction. Specific timelines will depend on team velocity and priorities."</action>

<template-output>roadmap_phases</template-output>
</step>

<step n="11" goal="Complete PRD and suggest next steps">
<template-output>product_value_summary</template-output>

<check if="standalone_mode != true">
  <action>Load the FULL file: {status_file}</action>
  <action>Update workflow_status["prd"] = "{default_output_file}"</action>
  <action>Save file, preserving ALL comments and structure</action>
</check>

<output>**✅ PRD Complete, {user_name}!**

Your comprehensive product requirements document is ready, including:

**Created:**

- **PRD.md** - Complete requirements document including:
  - Vision & Strategic Objectives
  - Business Model
  - User Personas & Customer Journeys
  - Functional & Non-Functional Requirements
  - High-Level Architecture
  - Epic Breakdown with User Stories
  - Roadmap

All adapted to {project_type} and {domain}.

**Next Steps:**

1. **UX Design** (If UI exists)
   Run: `workflow ux-design` for detailed interaction design and prototyping

2. **Detailed Architecture** (Recommended)
   Run: `workflow create-architecture` for in-depth technical architecture decisions

3. **Story Refinement**
   Use epics and stories from the PRD to start sprint planning and implementation

What makes your product special - {product_value_summary} - is woven throughout the PRD and will guide all design and development work.
</output>
</step>

</workflow>
