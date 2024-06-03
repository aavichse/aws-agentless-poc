import threading
import time
from typing import Dict, List
import boto3
from datetime import datetime, timedelta, timezone
from .vpclog_reader import FlowLogsS3Reader, FlowRecord, FlowRecordHeader
from .vpclog_reveal import Reveal
from .vpclog_processor import filter_and_aggregate_flowlog
from aggregator.model.integrations.v1.common.inventory import InventoryItem
from .ipmap_fetcher import fetch_ipmap
from common.logger import get_logger


LOG = get_logger(module_name=__name__)

class VPCLogWorker:

    def __init__(self, vpclog_s3_bucket_name: str, interval: int, 
                 reporting_entity_id: str, dump: bool = False) -> None:
        self.vpclog_s3_bucket_name = vpclog_s3_bucket_name
        self.interval = interval
        self.reader_lock = threading.Lock()
        self.running = False
        self.worker_thread = None
        self.boto_client = boto3.client("s3")
        self.reporting_entity_id = reporting_entity_id
        self.log_records = dump

    def start(self): 
        if not self.running: 
            LOG.info(f'Starting vpclog worker at S3 bucket={self.vpclog_s3_bucket_name}')
            self.start_time = datetime.now() - timedelta(hours=1)
            self.running = True
            self.worker_thread = threading.Thread(target=self.run)
            self.worker_thread.daemon = True
            self.worker_thread.start()
            self.worker_thread.join()

    def run(self):
        next_run_time = time.time() + self.interval
        while self.running:
            if not self.reader_lock.locked():
                start_time = time.time()
                self.reader_lock.acquire()
                try:
                    self.read_flowlogs()
                finally:
                    self.reader_lock.release()
                duration = time.time() - start_time
                sleep_time = max(0, next_run_time - time.time())
            else:
                LOG.warn("Reader is already working, skipping this cycle.")
                sleep_time = max(0, next_run_time - time.time())

            time.sleep(sleep_time)
            next_run_time += self.interval  # Schedule the next run

    def stop(self): 
        self.running = False
        if self.worker_thread: 
            self.worker_thread.join()
            LOG.info(f'Stop vpclog worker at S3 bucket ={self.vpclog_s3_bucket_name}')

    def read_flowlogs(self):
        ipmap = fetch_ipmap()
        reader = FlowLogsS3Reader(location=self.vpclog_s3_bucket_name, boto_client=self.boto_client, dump=self.log_records)
        reader.memorize_previous_runs()

        buffered_records = []

        for rec in reader: 
            buffered_records.append(rec)
            if len(buffered_records) > 200:   # FIXME: 200 , work in batch 
                self.publish(records=buffered_records, ipmap=ipmap)
                buffered_records = []   # next batch

        if len(buffered_records) > 0:   # publish left over
            self.publish(records=buffered_records, ipmap=ipmap)
         
    def publish(self, records:List[FlowRecord], ipmap: Dict[str, InventoryItem]): 
        revealed_records = filter_and_aggregate_flowlog(records)
        reveal = Reveal(ipmap=ipmap, reporting_entity_id=self.reporting_entity_id)
        for record in revealed_records: 
            reveal.send(record)