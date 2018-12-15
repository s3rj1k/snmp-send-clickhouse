CREATE DATABASE IF NOT EXISTS snmp;

CREATE TABLE IF NOT EXISTS snmp.table (
  Address                  String,
  Community                String,
  ElapsedSeconds           Int64,
  ElapsedSecondsSinceEpoch Int64,
  Index                    Int,
  Name                     String,
  Descr                    String,
  Alias                    String,
  Type                     String,
  Mtu                      Int64,
  Speed                    Int64,
  HighSpeed                Int64,
  PhysAddress              String,
  AdminStatus              String,
  OperStatus               String,
  InOctets                 Int64,
  InUcastPkts              Int64,
  InDiscards               Int64,
  InErrors                 Int64,
  InUnknownProtos          Int64,
  OutOctets                Int64,
  OutUcastPkts             Int64,
  OutDiscards              Int64,
  OutErrors                Int64,
  InMulticastPkts          Int64,
  InBroadcastPkts          Int64,
  OutMulticastPkts         Int64,
  OutBroadcastPkts         Int64,
  HCInOctets               Int64,
  HCInUcastPkts            Int64,
  HCInMulticastPkts        Int64,
  HCInBroadcastPkts        Int64,
  HCOutOctets              Int64,
  HCOutUcastPkts           Int64,
  HCOutMulticastPkts       Int64,
  HCOutBroadcastPkts       Int64,
  PromiscuousMode          String,
  ConnectorPresent         String,
  Date                     Date     MATERIALIZED toDate(ElapsedSecondsSinceEpoch),
  DateTime                 DateTime MATERIALIZED toDateTime(ElapsedSecondsSinceEpoch)
) ENGINE = MergeTree PARTITION BY (
  toDayOfMonth(toDate(ElapsedSecondsSinceEpoch)),
  Address
)
ORDER BY
  (Address, Community, Date);

CREATE TABLE IF NOT EXISTS snmp.table_buffer AS snmp.table ENGINE = Buffer(snmp, table, 16, 10, 100, 1000, 1000000, 1024, 1024000);
