from typing import Dict, List
from collections import defaultdict
from datetime import datetime
from .vpclog_reader import FlowRecord
from common.logger import get_logger

LOG = get_logger(module_name=__name__)

KEY_FIELDS = ('srcaddr', 'dstaddr', 'srcport')

class _FlowStats:
    """
    An aggregator for flow records. Sums bytes and packets and keeps track of
    the active time window.
    """
    reasonable_max_datetime = datetime(3000, 12, 31, 23, 59, 59)
    reasonable_min_datetime = datetime(1970, 1, 1, 0, 0, 0)

    def __init__(self):
        self.start = int(self.reasonable_max_datetime.timestamp())
        self.end = int(self.reasonable_min_datetime.timestamp())
        self.count = 0
        self.rec = None

    def update(self, flow_record):
        if self.rec is None: 
            event_data = flow_record.to_dict()
            self.rec = FlowRecord(event_data)
            self.rec.count = 0
            
        if self.rec.start < self.start:
            self.rec.start = flow_record.start
        
        if self.rec.end > self.end:
            self.rec.end = flow_record.end
        
        self.rec.count += 1

    def to_dict(self):
        return {x: getattr(self, x) for x in self.__slots__}


SYN_ACK = 18
SIN = 2
FIN = 1


def is_ephermal_port(record: FlowRecord) -> bool: 
    return record.dstport < record.srcport


def is_filtered(record: FlowRecord) -> bool:
    if record.srcaddr is None or record.dstaddr is None:
        return True
    
    if record.tcp_flags is None: 
        return True
    
    if record.tcp_flags >= SYN_ACK:
        return False
    
    if record.tcp_flags >= 1 and record.tcp_flags <= 3:
        return False
   
    # if record.tcp_flags == 2 or record.tcp_flags == 3:
    #     return False

    return True    



def filter_and_aggregate_flowlog(records: List[FlowRecord], key_fields=KEY_FIELDS) -> List[FlowRecord]:
    filtered = 0
    
    flow_table = defaultdict(_FlowStats)
    
    for rec in records: 
        
        if is_filtered(rec):
            filtered += 1
            continue
        
        key = tuple(getattr(rec, attr) for attr in key_fields)
        if any(x is None for x in key):
            continue

        flow_table[key].update(rec)
    
    aggregated_records =  [x.rec for x in flow_table.values()]
    squashed = len(records) - filtered - len(aggregated_records)
    LOG.info(f'origin={len(records)}, {filtered=}, aggregated={len(aggregated_records)}, {squashed=}')
    
    return aggregated_records
    
    