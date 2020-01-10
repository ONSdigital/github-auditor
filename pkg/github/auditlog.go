package github

import (
	"github.com/ONSdigital/graphql"
	"github.com/pkg/errors"
)

type (

	// Node represents a node in the returned results graph.
	Node struct {
		ID                   string `json:"id"`
		Type                 string `json:"-"`
		Action               string `json:"action"`
		ActorLogin           string `json:"actorLogin"`
		CreatedAt            string `json:"createdAt"`
		OauthApplicationName string `json:"oauthApplicationName,omitempty"`
		OrganizationName     string `json:"organizationName,omitempty"`
		UserLogin            string `json:"userLogin,omitempty"`
		RepositoryName       string `json:"repositoryName,omitempty"`
		TeamName             string `json:"teamName,omitempty"`
		Visibility           string `json:"visibility,omitempty"`
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
							actorLogin
							createdAt
							userLogin
						}
						... on OauthApplicationCreateAuditEntry {
							action
							actorLogin
							createdAt
							oauthApplicationName
							organizationName
						}
						... on RepoAccessAuditEntry {
							action
							actorLogin
							createdAt
							repositoryName
							visibility
						}
						... on RepoAddMemberAuditEntry {
							action
							actorLogin
							createdAt
							repositoryName
						}
						... on RepoArchivedAuditEntry {
							action
							actorLogin
							createdAt
							repositoryName
						}
						... on RepoCreateAuditEntry {
							action
							actorLogin
							createdAt
							repositoryName
							visibility
						}
						... on RepoDestroyAuditEntry {
							action
							actorLogin
							createdAt
							repositoryName
						}
						... on RepoRemoveMemberAuditEntry {
							action
							actorLogin
							createdAt
							repositoryName
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
