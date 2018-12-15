package main

func main() {

	// declare variable for local DB
	after := make(map[MergedTableMapKey][]IfMergedTable)

	// declare variable for printing
	delta := make(map[MergedTableMapKey][]IfMergedTable)

	// declare function to fush data into local DB
	flushLocalDB := func(after map[MergedTableMapKey][]IfMergedTable) {
		if err := WriteDataToFile(after, c.DB.Local.Path); err != nil {
			ErrorLogger.Fatal(err)
		}
	}

	// flush data to local DB on exit
	defer flushLocalDB(after)

	// read previous data from local DB
	before, err := ReadDataFromFile(c.DB.Local.Path)
	if err != nil {
		ErrorLogger.Print(err)
		return
	}

	// loop-over devices in config
	for _, key := range c.SNMP {

		var table []IfMergedTable

		// get table data
		table, err = GetMergedTableData(key.Community, key.Address)
		if err != nil {
			InfoLogger.Print(err)
			continue
		}

		// set table data for local DB
		after[key] = table

		// get previous table data
		if _, ok := before[key]; !ok {
			ErrorLogger.Printf("no previous table data for SNMP community='%s' on device='%s'", key.Community, key.Address)
			continue
		}

		// compute delta table
		delta[key] = DeltaIfMergedTables(before[key], table)
	}

	// send data to DB
	err = PushDataToDB(
		c.DB.Clickhouse.Address,
		c.DB.Clickhouse.Port,
		c.DB.Clickhouse.Table,
		delta,
	)
	if err != nil {
		ErrorLogger.Print(err)
		return
	}
}
