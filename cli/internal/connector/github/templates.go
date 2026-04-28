package github

// GraphQL queries and mutations used by the github connector. Lifted from the
// original `.archetipo/connectors/github.md` to preserve behaviour and tested
// against `gh api graphql` semantics.

// linkProjectMutation links a Projects v2 project to a repository so the
// project shows up in repo-context tools.
const linkProjectMutation = `
mutation($projectId: ID!, $repoId: ID!) {
  linkProjectV2ToRepository(input: {projectId: $projectId, repositoryId: $repoId}) {
    repository { id nameWithOwner }
  }
}`

// addProjectItemMutation adds an issue to a project board and returns the new item id.
const addProjectItemMutation = `
mutation($projectId: ID!, $contentId: ID!) {
  addProjectV2ItemById(input: {projectId: $projectId, contentId: $contentId}) {
    item { id }
  }
}`

// updateSingleSelectFieldMutation updates a single-select field for one item.
const updateSingleSelectFieldMutation = `
mutation($projectId: ID!, $itemId: ID!, $fieldId: ID!, $optionId: String!) {
  updateProjectV2ItemFieldValue(input: {
    projectId: $projectId,
    itemId: $itemId,
    fieldId: $fieldId,
    value: { singleSelectOptionId: $optionId }
  }) { projectV2Item { id } }
}`

// updateNumberFieldMutation updates a number field for one item.
const updateNumberFieldMutation = `
mutation($projectId: ID!, $itemId: ID!, $fieldId: ID!, $value: Float!) {
  updateProjectV2ItemFieldValue(input: {
    projectId: $projectId,
    itemId: $itemId,
    fieldId: $fieldId,
    value: { number: $value }
  }) { projectV2Item { id } }
}`
