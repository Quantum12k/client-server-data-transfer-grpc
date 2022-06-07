package main

import "context"

type Authentication struct {
	Login    string
	Password string
}

func (c *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    c.Login,
		"password": c.Password,
	}, nil
}

func (c *Authentication) RequireTransportSecurity() bool {
	return true
}
