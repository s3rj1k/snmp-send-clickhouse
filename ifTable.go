package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const snmpTableValueSeparator = "-=-"

// http://www.net-snmp.org/docs/mibs/interfaces.html

// IfTable partially describes SNMP IfTable
type IfTable struct {
	Index       int
	Descr       string
	Type        string
	Mtu         int64
	Speed       int64 // GAUGE (bps)
	PhysAddress string
	AdminStatus string
	OperStatus  string
	// LastChange
	InOctets    int64 // COUNTER
	InUcastPkts int64 // COUNTER
	// InNUcastPkts
	InDiscards      int64 // COUNTER
	InErrors        int64 // COUNTER
	InUnknownProtos int64 // COUNTER
	OutOctets       int64 // COUNTER
	OutUcastPkts    int64 // COUNTER
	// OutNUcastPkts
	OutDiscards int64 // COUNTER
	OutErrors   int64 // COUNTER
	// OutQLen
	// Specific
}

// CmdIfTable - gets SNMP IfTable data for specified community and host
func CmdIfTable(community, address string) (map[int]IfTable, error) {

	// declare output variable
	out := make(map[int]IfTable)

	// run snmptable
	cmdOut := RunCommand("snmptable", "-CB", "-Cf", snmpTableValueSeparator, "-Cl", "-CH", "-v2c", fmt.Sprintf("-c%s", community), address, "ifTable")
	if cmdOut.ReturnCode != 0 {
		return nil, fmt.Errorf("failed to run '%s' with return code=%d: %s", cmdOut.Command, cmdOut.ReturnCode, cmdOut.CombinedOutput)
	}

	// prepare scanner interface for snmptable output
	scanner := bufio.NewScanner(bytes.NewReader(cmdOut.CombinedOutput))

	// scan output lines
	for scanner.Scan() {

		var t IfTable
		var err error

		// split by custom field separator
		fields := strings.Split(scanner.Text(), snmpTableValueSeparator)

		// validate number of fields
		if len(fields) != 22 {
			return nil, fmt.Errorf("wrong number of parameters, must be 22 but got %d", len(fields))
		}

		if t.Index, err = strconv.Atoi(fields[0]); err != nil {
			return nil, fmt.Errorf("failed to parse IfIndex (IfTable): %s", fields[0])
		}

		t.Descr = fields[1]
		t.Type = fields[2]

		if t.Mtu, err = strconv.ParseInt(fields[3], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfMtu (IfTable): %s", fields[3])
		}

		if t.Speed, err = strconv.ParseInt(fields[4], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfSpeed (IfTable): %s", fields[4])
		}

		t.PhysAddress = fields[5]
		t.AdminStatus = fields[6]
		t.OperStatus = fields[7]

		if t.InOctets, err = strconv.ParseInt(fields[9], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInOctets (IfTable): %s", fields[9])
		}

		if t.InUcastPkts, err = strconv.ParseInt(fields[10], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInUcastPkts (IfTable): %s", fields[10])
		}

		if t.InDiscards, err = strconv.ParseInt(fields[12], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInDiscards (IfTable): %s", fields[12])
		}

		if t.InErrors, err = strconv.ParseInt(fields[13], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInErrors (IfTable): %s", fields[13])
		}

		if t.InUnknownProtos, err = strconv.ParseInt(fields[14], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfInUnknownProtos (IfTable): %s", fields[14])
		}

		if t.OutOctets, err = strconv.ParseInt(fields[15], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfOutOctets (IfTable): %s", fields[15])
		}

		if t.OutUcastPkts, err = strconv.ParseInt(fields[16], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfOutUcastPkts (IfTable): %s", fields[16])
		}

		if t.OutDiscards, err = strconv.ParseInt(fields[18], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfOutDiscards (IfTable): %s", fields[18])
		}

		if t.OutErrors, err = strconv.ParseInt(fields[19], 10, 64); err != nil {
			return nil, fmt.Errorf("failed to parse IfOutErrors (IfTable): %s", fields[19])
		}

		// populate output
		out[t.Index] = t
	}

	// fail no scanner error
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("faild to read command output: %s", err.Error())
	}

	return out, nil
}
