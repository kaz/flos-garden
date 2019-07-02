package lifeline

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kaz/flos-hortus/common"
	"github.com/kaz/flos-hortus/database"
	"github.com/kaz/flos/lifeline"
)

type (
	collector struct {
		*common.Collector
	}
)

func newCollector(ctx context.Context, host string) (*collector, error) {
	if _, err := database.DB().Exec("CREATE TABLE IF NOT EXISTS lifeline_data (host TEXT, name TEXT, success BOOL, output MEDIUMTEXT, updated DATETIME(6), PRIMARY KEY(host(128), name(128)))" + common.TABLE_OPTION); err != nil {
		return nil, err
	}

	c := &collector{
		&common.Collector{
			Context: ctx,
			Logger:  log.New(os.Stdout, "[lifeline host="+host+"] ", log.Ltime),
			Host:    host,
			Path:    "/lifeline",
		},
	}

	c.RegisterCollectFunc(c.collect)
	return c, nil
}

func (c *collector) collect() error {
	resp, err := c.DoRequest(http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to get books: %v\n", msg)
	}

	var results map[string]*lifeline.Result
	if err := common.ReadBody(resp, &results); err != nil {
		return fmt.Errorf("failed to decode resp body: %v\n", err)
	}

	stmt, err := database.DB().PrepareNamedContext(c.Context, "REPLACE INTO lifeline_data VALUES (:host, :name, :success, :output, :updated)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v\n", err)
	}

	for _, result := range results {
		_, err := stmt.ExecContext(c.Context, map[string]interface{}{
			"host":    c.Host,
			"name":    result.Name,
			"success": result.Success,
			"output":  result.Output,
			"updated": common.Time(result.Timestamp),
		})
		if err != nil {
			return fmt.Errorf("failed to insert record: %v\n", err)
		}
	}

	return nil
}
