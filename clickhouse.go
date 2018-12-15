package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mailru/go-clickhouse"
)

// Debug: select * from snmp.table_buffer where Name='vlan2392' \G

// PushDataToDB - pushes SNMP data to local clickhouse instance
func PushDataToDB(address string, port int, table string, data map[MergedTableMapKey][]IfMergedTable) error {

	// open connection to local clickhouse instance
	connect, err := sql.Open("clickhouse", fmt.Sprintf("http://%s:%d", address, port))
	if err != nil {
		return fmt.Errorf("failed to push SNMP data to DB: %s", err.Error())
	}

	// check server availability
	if err = connect.Ping(); err != nil {
		return fmt.Errorf("failed to push SNMP data to DB: %s", err.Error())
	}

	// begin SQL transaction
	tx, err := connect.Begin()
	if err != nil {
		return fmt.Errorf("failed to push SNMP data to DB: %s", err.Error())
	}

	// define SQL statement template
	sql := `
	INSERT INTO ###TABLE_NAME### (
		Address,
		Community,
		ElapsedSeconds,
		ElapsedSecondsSinceEpoch,
		Index,
		Name,
		Descr,
		Alias,
		Type,
		Mtu,
		Speed,
		HighSpeed,
		PhysAddress,
		AdminStatus,
		OperStatus,
		InOctets,
		InUcastPkts,
		InDiscards,
		InErrors,
		InUnknownProtos,
		OutOctets,
		OutUcastPkts,
		OutDiscards,
		OutErrors,
		InMulticastPkts,
		InBroadcastPkts,
		OutMulticastPkts,
		OutBroadcastPkts,
		HCInOctets,
		HCInUcastPkts,
		HCInMulticastPkts,
		HCInBroadcastPkts,
		HCOutOctets,
		HCOutUcastPkts,
		HCOutMulticastPkts,
		HCOutBroadcastPkts,
		PromiscuousMode,
		ConnectorPresent
	) VALUES (
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?
	)`

	// replace table name inside SQL template, `table` is validated during application configuration
	sql = strings.Replace(sql, "###TABLE_NAME###", strings.TrimSpace(table), 1)

	// prepare SQL statement
	stmt, err := tx.Prepare(sql) // nolint: safesql
	if err != nil {
		return fmt.Errorf("failed to push SNMP data to DB: %s", err.Error())
	}

	// batch insert data into prepared statement
	for k := range data {
		for i := range data[k] {
			if _, err := stmt.Exec(
				data[k][i].Address,
				data[k][i].Community,
				data[k][i].ElapsedSeconds,
				data[k][i].ElapsedSecondsSinceEpoch,
				data[k][i].Index,
				data[k][i].Name,
				data[k][i].Descr,
				data[k][i].Alias,
				data[k][i].Type,
				data[k][i].Mtu,
				data[k][i].Speed,
				data[k][i].HighSpeed,
				data[k][i].PhysAddress,
				data[k][i].AdminStatus,
				data[k][i].OperStatus,
				data[k][i].InOctets,
				data[k][i].InUcastPkts,
				data[k][i].InDiscards,
				data[k][i].InErrors,
				data[k][i].InUnknownProtos,
				data[k][i].OutOctets,
				data[k][i].OutUcastPkts,
				data[k][i].OutDiscards,
				data[k][i].OutErrors,
				data[k][i].InMulticastPkts,
				data[k][i].InBroadcastPkts,
				data[k][i].OutMulticastPkts,
				data[k][i].OutBroadcastPkts,
				data[k][i].HCInOctets,
				data[k][i].HCInUcastPkts,
				data[k][i].HCInMulticastPkts,
				data[k][i].HCInBroadcastPkts,
				data[k][i].HCOutOctets,
				data[k][i].HCOutUcastPkts,
				data[k][i].HCOutMulticastPkts,
				data[k][i].HCOutBroadcastPkts,
				data[k][i].PromiscuousMode,
				data[k][i].ConnectorPresent,
			); err != nil {
				return fmt.Errorf("failed to push SNMP data to DB: %s", err.Error())
			}
		}
	}

	// commit prepared statement
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to push SNMP data to DB: %s", err.Error())
	}

	return nil
}
