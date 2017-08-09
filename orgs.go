package cf

import "io/ioutil"

func (c *Client) Orgs() ([]byte, error) {
	resp, err := c.httpClient.Get(c.Endpoint + "/v2/organizations")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
