package bookshelf

import "context"

func RunArchiveCollector(ctx context.Context, host string) error {
	archiveCollector, err := newBookshelfCollector("archive", host, "/archive/snapshots", "LONGBLOB")
	if err != nil {
		return err
	}

	go archiveCollector.Collect(ctx)
	return nil
}
