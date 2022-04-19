package client

import (
	"fmt"
	"net/url"
)

type Project struct {
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	ID               string  `json:"id"`
	OrganizationName string  `json:"organization_name"`
	OrganizationSlug string  `json:"organization_slug"`
	OrganizationID   string  `json:"organization_id"`
	VCSInfo          VCSInfo `json:"vcs_info"`
}

type VCSInfo struct {
	URL           string `json:"vcs_url"`
	Provider      string `json:"provider"`
	DefaultBranch string `json:"default_branch"`
}

// GetProject gets an existing project by its project slug (vcs-slug/org-name/repo-name)
func (c *Client) GetProject(project string) (*Project, error) {
	slug, err := c.Slug(project)
	if err != nil {
		return nil, err
	}

	req, err := c.rest.NewRequest("GET", &url.URL{Path: fmt.Sprintf("project/%s", slug)}, nil)
	if err != nil {
		return nil, err
	}

	p := &Project{}

	if sts, err := c.rest.DoRequest(req, p); err != nil {
		return nil, fmt.Errorf("could not find project: %v %v", sts, err)
	}

	return p, nil
}
