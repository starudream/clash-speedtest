package clash

import (
	"fmt"
)

type Version struct {
	Version string `json:"version"`
	Premium bool   `json:"premium,omitempty"` // clash
	Meta    bool   `json:"meta,omitempty"`    // clash.Meta
}

func (c *Client) GetVersion() (version *Version, _ error) {
	resp, err := c.R().SetResult(&Version{}).Get(c.Addr + "/version")
	if err != nil {
		return version, err
	}

	if !resp.IsSuccess() {
		return version, fmt.Errorf("response status %s", resp.Status())
	}

	return resp.Result().(*Version), nil
}
