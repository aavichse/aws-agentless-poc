from typing import List
import pandas as pd
from io import StringIO
from reveal.vpclog_reader import FlowLogsS3Reader, DEFAULT_FIELDS
from reveal.vpclog_processor import filter_and_aggregate_flowlog


class MockFlowLogsS3Reader(FlowLogsS3Reader): 
    HEADER = "action dstaddr dstport end flow-direction interface-id tcp-flags start srcport srcaddr vpc-id"
    def __init__(self, test_data: List[str], **kwargs): 
    
        self.lines = [self.HEADER] + test_data
        super().__init__(location='', **kwargs)
        
        
    def _get_all_keys(self):
        yield ""

    def _read_streams(self):
        for key in self._get_all_keys():
            yield from self._read_file(key)

    def _read_file(self, key):
        body = "\n".join(self.lines)
        file = StringIO(body)
        df = pd.read_csv(file, sep="\s+", header=0, low_memory=True)
        df.rename(columns=lambda x: x.replace("-", "_"), inplace=True)
        df.drop(columns=df.columns.difference(DEFAULT_FIELDS), inplace=True)
        for index, row in df.iterrows():
            yield row.to_dict()
            
            
def test_ignore_NODATA():
    test_data = [
        "- - - 1717931268 - eni-0018d12067d83428f - 1717931236 - - vpc-0bf8ac1eaac82eb4f",
        "- - - 1717931306 - eni-007f7adc64ab3cf8d - 1717931276 - - vpc-0bf8ac1eaac82eb4f"
    ]

    mock_reader = MockFlowLogsS3Reader(boto_client=None, test_data=test_data)

    records = []
    for rec in mock_reader:
        records.append(rec)
        
    revealed_records = filter_and_aggregate_flowlog(records)
        
    assert(len(revealed_records) == 0)
    

def test_filter_drop_non_syn_ack():
    test_data = [
        "ACCEPT 10.156.53.94 443 1717931126 ingress eni-0d597f1ad6ba4d433 3 1717931111 46618 10.135.159.179 vpc-0bf8ac1eaac82eb4f",
        "ACCEPT 10.156.53.94 443 1717931126 ingress eni-0d597f1ad6ba4d433 3 1717931111 46628 10.135.159.179 vpc-0bf8ac1eaac82eb4f"
    ]

    mock_reader = MockFlowLogsS3Reader(boto_client=None, test_data=test_data)

    records = []
    for rec in mock_reader:
        records.append(rec)
        
    revealed_records = filter_and_aggregate_flowlog(records)
        
    assert(len(revealed_records) == 0)        


def test_not_filter_syn_ack():
    test_data = [
        "ACCEPT 10.156.53.90 19876 1717931102 ingress eni-03a9b8d35cb797b74 19 1717931077 9997 10.135.159.155 vpc-0bf8ac1eaac82eb4f",
        "ACCEPT 10.156.53.90 55422 1717931162 ingress eni-03a9b8d35cb797b74 23 1717931137 3128 10.39.66.83 vpc-0bf8ac1eaac82eb4f",
        "ACCEPT 10.156.53.90 12122 1717931134 ingress eni-03a9b8d35cb797b74 18 1717931105 3128 10.39.66.227 vpc-0bf8ac1eaac82eb4f"
    ]

    mock_reader = MockFlowLogsS3Reader(boto_client=None, test_data=test_data)

    records = []
    for rec in mock_reader:
        records.append(rec)
        
    revealed_records = filter_and_aggregate_flowlog(records)
        
    assert(len(revealed_records) == 3)    
    
    
def test_aggregated_records():
    test_data = [
        "ACCEPT 10.156.53.90 20876 1717931102 ingress eni-03a9b8d35cb797b74 19 1717931077 9997 10.135.159.155 vpc-0bf8ac1eaac82eb4f",
        "ACCEPT 10.156.53.90 29876 1717931102 ingress eni-03a9b8d35cb797b74 19 1717931077 9997 10.135.159.155 vpc-0bf8ac1eaac82eb4f",
        "ACCEPT 10.156.53.90 39876 1717931102 ingress eni-03a9b8d35cb797b74 19 1717931077 9997 10.135.159.155 vpc-0bf8ac1eaac82eb4f",
    ]

    mock_reader = MockFlowLogsS3Reader(boto_client=None, test_data=test_data)

    records = []
    for rec in mock_reader:
        records.append(rec)
        
    revealed_records = filter_and_aggregate_flowlog(records)
        
    assert(len(revealed_records) == 1)  
    assert(revealed_records[0].count == 3)
    
    test_data = test_data + [
        "ACCEPT 10.156.54.148 25366 1717931146 egress eni-0b3db7c96411f608f 19 1717931118 443 10.156.52.47 vpc-0bf8ac1eaac82eb4f",
        "ACCEPT 10.156.54.148 22366 1717931146 egress eni-0b3db7c96411f608f 19 1717931118 443 10.156.52.47 vpc-0bf8ac1eaac82eb4f",
    ]
    
    mock_reader = MockFlowLogsS3Reader(boto_client=None, test_data=test_data)
    records = []
    for rec in mock_reader:
        records.append(rec)
        
    revealed_records = filter_and_aggregate_flowlog(records)
    
    assert(len(revealed_records) == 2)  
    assert(revealed_records[1].count == 2)
    
    