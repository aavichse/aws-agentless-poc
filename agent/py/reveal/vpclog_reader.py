import os
from os.path import basename
from calendar import timegm
from datetime import datetime, timedelta, timezone
from typing import List
from dateutil.rrule import rrule, DAILY
import pandas as pd
from io import BytesIO
import boto3
from common.logger import get_logger

LOG = get_logger(module_name=__name__)


DEFAULT_FIELDS = (
    "version",
    "account_id",
    "interface_id",
    "srcaddr",
    "dstaddr",
    "srcport",
    "dstport",
    "protocol",
    "packets",
    "bytes",
    "start",
    "end",
    "action",
    "log_status",
    "vpc_id",
    "subnet_id",
    "instance_id",
    "tcp_flags",
    "type",
    "pkt_srcaddr",
    "pkt_dstaddr",
    "region",
    "az_id",
    "sublocation_type",
    "sublocation_id",
    "pkt_src_aws_service",
    "pkt_dst_aws_service",
    "flow_direction",
    "traffic_path",
    "resource_type",
    "tgw_id",
    "tgw_attachment_id",
    "tgw_src_vpc_account_id",
    "tgw_dst_vpc_account_id",
    "tgw_src_vpc_id",
    "tgw_dst_vpc_id",
    "tgw_src_subnet_id",
    "tgw_dst_subnet_id",
    "tgw_src_eni",
    "tgw_dst_eni",
    "tgw_src_az_id",
    "tgw_dst_az_id",
    "tgw_pair_attachment_id",
    "packets_lost_no_route",
    "packets_lost_blackhole",
    "packets_lost_mtu_exceeded",
    "packets_lost_ttl_expired",
)

# S3_CSV_FILTER_QUERY = "log_status == 'OK' and " \
#                       "tcp_flags in ('2', '3', '4', '18', '19') and " \
#                       "protocol in ('6', '17') and " \
#                       "srcport != '0' and " \
#                       "dstport != '0'"
S3_CSV_FILTER_QUERY = "dstport == '9000'"


class FlowRecordHeader: 
    def __init__(self, fields: List[str]):
        self.header = fields

