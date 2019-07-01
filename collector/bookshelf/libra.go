package bookshelf

import "context"

func RunLibraCollector(ctx context.Context, host string) error {
	libraColelctor, err := newBookshelfCollector(ctx, "libra", host, "/libra/books", "LONGTEXT")
	if err != nil {
		return err
	}

	go libraColelctor.Collect()
	return nil
}
