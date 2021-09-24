package shortcutclient

import "fmt"

type Member struct {
	Disabled bool    `json:"disabled"`
	ID       string  `json:"id"`
	Profile  Profile `json:"profile"`
	Role     string  `json:"role"`
	State    string  `json:"state"`
}

type Profile struct {
	Deactivated            bool   `json:"deactivated"`
	EmailAddress           string `json:"email_address"`
	EntityType             string `json:"entity_type"`
	ID                     string `json:"id"`
	MentionName            string `json:"mention_name"`
	Name                   string `json:"name"`
	TwoFactorAuthActivated bool   `json:"two_factor	_auth_activated"`
}

func (c *Client) ListMembers() ([]Member, error) {
	path := "/members"

	var members []Member
	if err := c.get(path, &members); err != nil {
		return nil, err
	}

	return members, nil
}
func (c *Client) GetMember(memberID string) (*Member, error) {
	path := fmt.Sprintf("/members/%s", memberID)

	var member *Member
	if err := c.get(path, &member); err != nil {
		return nil, err
	}

	return member, nil
}
