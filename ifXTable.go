package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// http://net-snmp.sourceforge.net/docs/mibs/ifMIBObjects.html

// IfXTable partially describes SNMP IfXTable
type IfXTable struct {
	Index              int // snmptable option: (-Ci) This option prepends the index of the entry to all printed lines.
	Name               string
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
	// LinkUpDownTrapEnable
	HighSpeed        int64 // GAUGE (Mbps)
	PromiscuousMode  string
	ConnectorPresent string
	Alias            string
	// CounterDiscontinuityTime
}

// CmdIfXTable - gets SNMP IfXTable data for specified community and host
func CmdIfXTable(community, address string) (map[int]IfXTable, error) {

	// declare output variable
	out := make(map[int]IfXTable)

	// run snmptable
	cmdOut := RunCommand("snmptable", "-Ci", "-CB", "-Cf", snmpTableValueSeparator, "-Cl", "-CH", "-v2c", fmt.Sprintf("-c%s", community), address, "ifXTable")
	if cmdOut.ReturnCode != 0 {
		return nil, fmt.Errorf("failed to run '%s' with return code=%d: %s", cmdOut.Command, cmdOut.ReturnCode, cmdOut.CombinedOutput)
	}

	// prepare scanner interface for snmptable output
	scanner := bufio.NewScanner(bytes.NewReader(cmdOut.CombinedOutput))

	// scan output lines
	for scanner.Scan() {

		var t IfXTable
		var err error

		// split by custom field separator
		fields := strings.Split(scanner.Text(), snmpTableValueSeparator)

		// validate number of fields
		if len(fields) != 20 {
			return nil, fmt.Errorf("wrong number of parameters, must be 19 but got %d", len(fields))
		}

		if t.Index, err = strconv.Atoi(fields[0]); err != nil {
			return nil, fmt.Errorf("failed to parse index of entry (IfXTable): %s", fields[0])
		}

		t.Name = fields[1]

		if t.InMulticastPkts, err = strconv.ParseInt(fields[2], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInMulticastPkts (IfXTable): %s", fields[2])
		}

		if t.InBroadcastPkts, err = strconv.ParseInt(fields[3], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInBroadcastPkts (IfXTable): %s", fields[3])
		}

		if t.OutMulticastPkts, err = strconv.ParseInt(fields[4], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfOutMulticastPkts (IfXTable): %s", fields[4])
		}

		if t.OutBroadcastPkts, err = strconv.ParseInt(fields[5], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfOutBroadcastPkts (IfXTable): %s", fields[5])
		}

		if t.HCInOctets, err = strconv.ParseInt(fields[6], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCInOctets (IfXTable): %s", fields[6])
		}

		if t.HCInUcastPkts, err = strconv.ParseInt(fields[7], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCInUcastPkts (IfXTable): %s", fields[7])
		}

		if t.HCInMulticastPkts, err = strconv.ParseInt(fields[8], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCInMulticastPkts (IfXTable): %s", fields[8])
		}

		if t.HCInBroadcastPkts, err = strconv.ParseInt(fields[9], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCInBroadcastPkts (IfXTable): %s", fields[9])
		}

		if t.HCOutOctets, err = strconv.ParseInt(fields[10], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCOutOctets (IfXTable): %s", fields[10])
		}

		if t.HCOutUcastPkts, err = strconv.ParseInt(fields[11], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCOutUcastPkts (IfXTable): %s", fields[11])
		}

		if t.HCOutMulticastPkts, err = strconv.ParseInt(fields[12], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCOutMulticastPkts (IfXTable): %s", fields[12])
		}

		if t.HCOutBroadcastPkts, err = strconv.ParseInt(fields[13], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHCOutBroadcastPkts (IfXTable): %s", fields[13])
		}

		if t.HighSpeed, err = strconv.ParseInt(fields[15], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfHighSpeed (IfXTable): %s", fields[15])
		}

		t.PromiscuousMode = fields[16]
		t.ConnectorPresent = fields[17]
		t.Alias = fields[18]

		// populate output
		out[t.Index] = t
	}

	// fail no scanner error
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("faild to read command output: %s", err.Error())
	}

	return out, nil
}
