package github

import (
	"github.com/ONSdigital/graphql"
	"github.com/pkg/errors"
)

type (

	// Node represents a node in the returned results graph.
	Node struct {
		ID             string `json:"id"`
		Action         string `json:"action"`
		ActorLogin     string `json:"actorLogin"`
		CreatedAt      string `json:"createdAt"`
		UserLogin      string `json:"userLogin,omitempty"`
		RepositoryName string `json:"repositoryName,omitempty"`
		TeamName       string `json:"teamName,omitempty"`
		Visibility     string `json:"visibility,omitempty"`
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

// FetchAllAuditLogEntries returns all audit log entries for the specified organisation.
func (c Client) FetchAllAuditLogEntries(organisation string) ([]Node, error) {
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
						# An entry in the audit log.
						... on AuditEntry {
							action
							actorLogin
							createdAt
							userLogin							
						}
						# Triggered when a repository owned by an organisation is switched from private to public (or vice-versa).
						... on RepoAccessAuditEntry { 
							action
							actorLogin
							createdAt
							userLogin
							repositoryName
							visibility						
						}
						# Triggered when a user accepts an invitation to have collaboration access to a repository.
						... on RepoAddMemberAuditEntry {
							action
							actorLogin
							createdAt
							userLogin
							repositoryName
						}
						# Triggered when a new repository is created.
						... on RepoCreateAuditEntry {
							action
							actorLogin
							createdAt
							userLogin
							repositoryName
						}
						# Triggered when a repository is deleted.
						... on RepoDestroyAuditEntry {
							action
							actorLogin
							createdAt
							userLogin
							repositoryName
						}						
						# Triggered when a team is given control of a repository.
						... on TeamAddRepositoryAuditEntry {
							action
							actorLogin
							createdAt
							userLogin
							repositoryName
							teamName
						}
						# Triggered when a repository is no longer under a team's control.
						... on TeamRemoveRepositoryAuditEntry {
							action
							actorLogin
							createdAt
							userLogin
							repositoryName
							teamName
						}											
					}
				}
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
