# PRD Information Checklist

This checklist guides the party-mode agents to collect all necessary information for generating a comprehensive Product Requirements Document.

## Information Collection Status

### 1. Vision & Strategy

- [ ] **Vision Statement** - One compelling sentence capturing the future state
- [ ] **Strategic Objectives** - 3-5 key objectives that move toward the vision
- [ ] **Long-term Impact** - The transformational change this product will create
- [ ] **Product Differentiator** - What makes this product special and unique

### 2. Business Model

- [ ] **Value Proposition** - What unique value does this create?
- [ ] **Customer Segments** - Who are the target customers?
- [ ] **Revenue Model** - How will this generate revenue/value?
- [ ] **Key Resources** - What critical resources are needed?
- [ ] **Cost Structure** - What are the main cost drivers?

### 3. Target Users (Personas)

- [ ] **Persona 1 - Name and Role** (fictional but representative)
- [ ] **Persona 1 - Background** - Context, experience level, environment
- [ ] **Persona 1 - Goals** - What they want to achieve
- [ ] **Persona 1 - Pain Points** - Current frustrations and challenges
- [ ] **Persona 1 - Behaviors** - How they work, what tools they use
- [ ] **Persona 1 - Motivations** - What drives them
- [ ] **Persona 1 - Tech Savviness** - Their comfort level with technology

- [ ] **Persona 2 - Name and Role**
- [ ] **Persona 2 - Background**
- [ ] **Persona 2 - Goals**
- [ ] **Persona 2 - Pain Points**
- [ ] **Persona 2 - Behaviors**
- [ ] **Persona 2 - Motivations**
- [ ] **Persona 2 - Tech Savviness**

### 4. Customer Journey

- [ ] **Persona 1 Journey - Awareness** - How do they discover the product?
- [ ] **Persona 1 Journey - Consideration** - What makes them consider using it?
- [ ] **Persona 1 Journey - First Use** - Their initial experience
- [ ] **Persona 1 Journey - Regular Use** - How they engage day-to-day
- [ ] **Persona 1 Journey - Advocacy** - What makes them recommend it?

- [ ] **Persona 2 Journey** - Same stages as Persona 1

### 5. Success Criteria

- [ ] **Success Criteria** - What winning looks like for THIS specific product
- [ ] **Business Metrics** (if applicable) - Measurable business outcomes

### 6. Product Scope

- [ ] **MVP Scope** - What must work for this to be useful?
- [ ] **Growth Features** - What makes it competitive?
- [ ] **Vision Features** - What's the dream version?

### 7. Project Classification

- [ ] **Project Type** - API, Mobile, Web, CLI, SDK, SaaS, etc. (from project-types.csv)
- [ ] **Domain Type** - Medical, Finance, Government, Education, etc. (from domain-complexity.csv)
- [ ] **Complexity Level** - Low, Medium, High, Critical
- [ ] **Domain Context Summary** (if complex domain)

### 8. Domain-Specific Requirements (if applicable)

- [ ] **Regulatory Requirements**
- [ ] **Compliance Needs**
- [ ] **Industry Standards**
- [ ] **Safety/Risk Factors**
- [ ] **Required Validations**

### 9. Innovation Patterns (if applicable)

- [ ] **Innovation Description** - What makes it unique/novel?
- [ ] **Validation Approach** - How to validate the innovation?

### 10. Project-Specific Requirements

Based on project type, collect relevant information:

**For API/Backend:**

- [ ] Endpoint specification
- [ ] Authentication model
- [ ] Authorization model
- [ ] Error codes and rate limits
- [ ] Data schemas

**For Mobile:**

- [ ] Platform requirements (iOS/Android/both)
- [ ] Device features needed
- [ ] Offline capabilities
- [ ] Store compliance

**For SaaS B2B:**

- [ ] Multi-tenant architecture
- [ ] Permission models
- [ ] Subscription tiers
- [ ] Critical integrations

**For Web Applications:**

- [ ] Browser support
- [ ] Responsive requirements
- [ ] Progressive Web App features

### 11. UX Principles (if UI exists)

- [ ] **Visual Personality** - How should this feel to use?
- [ ] **Key Interaction Patterns**
- [ ] **Critical User Flows**

### 12. Functional Requirements

- [ ] **Complete FR List** - Comprehensive list of all capabilities
- [ ] **FR Organization** - Grouped by capability area
- [ ] **FR Coverage** - All capabilities from vision, scope, and domain requirements represented

Functional Requirements must include:

- User-facing capabilities
- System capabilities
- Integration capabilities
- Data management capabilities
- Security and access capabilities

### 13. Non-Functional Requirements

Collect only those that matter for THIS product:

- [ ] **Performance Requirements** (if user-facing impact)
- [ ] **Security Requirements** (if handling sensitive data)
- [ ] **Scalability Requirements** (if growth expected)
- [ ] **Accessibility Requirements** (if broad audience)
- [ ] **Integration Requirements** (if connecting systems)

### 14. High-Level Architecture

- [ ] **Overall System Architecture** - Monolith, microservices, serverless, etc.
- [ ] **Key Components** - Main components and their relationships
- [ ] **Technology Stack** - Primary languages, frameworks, database
- [ ] **Database Architecture** - Type, technology, data modeling approach
- [ ] **Frameworks and Libraries** - Core frameworks for each layer
- [ ] **Infrastructure Overview** (optional) - Cloud provider, hosting approach

### 15. Epic Breakdown

- [ ] **Epics Overview** - Summary of all epics
- [ ] **Epic Details** - For each epic:
  - Epic title and goal
  - Which FRs it implements
  - 5-15 high-level user stories
  - Dependencies on other epics
  - Estimated complexity (Small/Medium/Large)
  - Priority (MVP/Growth/Vision)

### 16. Roadmap

- [ ] **Phase 1 (Q1)** - Foundation & MVP
- [ ] **Phase 2 (Q2)** - Growth Features
- [ ] **Phase 3 (Q3)** - Scale & Optimization
- [ ] **Phase 4 (Q4)** - Vision Features

For each phase:

- Goals
- Epics included
- Key milestones
- Success metrics
- Dependencies

---

## Completeness Gate

Before generating the PRD, verify:

1. ✅ **Minimum Required Information** (Cannot proceed without these):
   - Vision statement
   - At least 1 persona with complete profile
   - Product scope (at least MVP defined)
   - Project classification (type, domain, complexity)
   - At least 10 functional requirements
   - High-level architecture

2. ✅ **Recommended Information** (Should have most of these):
   - 2 personas with customer journeys
   - Business model
   - Success criteria
   - Project-specific requirements
   - Non-functional requirements (relevant ones)
   - Epic breakdown
   - Roadmap

3. ✅ **Optional Information** (Nice to have):
   - Domain-specific requirements
   - Innovation patterns
   - UX principles
   - Infrastructure overview

## Agent Conversation Flow

### Phase 1: Discovery (Steps 1-7 of checklist)

Focus on understanding the product vision, users, and business context.
Questions should be open-ended and exploratory.

### Phase 2: Requirements (Steps 8-13 of checklist)

Focus on defining what the product must do and how well it must do it.
Questions should be specific and technical.

### Phase 3: Planning (Steps 14-16 of checklist)

Focus on architecture, implementation strategy, and timeline.
Questions should be about feasibility and approach.

### Phase 4: Validation

Review all collected information for completeness and clarity.
Fill any gaps with targeted questions.

### Phase 5: Generation

Synthesize all information and generate the PRD document.
