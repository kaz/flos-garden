package bookshelf

import "context"

func RunArchiveCollector(ctx context.Context, host string) error {
	archiveCollector, err := newBookshelfCollector("archive", host, "/archive/snapshots")
	if err != nil {
		return err
	}

	go archiveCollector.Collect(ctx)
	return nil
}
