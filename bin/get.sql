SELECT
    host,
    remote_id,
    series
FROM
    bookshelf_data_archive
WHERE
    (host, remote_id) IN (
        SELECT
            host,
            MAX(remote_id)
        FROM
            bookshelf_data_archive
        WHERE
            host = 'flos-node-1' AND
            series LIKE '/home/kaz/%'
        GROUP BY
            series
    )
;
