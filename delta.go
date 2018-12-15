package main

import (
	"sort"
)

// DeltaIfMergedTable - computes delta values between current SNMP stats and previously acquired stats
func DeltaIfMergedTable(before, after IfMergedTable) IfMergedTable {

	// declare output variable
	var out IfMergedTable

	// declare anonymous delta helper function
	deltaCounter := func(x, y int64) int64 {

		// counter can only increase in value
		if y-x < 0 {
			return 0
		}

		// return non-zero delta
		return y - x
	}

	// declare anonymous divide helper function (numerator/denominator) then denominator >= 1
	divide := func(numerator, denominator int64) int64 {

		// ensure no zero division
		if denominator <= 0 {
			denominator = 1
		}

		return numerator / denominator
	}

	out.Address = after.Address
	out.Community = after.Community

	out.ElapsedSeconds = deltaCounter(before.ElapsedSeconds, after.ElapsedSeconds)
	out.ElapsedSecondsSinceEpoch = after.ElapsedSeconds

	out.Index = after.Index

	out.Name = after.Name
	out.Descr = after.Descr
	out.Alias = after.Alias
	out.Type = after.Type

	out.Mtu = after.Mtu

	out.Speed = after.Speed
	out.HighSpeed = after.HighSpeed

	out.PhysAddress = after.PhysAddress
	out.AdminStatus = after.AdminStatus
	out.OperStatus = after.OperStatus

	out.InOctets =
		divide(
			deltaCounter(before.InOctets, after.InOctets),
			out.ElapsedSeconds,
		)
	out.InUcastPkts =
		divide(
			deltaCounter(before.InUcastPkts, after.InUcastPkts),
			out.ElapsedSeconds,
		)
	out.InDiscards =
		divide(
			deltaCounter(before.InDiscards, after.InDiscards),
			out.ElapsedSeconds,
		)
	out.InErrors =
		divide(
			deltaCounter(before.InErrors, after.InErrors),
			out.ElapsedSeconds,
		)
	out.InUnknownProtos =
		divide(
			deltaCounter(before.InUnknownProtos, after.InUnknownProtos),
			out.ElapsedSeconds,
		)
	out.OutOctets =
		divide(
			deltaCounter(before.OutOctets, after.OutOctets),
			out.ElapsedSeconds,
		)
	out.OutUcastPkts =
		divide(
			deltaCounter(before.OutUcastPkts, after.OutUcastPkts),
			out.ElapsedSeconds,
		)
	out.OutDiscards =
		divide(
			deltaCounter(before.OutDiscards, after.OutDiscards),
			out.ElapsedSeconds,
		)
	out.OutErrors =
		divide(
			deltaCounter(before.OutErrors, after.OutErrors),
			out.ElapsedSeconds,
		)

	out.InMulticastPkts =
		divide(
			deltaCounter(before.InMulticastPkts, after.InMulticastPkts),
			out.ElapsedSeconds,
		)
	out.InBroadcastPkts =
		divide(
			deltaCounter(before.InBroadcastPkts, after.InBroadcastPkts),
			out.ElapsedSeconds,
		)
	out.OutMulticastPkts =
		divide(
			deltaCounter(before.OutMulticastPkts, after.OutMulticastPkts),
			out.ElapsedSeconds,
		)
	out.OutBroadcastPkts =
		divide(
			deltaCounter(before.OutBroadcastPkts, after.OutBroadcastPkts),
			out.ElapsedSeconds,
		)
	out.HCInOctets =
		divide(
			deltaCounter(before.HCInOctets, after.HCInOctets),
			out.ElapsedSeconds,
		)
	out.HCInUcastPkts =
		divide(
			deltaCounter(before.HCInUcastPkts, after.HCInUcastPkts),
			out.ElapsedSeconds,
		)
	out.HCInMulticastPkts =
		divide(
			deltaCounter(before.HCInMulticastPkts, after.HCInMulticastPkts),
			out.ElapsedSeconds,
		)
	out.HCInBroadcastPkts =
		divide(
			deltaCounter(before.HCInBroadcastPkts, after.HCInBroadcastPkts),
			out.ElapsedSeconds,
		)
	out.HCOutOctets =
		divide(
			deltaCounter(before.HCOutOctets, after.HCOutOctets),
			out.ElapsedSeconds,
		)
	out.HCOutUcastPkts =
		divide(
			deltaCounter(before.HCOutUcastPkts, after.HCOutUcastPkts),
			out.ElapsedSeconds,
		)
	out.HCOutMulticastPkts =
		divide(
			deltaCounter(before.HCOutMulticastPkts, after.HCOutMulticastPkts),
			out.ElapsedSeconds,
		)
	out.HCOutBroadcastPkts =
		divide(
			deltaCounter(before.HCOutBroadcastPkts, after.HCOutBroadcastPkts),
			out.ElapsedSeconds,
		)

	out.PromiscuousMode = after.PromiscuousMode
	out.ConnectorPresent = after.ConnectorPresent

	return out
}

// DeltaIfMergedTables - wrapper for deltaIfMergedTable, multiple values
func DeltaIfMergedTables(before, after []IfMergedTable) []IfMergedTable {

	// declare map for input data, for easier lookups
	beforeMap := make(map[string]IfMergedTable)
	afterMap := make(map[string]IfMergedTable)

	// populate before map
	for i := range before {
		beforeMap[before[i].Name] = before[i]
	}

	// populate after map
	for i := range after {
		afterMap[after[i].Name] = after[i]
	}

	// declare output variable
	out := make([]IfMergedTable, 0)

	// do delte only then befare and after values are defined for specified key
	for k := range afterMap {
		if _, ok := beforeMap[k]; !ok {
			continue
		}

		// populate output
		out = append(out, DeltaIfMergedTable(beforeMap[k], afterMap[k]))
	}

	// sort output
	sort.Slice(out, func(i, j int) bool { return out[i].Index < out[j].Index })

	return out
}
