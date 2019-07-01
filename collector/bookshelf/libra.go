package bookshelf

import "context"

func RunLibraCollector(ctx context.Context, host string) error {
	libraColelctor, err := newBookshelfCollector("libra", host, "/libra/books", "LONGTEXT")
	if err != nil {
		return err
	}

	go libraColelctor.Collect(ctx)
	return nil
}
