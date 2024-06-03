from reveal.vpclog_reader import FlowRecord
from .test_vpclog_reader import MockFlowLogsS3Reader

def read_record(line: str) -> FlowRecord:
    test_data = [line]
    mock_reader = MockFlowLogsS3Reader(boto_client=None, test_data=test_data)
    for record in mock_reader:
        return record

      
def test_read_single_record():
    line = "- - - 1717931268 - eni-0018d12067d83428f - 1717931236 - - vpc-0bf8ac1eaac82eb4f"
    record = read_record(line)
    print(record)
    
    
def test_read_single_record():
    line = "- - - 1717931268 - eni-0018d12067d83428f - 1717931236 - - vpc-0bf8ac1eaac82eb4f"
    record = read_record(line)
    print(record)
    



