import json
import boto3
import threading
from typing import List
from datetime import datetime

from pydantic import BaseModel, Field
from .aws_resource_reader import AWSEC2Reader
from aggregator.model.integrations.v1.common.inventory import InventoryItem

from common.logger import get_logger

LOG = get_logger(module_name=__name__)


class AwsInventory(BaseModel):
    version: int = 0
    items: List[InventoryItem] = Field([])

    def write_to_json_file(self, file_path: str):
        with open(file_path, 'w') as f:
            f.write(self.model_dump_json(indent=4))


class AWSInventoryFetcher:

    def __init__(self, version: int, regions: list[str]):
        self.regions = regions
        self.inventory = AwsInventory(version=version)

    def go(self) -> AwsInventory:
        """for the POC we will implement a single threaded fetcher 
        """
        for region in self.regions: 
            boto_session = boto3.Session(region_name=region)
            ec2_reader = AWSEC2Reader(session=boto_session)
            for item in ec2_reader: 
                self.inventory.items.append(item)

        return self.inventory


class AWSOrchestrator:

    def __init__(self, regions: list[str], interval: int = 120):
        self.regions = regions

        self.timer = None
        self.interval = interval
        self.is_running = False

        self.fetcher = None
        self.inventory = AwsInventory()

    def start(self):
        self.is_running = True
        self.start_time = datetime.now()
        self._schedule_next_run(0)

    def stop(self):
        self.is_running = False
        if self.timer:
            self.timer.cancel()

    def _schedule_next_run(self, delay: int = None):
        if self.is_running:
            if delay is None:
                current_time = datetime.now()
                elapsed = (current_time - self.start_time).total_seconds()
                delay = self.interval - (elapsed % self.interval)

            self.timer = threading.Timer(delay, self._run)
            self.timer.start()

    def _run(self):
        LOG.info(f"Attempting to fetch inventory at {datetime.now()}")

        if self.fetcher:
            LOG.warn("Previous fetching in progress")
        else:
            new_ver = self.inventory.version + 1
            self.fetcher = AWSInventoryFetcher(new_ver, self.regions)
            try:
                self.inventory = self.fetcher.go()
                self.inventory.write_to_json_file('inventory.json')
                LOG.info(f"Fetched inventory {new_ver=}")
            except Exception as e:
                LOG.error(f"Failed to fetch inventory {new_ver=}, {e=}")
            finally:
                self.fetcher = None

        self._schedule_next_run()
