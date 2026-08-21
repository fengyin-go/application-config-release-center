package client

type Client struct {
	Called chan struct{}
	Calls  int
}

func New() *Client      { return &Client{Called: make(chan struct{}, 3)} }
func (c *Client) Send() { c.Calls++; c.Called <- struct{}{} }
