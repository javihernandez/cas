package api

import (
	"context"
	"fmt"
	"github.com/vchain-us/ledger-compliance-go/schema"
)

func (u LcUser) RequireFeatOrErr(feat string) error {
	f, err := u.Client.Feats(context.Background())
	if err != nil {
		return err
	}
	if _, ok := f.Map()[feat]; !ok {
		return fmt.Errorf("seems that the connected server component `%s` at version `%s` builded at `%s` doesn't support %s feature. Please contact a system administrator", f.Component, f.Version, f.BuildTime, schema.FeatBoM)
	}
	return nil
}