class FlowRecord:
    """
    Given a VPC Flow Logs event dictionary, returns a Python object whose
    attributes match the field names in the event record. Integers are stored
    as Python int objects; timestamps are stored as Python datetime objects.
    """

    __slots__ = [
        "version",
        "account_id",
        "interface_id",
        "action",
        "flow_direction",
        "srcaddr",
        "srcport",
        "dstaddr",
        "dstport",
        "protocol",
        "start",
        "end",
        "packets",
        "bytes",
        "log_status",
        "instance_id",
        "vpc_id",
        "subnet_id",
        "tcp_flags",
        "type",
        "pkt_srcaddr",
        "pkt_dstaddr",
        "region",
        "az_id",
        "sublocation_type",
        "sublocation_id",
        "pkt_src_aws_service",
        "pkt_dst_aws_service",
        "traffic_path",
        "resource_type",
        "tgw_id",
        "tgw_attachment_id",
        "tgw_src_vpc_account_id",
        "tgw_dst_vpc_account_id",
        "tgw_src_vpc_id",
        "tgw_dst_vpc_id",
        "tgw_src_subnet_id",
        "tgw_dst_subnet_id",
        "tgw_src_eni",
        "tgw_dst_eni",
        "tgw_src_az_id",
        "tgw_dst_az_id",
        "tgw_pair_attachment_id",
        "packets_lost_no_route",
        "packets_lost_blackhole",
        "packets_lost_mtu_exceeded",
        "packets_lost_ttl_expired",
        "count",
    ]

    def __init__(self, event_data, EPOCH_32_MAX=2147483647):
        # Contra the docs, the start and end fields can contain
        # millisecond-based timestamps.
        # http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/flow-logs.html
        if "start" in event_data:
            start = int(event_data["start"])
            if start > EPOCH_32_MAX:
                start /= 1000
            self.start = start  # datetime.utcfromtimestamp(start)
        else:
            self.start = None

        if "end" in event_data:
            end = int(event_data["end"])
            if end > EPOCH_32_MAX:
                end /= 1000
            self.end = end  # datetime.utcfromtimestamp(end)
        else:
            self.end = start + 1000

        for key, func in (
            ("version", int),
            ("account_id", str),
            ("interface_id", str),
            ("srcaddr", str),
            ("dstaddr", str),
            ("srcport", int),
            ("dstport", int),
            ("protocol", int),
            ("packets", int),
            ("bytes", int),
            ("action", str),
            ("log_status", str),
            ("vpc_id", str),
            ("subnet_id", str),
            ("instance_id", str),
            ("tcp_flags", int),
            ("type", str),
            ("pkt_srcaddr", str),
            ("pkt_dstaddr", str),
            ("region", str),
            ("az_id", str),
            ("sublocation_type", str),
            ("sublocation_id", str),
            ("pkt_src_aws_service", str),
            ("pkt_dst_aws_service", str),
            ("flow_direction", str),
            ("traffic_path", int),
            ("resource_type", str),
            ("tgw_id", str),
            ("tgw_attachment_id", str),
            ("tgw_src_vpc_account_id", str),
            ("tgw_dst_vpc_account_id", str),
            ("tgw_src_vpc_id", str),
            ("tgw_dst_vpc_id", str),
            ("tgw_src_subnet_id", str),
            ("tgw_dst_subnet_id", str),
            ("tgw_src_eni", str),
            ("tgw_dst_eni", str),
            ("tgw_src_az_id", str),
            ("tgw_dst_az_id", str),
            ("tgw_pair_attachment_id", str),
            ("packets_lost_no_route", int),
            ("packets_lost_blackhole", int),
            ("packets_lost_mtu_exceeded", int),
            ("packets_lost_ttl_expired", int),
            ("count", int),
            
        ):
            value = event_data.get(key, "-")
            if value == "-" or value == "None" or value is None:
                value = None
            else:
                value = func(value)

            setattr(self, key, value)

    def __eq__(self, other):
        try:
            return all(getattr(self, x) == getattr(other, x) for x in self.__slots__)
        except AttributeError:
            return False

    # FIXME -
    def __hash__(self):
        return hash(tuple(getattr(self, x) for x in self.__slots__))

    def __str__(self):
        ret = []
        for key in self.__slots__:
            value = getattr(self, key)
            if value is not None:
                ret.append("{}: {}".format(key, value))
        return ", ".join(ret)

    def to_dict(self):
        ret = {}
        for key in self.__slots__:
            value = getattr(self, key)
            if value is not None:
                ret[key] = value

        return ret

    def to_message(self):
        D_transform = {
            # 'start': lambda dt: str(timegm(dt.utctimetuple())),
            # 'end': lambda dt: str(timegm(dt.utctimetuple())),
        }

        ret = []
        for attr in self.__slots__:
            transform = D_transform.get(attr, lambda x: str(x) if x else "-")
            ret.append(transform(getattr(self, attr)))

        return " ".join(ret)
    
    def to_short_message(self):
        return f'{self.action}  {self.srcaddr}:{self.srcport}-{self.dstaddr}:{self.dstport};{self.count}'


class BaseReader:
    def __init__(
        self,
        boto_client: boto3.client,
        start_time: datetime = None,
        end_time: datetime = None,
        dump: bool = False, 
    ):
        self.boto_client = boto_client
        self.dump = dump

        now = datetime.now()
        self.start_time = start_time or now - timedelta(hours=8)
        self.end_time = end_time or now + timedelta(hours=1)

        LOG.info(f"Using time range: {self.start_time} - {self.end_time}")

        self.iterator = self._reader()

    def _reader(self):
        raise NotImplementedError()

    def __iter__(self):
        return self

    def __next__(self):
        return next(self.iterator)


