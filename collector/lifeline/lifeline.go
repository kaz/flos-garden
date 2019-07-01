package lifeline

import (
	"context"
)

func RunLifelineCollector(ctx context.Context, host string) error {
	lifelineCollector, err := newCollector(ctx, host)
	if err != nil {
		return err
	}

	go lifelineCollector.Collect()
	return nil
}
