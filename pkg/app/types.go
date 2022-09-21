package app

import "time"

const (
	createSchema = `
CREATE TABLE IF NOT EXISTS history (
    job_name varchar,
	timestamp integer NOT NULL,
	check_command varchar NULL,
	check_command_output varchar NULL
);

CREATE INDEX IF NOT EXISTS history__timestamp ON history(timestamp);
CREATE INDEX IF NOT EXISTS history__check_command ON history(check_command);
`
	getLastHistory = `
SELECT rowid, job_name, timestamp, check_command, check_command_output
FROM history
WHERE check_command = ?
ORDER BY timestamp DESC
LIMIT 1
`
	addToHistory = `
INSERT INTO history (
	job_name, timestamp, check_command, check_command_output    
) VALUES (
    ?, ?, ?, ?
)
`
	pruneHistory = `
DELETE FROM history 
WHERE rowid < (SELECT max(rowid) FROM history LIMIT 1) - 1024;
`
)

type History struct {
	RowId              int64
	jobName            string
	Timestamp          time.Time
	CheckCommand       string
	CheckCommandOutput string
}
