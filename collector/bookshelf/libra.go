package bookshelf

import "context"

func RunLibraCollector(ctx context.Context, host string) error {
	libraColelctor, err := newBookshelfCollector("libra", host, "/libra/books")
	if err != nil {
		return err
	}

	go libraColelctor.Collect(ctx)
	return nil
}
