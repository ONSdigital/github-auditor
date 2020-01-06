package github

import (
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

type (

	// AuditEntry represents an entry in the GitHub audit log.
	AuditEntry struct {
		ID         string `json:"id"`
		Action     string `json:"action"`
		ActorLogin string `json:"actorLogin"`
		CreatedAt  string `json:"createdAt"`
		UserLogin  string `json:"userLogin,omitempty"`
	}

	// Organization represents a GitHub organisation.
	Organization struct {
		AuditLog struct {
			TotalCount int
			PageInfo   PageInfo
			Nodes      []AuditEntry
		}
	}
)

// FetchAllAuditLogEntries returns all audit log entries for the specified organisation.
func (c *Client) FetchAllAuditLogEntries(organisation string) ([]AuditEntry, error) {
	var auditEntries []AuditEntry
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
					}
				}
			}
		}
	`)

	req.Var("login", organisation)

	page := 0
	hasNextPage := true

	for hasNextPage {
		page++
		res := &struct{ Organization Organization }{}
		req.Var("after", endCursor)

		if err := c.Run(req, &res); err != nil {
			return nil, errors.Wrap(err, "failed to fetch audit log entries for organisation")
		}

		auditEntries = append(auditEntries, res.Organization.AuditLog.Nodes...)
		endCursor = &res.Organization.AuditLog.PageInfo.EndCursor
		hasNextPage = res.Organization.AuditLog.PageInfo.HasNextPage
	}

	return auditEntries, nil
}
