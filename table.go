package main

import (
	"fmt"
	"sort"
	"time"
)

// IfMergedTable merged IfTable, IfXTable data
type IfMergedTable struct {
	Address   string
	Community string

	ElapsedSeconds           int64
	ElapsedSecondsSinceEpoch int64

	Index int

	Name  string
	Descr string
	Alias string
	Type  string

	Mtu int64

	Speed     int64 // GAUGE (bps)
	HighSpeed int64 // GAUGE (Mbps)

	PhysAddress string
	AdminStatus string
	OperStatus  string

	InOctets        int64 // COUNTER
	InUcastPkts     int64 // COUNTER
	InDiscards      int64 // COUNTER
	InErrors        int64 // COUNTER
	InUnknownProtos int64 // COUNTER
	OutOctets       int64 // COUNTER
	OutUcastPkts    int64 // COUNTER
	OutDiscards     int64 // COUNTER
	OutErrors       int64 // COUNTER

	InMulticastPkts    int64 // COUNTER
	InBroadcastPkts    int64 // COUNTER
	OutMulticastPkts   int64 // COUNTER
	OutBroadcastPkts   int64 // COUNTER
	HCInOctets         int64 // COUNTER
	HCInUcastPkts      int64 // COUNTER
	HCInMulticastPkts  int64 // COUNTER
	HCInBroadcastPkts  int64 // COUNTER
	HCOutOctets        int64 // COUNTER
	HCOutUcastPkts     int64 // COUNTER
	HCOutMulticastPkts int64 // COUNTER
	HCOutBroadcastPkts int64 // COUNTER

	PromiscuousMode  string
	ConnectorPresent string
}

// MergeTables - merges IfTable with IfXTable
func MergeTables(table IfTable, xTable IfXTable) IfMergedTable {

	var out IfMergedTable

	out.Index = table.Index
	out.Descr = table.Descr
	out.Type = table.Type
	out.Mtu = table.Mtu
	out.Speed = table.Speed
	out.PhysAddress = table.PhysAddress
	out.AdminStatus = table.AdminStatus
	out.OperStatus = table.OperStatus
	out.InOctets = table.InOctets
	out.InUcastPkts = table.InUcastPkts
	out.InDiscards = table.InDiscards
	out.InErrors = table.InErrors
	out.InUnknownProtos = table.InUnknownProtos
	out.OutOctets = table.OutOctets
	out.OutUcastPkts = table.OutUcastPkts
	out.OutDiscards = table.OutDiscards
	out.OutErrors = table.OutErrors

	out.Name = xTable.Name
	out.InMulticastPkts = xTable.InMulticastPkts
	out.InBroadcastPkts = xTable.InBroadcastPkts
	out.OutMulticastPkts = xTable.OutMulticastPkts
	out.OutBroadcastPkts = xTable.OutBroadcastPkts
	out.HCInOctets = xTable.HCInOctets
	out.HCInUcastPkts = xTable.HCInUcastPkts
	out.HCInMulticastPkts = xTable.HCInMulticastPkts
	out.HCInBroadcastPkts = xTable.HCInBroadcastPkts
	out.HCOutOctets = xTable.HCOutOctets
	out.HCOutUcastPkts = xTable.HCOutUcastPkts
	out.HCOutMulticastPkts = xTable.HCOutMulticastPkts
	out.HCOutBroadcastPkts = xTable.HCOutBroadcastPkts
	out.HighSpeed = xTable.HighSpeed
	out.PromiscuousMode = xTable.PromiscuousMode
	out.ConnectorPresent = xTable.ConnectorPresent
	out.Alias = xTable.Alias

	return out
}

// GetMergedTableData - aquires merged table output
func GetMergedTableData(community, address string) ([]IfMergedTable, error) {

	// declare output variable
	out := make([]IfMergedTable, 0)

	// declare chanels
	errChan := make(chan error, 2)
	ifTableChan := make(chan map[int]IfTable, 1)
	ifXTableChan := make(chan map[int]IfXTable, 1)

	// SNMP IfTable
	go func(community, address string, ifTableChan chan map[int]IfTable, errChan chan error) {
		ifTable, err := CmdIfTable(community, address)
		if err != nil {
			// send error to error chanel
			errChan <- err
			// close data chanel
			close(ifTableChan)
		}

		// send no-error to error chanel
		errChan <- nil

		// send data to ifTableChan
		ifTableChan <- ifTable
		// close data chanel
		close(ifTableChan)

	}(community, address, ifTableChan, errChan)

	// SNMP IfXTable
	go func(community, address string, ifXTableChan chan map[int]IfXTable, errChan chan error) {
		ifXTable, err := CmdIfXTable(community, address)
		if err != nil {
			// send error to error chanel
			errChan <- err
			// close data chanel
			close(ifXTableChan)
		}

		// send no-error to error chanel
		errChan <- nil

		// send data to ifXTableChan
		ifXTableChan <- ifXTable
		// close data chanel
		close(ifXTableChan)

	}(community, address, ifXTableChan, errChan)

	// read data from chanels
	ifTable := <-ifTableChan
	ifXTable := <-ifXTableChan

	// close error chanel, we already should have send all data to it, ifTableChan and ifXTable are used as locks
	close(errChan)

	// check for errors in errors chanel
	for e := range errChan {
		if e != nil {
			return nil, fmt.Errorf("failed to parse merged table elements for SNMP community='%s' on device='%s', %s", community, address, e.Error())
		}
	}

	// validate length
	if len(ifTable) != len(ifXTable) {
		return nil, fmt.Errorf("failed to parse merged table elements for SNMP community='%s' on device='%s', IfTable elements size does not match IfXTable elements size", community, address)
	}

	// merge IfTable with IfXTable
	for k := range ifTable {
		if _, ok := ifXTable[k]; !ok {
			return nil, fmt.Errorf("failed to parse merged table elements for SNMP community='%s' on device='%s', element index '%d' is absent from IfXTable", community, address, k)
		}

		t := MergeTables(ifTable[k], ifXTable[k])

		// add additional data to output
		t.Address = address
		t.Community = community
		t.ElapsedSeconds = time.Now().Unix()
		t.ElapsedSecondsSinceEpoch = time.Now().Unix()

		// populate output
		out = append(out, t)
	}

	// sort output
	sort.Slice(out, func(i, j int) bool { return out[i].Index < out[j].Index })

	return out, nil
}
