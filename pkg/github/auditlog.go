package github

import (
	"github.com/ONSdigital/graphql"
	"github.com/pkg/errors"
)

type (

	// Actor represents the GitHub user who initiated the audit action.
	Actor struct {
		Type  string `json:"__typename"`
		Login string `json:"login,omitempty"` // Bot
		Name  string `json:"name,omitempty"`  // Organization or User
	}

	// Node represents a node in the returned results graph.
	Node struct {
		ID                   string `json:"id"`
		Action               string `json:"action"`
		Actor                Actor
		BlockedUser          Actor
		CreatedAt            string `json:"createdAt"`
		MergeType            string `json:"mergeType,omitempty"`
		OauthApplicationName string `json:"oauthApplicationName,omitempty"`
		OrganizationName     string `json:"organizationName,omitempty"`
		Permission           string `json:"permission,omitempty"`
		PermissionWas        string `json:"permissionWas,omitempty"`
		RepositoryName       string `json:"repositoryName,omitempty"`
		TeamName             string `json:"teamName,omitempty"`
		User                 Actor
		Visibility           string `json:"visibility,omitempty"`
	}

	// PageInfo represents the pagination information returned from the query.
	PageInfo struct {
		StartCursor     string
		EndCursor       string
		HasPreviousPage bool
		HasNextPage     bool
	}

	// Organization represents a GitHub organisation.
	Organization struct {
		AuditLog struct {
			TotalCount int
			PageInfo   PageInfo
			Nodes      []Node
		}
	}
)

// FetchAllAuditEvents returns all audit log events for the passed organisation.
func (c Client) FetchAllAuditEvents(organisation string) (events []Node, err error) {
	var endCursor *string // Using a pointer type allows this to be nil (an empty string isn't a valid cursor).

	req := graphql.NewRequest(`
		query GitHubAuditEntries($login: String!, $after: String) {
			organization(login: $login) {
				auditLog(first: 100, after: $after) {
					totalCount
					pageInfo {
						startCursor
						endCursor
						hasNextPage
						hasPreviousPage
					}
					nodes {
						... on Node {
							id
						}
						... on AuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							user {
								...userFields
							}
						}
						... on OauthApplicationCreateAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							oauthApplicationName
							organizationName
						}
						... on OrgAddBillingManagerAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
							user {
								...userFields
							}
						... on OrgAddMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							user {
								...userFields
							}
						}
						... on OrgBlockUserAuditEntry {
							action
							actor {
								...actorFields
							}
							blockedUser
							createdAt
							organizationName
						}
						... on OrgCreateAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
						}
						... on OrgDisableSamlAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
						}
						... on OrgDisableTwoFactorRequirementAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
						}
						... on OrgEnableOauthAppRestrictionsAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
						}
						... on OrgEnableSamlAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
						}
						... on OrgEnableTwoFactorRequirementAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
						}
						... on OrgInviteMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
							user {
								...userFields
							}
						}
						... on OrgOauthAppAccessApprovedAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							oauthApplicationName
							organizationName
						}
						... on OrgOauthAppAccessDeniedAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							oauthApplicationName
							organizationName
						}
						... on OrgOauthAppAccessRequestedAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							oauthApplicationName
							organizationName
						}
						... on OrgRemoveBillingManagerAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
							user {
								...userFields
							}
						}
						... on OrgRemoveMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
							user {
								...userFields
							}
						}
						... on OrgRemoveOutsideCollaboratorAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
							user {
								...userFields
							}
						}
						... on OrgUpdateMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							organizationName
							permission
							permissionWas
							user {
								...userFields
							}
						}
						... on RepoAccessAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
							visibility
						}
						... on RepoAddMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
							user {
								...userFields
							}
						}
						... on RepoArchivedAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
						}
						... on RepoChangeMergeSettingAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							mergeType
							repositoryName
						}
						... on RepoCreateAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
							visibility
						}
						... on RepoDestroyAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
						}
						... on RepoRemoveMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
							user {
								...userFields
							}
						}
						... on TeamAddMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							teamName
							user {
								...userFields
							}
						}
						... on TeamAddRepositoryAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
							teamName
						}
						... on TeamRemoveMemberAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							teamName
							user {
								...userFields
							}
						}
						... on TeamRemoveRepositoryAuditEntry {
							action
							actor {
								...actorFields
							}
							createdAt
							repositoryName
							teamName
						}
					}
				}
			}
		}
		fragment userFields on User {
			__typename
			login
			name
		}
		fragment actorFields on Actor {
			__typename
			... on Bot {
				login
			}
			... on User {
				...userFields
			}
			... on Organization {
				name
			}
		}
	`)

	req.Var("login", organisation)

	page := 0
	hasNextPage := true
	var nodes []Node

	for hasNextPage {
		page++
		res := &struct{ Organization Organization }{}
		req.Var("after", endCursor)

		if err := c.Run(req, &res); err != nil {
			return nil, errors.Wrap(err, "failed to fetch audit log entries for organisation")
		}

		nodes = append(nodes, res.Organization.AuditLog.Nodes...)
		endCursor = &res.Organization.AuditLog.PageInfo.EndCursor
		hasNextPage = res.Organization.AuditLog.PageInfo.HasNextPage
	}

	return nodes, nil
}
