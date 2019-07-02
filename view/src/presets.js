export default [
  {
    name: "File access log (recent 1000 entries)",
    sql: `
SELECT
  host,
  BINARY created AS time,
  contents
FROM
  bookshelf_data_libra
WHERE
  series = 'audit'
ORDER BY
  created DESC
    `.trim(),
  },
  {
    name: "File snapshot log (recent 1000 entries)",
    sql: `
SELECT
  host,
  BINARY created AS time,
  contents
FROM
  bookshelf_data_libra
WHERE
  series = 'archive'
ORDER BY
  created DESC
    `.trim(),
  },
  {
    name: "Log tails (recent 1000 entries)",
    sql: `
SELECT
  host,
  BINARY created AS time,
  contents
FROM
  bookshelf_data_libra
WHERE
  series = 'tail'
ORDER BY
  created DESC
    `.trim(),
  },
  {
    name: "Monitoring results (recent 1000 entries)",
    sql: `
SELECT
  host,
  BINARY created AS time,
  contents
FROM
  bookshelf_data_libra
WHERE
  series = 'lifeline'
ORDER BY
  created DESC
    `.trim(),
  },
];
