package cmd

import "context"

type IUser interface {
	Register(ctx context.Context) error
	Login(ctx context.Context) error
}
