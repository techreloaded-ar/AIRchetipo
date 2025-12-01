# {{project_name}} - Product Backlog

**Author:** Archetipo
**Date:** {{date}}

---

## Overview

This document provides the complete product backlog for {{project_name}}, with the decomposition of requirements from the [PRD](./PRD.md) into implementable epics and user stories.

**Important Note:** User stories are ordered by priority - the most important ones for the MVP are positioned at the top of the backlog.

{{backlog_summary}}

**Statistics:**
- Total Epics: {{epic_count}}
- Total User Stories: {{story_count}}

---

## Functional Requirements Inventory

{{fr_inventory}}

---

## FR Coverage Map

{{fr_coverage_map}}

---

<!-- Repeat for each epic (N = 1, 2, 3...) -->

## Epic {{N}}: {{epic_title_N}}

**Goal:**

As a {{epic_user_type_N}},
I want {{epic_capability_N}},
So that {{epic_value_N}}.

**Covered FRs:** {{epic_fr_list_N}}

---

<!-- Repeat for each story (M = 1, 2, 3...) within epic N -->

### Story {{N}}.{{M}}: {{story_title_N_M}}

**User Story:**

As a {{user_type_N_M}},
I want {{capability_N_M}},
So that {{value_benefit_N_M}}.

**Acceptance Criteria (complete at least two; three recommended):**

**Criterion 1**
- **Given** {{precondition_1_N_M}}
- **When** {{action_1_N_M}}
- **Then** {{expected_outcome_1_N_M}}

**Criterion 2**
- **Given** {{precondition_2_N_M}}
- **When** {{action_2_N_M}}
- **Then** {{expected_outcome_2_N_M}}

**Criterion 3**
- **Given** {{precondition_3_N_M}}
- **When** {{action_3_N_M}}
- **Then** {{expected_outcome_3_N_M}}

**Test Scenarios:**

{{test_scenarios_N_M}}

**Dependencies:** {{dependencies_N_M}}

**Technical Notes:** {{technical_notes_N_M}}

---

<!-- End story repetition -->

<!-- End epic repetition -->

---

## FR Coverage Matrix

{{fr_coverage_matrix}}

---

## Summary

{{backlog_final_summary}}

---

_This document will be updated after the UX Design and Architecture workflows to incorporate interaction details and technical decisions._