class FlowLogsS3Reader(BaseReader):
    def __init__(
        self,
        location: str,
        **kwargs,
    ):
        super().__init__(**kwargs)
        location_parts = (location.rstrip("/") + "/").split("/", 1)
        self.bucket, self.prefix = location_parts
        self.done_keys = set()
        self.done_filename = self.bucket+".done.txt"
        self.current_key = None

    def _get_account_prefixes(self):
        """
        Yield each prefix of the type:
        <bucket-name>/AWSLogs/account-id=<acount>/

        Yields:
            _type_: string
        """
        prefix = self.prefix.strip("/") + "/AWSLogs/"
        prefix = prefix.lstrip("/")
        resp = self.boto_client.list_objects_v2(
            Bucket=self.bucket, Delimiter="/", Prefix=prefix
        )
        for item in resp.get("CommonPrefixes", []):
            prefix = item["Prefix"]
            #account_id = prefix.rsplit("=", 2)[1][:-1]
            account_id = '/'+prefix.split('/')[1]
            # LOG.info(f"Found account: {account_id}, {prefix=}")
            yield prefix

    def _get_region_prefixes(self, account_prefix):
        """Yield each prefix of the type:
        base_location/AWSLogs/account_number/vpcflowlogs/region_name/

        Args:
            account_prefix (_type_): _description_

        Yields:
            _type_: string
        """
        resp = self.boto_client.list_objects_v2(
            Bucket=self.bucket,
            Delimiter="/",
            Prefix=account_prefix + "vpcflowlogs/",
        )
        for item in resp.get("CommonPrefixes", []):
            prefix = item["Prefix"]
            #region_name = prefix.rsplit("=", 2)[2][:-1]
            region_name = '/'+prefix.split('/')[-2]
            # LOG.info(f"Found region: {region_name}, {prefix=}")
            yield prefix

    def _get_date_prefixes(self):
        """Each base_location/AWSLogs/account_number/vpcflowlogs/region_name/
        prefix has files organized in year/month/day directories.
        Yield the year/month/day/ fragments that are relevant to our time range

        Yields:
            _type_: _description_
        """
        dtstart = self.start_time.replace(hour=0, minute=0, second=0, microsecond=0)
        until = self.end_time.replace(hour=0, minute=0, second=0, microsecond=0)
        for dt in rrule(freq=DAILY, dtstart=dtstart, until=until):
            date_prefix = dt.strftime("%Y/%m/%d/")
            date = date_prefix[:-1]
            # LOG.info(f"Found date: {date_prefix}")
            yield date_prefix

    def _get_keys(self, prefix):
        """S3 keys have a file name like:
        <account>_vpcflowlogs_<region>_<flow-logs-id>_<datetime_hash>.log.gz
        Yield the keys for files relevant to our time range

        Args:
            prefix (_type_): _description_

        Yields:
            _type_: _description_
        """
        paginator = self.boto_client.get_paginator("list_objects_v2")
        all_pages = paginator.paginate(Bucket=self.bucket, Prefix=prefix)
        for page in all_pages:
            for item in page.get("Contents", []):
                key = item["Key"]
                file_name = basename(key)
                try:
                    dt = datetime.strptime(file_name.rsplit("_", 2)[1], "%Y%m%dT%H%MZ")
                except (IndexError, ValueError):
                    continue

                if self.start_time <= dt < self.end_time:
                    if key in self.done_keys:
                        #LOG.info(f"Skipping key: {key}")
                        continue

                    LOG.info(f"Found key: {key}")
                    yield key

    def _get_all_keys(self):
        for account_prefix in self._get_account_prefixes():
            for region_prefix in self._get_region_prefixes(account_prefix):
                for day_prefix in self._get_date_prefixes():
                    prefix = region_prefix + day_prefix
                    for key in self._get_keys(prefix):
                        yield key

    def dump_csv(self, csvname: str, values: List[str]):
        with open(csvname, 'a') as csv:
            csv.write(" ".join(values) + "\n")

    def _read_file(self, key):
        self.current_key = key
        
        csv_file = None
        if self.dump: 
            csv_file = key.split('/')[-1]+".csv"
        
        resp = self.boto_client.get_object(Bucket=self.bucket, Key=key)
        body = resp["Body"].read()

        if key.endswith(".parquet"):
            raise NotImplementedError("Parquet files are not yet supported")
        else:
            file = BytesIO(body)
            df = pd.read_csv(
                file, compression="gzip", sep="\s+", header=0, low_memory=True
            )
            
            if csv_file:
                self.dump_csv(csv_file, df.columns.to_list())

            df.rename(columns=lambda x: x.replace("-", "_"), inplace=True)
            df.drop(columns=df.columns.difference(DEFAULT_FIELDS), inplace=True)
            for index, row in df.iterrows():
                if csv_file:
                    self.dump_csv(csv_file, map(str, row.values))
                    
                yield row.to_dict()
        self.mark_done()

    def _read_streams(self):
        for key in self._get_all_keys():
            yield from self._read_file(key)

    def _reader(self):
        for event_data in self._read_streams():
            try:
                if isinstance(event_data, list):
                    yield FlowRecordHeader(event_data)
                else:
                    yield FlowRecord(event_data)
            except Exception as e:
                LOG.error(f"Error reading record: {e}")

    def memorize_previous_runs(self):
        if not os.path.exists(self.done_filename):
            return

        with open(self.done_filename, "r") as file:
            for line in file:
                # Strip whitespace from the beginning and end of the line
                # This helps in removing newline characters and spaces
                cleaned_line = line.strip()

                # Add the cleaned line to the set
                self.done_keys.add(cleaned_line)

    def mark_done(self):
        with open(self.done_filename, "a") as file:
            file.write(self.current_key + "\n")