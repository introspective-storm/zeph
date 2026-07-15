package picker

func DataPicker(height int) Model {
	allowedExt := []string{
		".json",
		".ndjson",
		".jsonl",
		".parquet",
		".csv",
		".csv.gz",
		".tsv",
		".txt",
		".xlsx",
		".avro",
		".db",
		".ddb",
		".sql",
		".sqlite",
		".sqlite3",
		".duckdb",
	}
	return New(allowedExt, height)
}
