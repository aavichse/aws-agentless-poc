import json
import boto3
import threading
from typing import List, Dict, Optional
from datetime import datetime

from pydantic import BaseModel, Field
from .resource_reader import AWSEC2Reader
from aggregator.model.integrations.v1.common.inventory import InventoryItem

from common.logger import get_logger

LOG = get_logger(module_name=__name__)


class AwsInventory(BaseModel):
    version: int = 0
    items: List[InventoryItem] = Field([])

    # lookups
    items_lookup: Optional[Dict[str, InventoryItem]] = Field({}, exclude=True)
    private_ips_lookup: Optional[Dict[str, InventoryItem]] = Field({}, exclude=True)

    def write_to_json_file(self, file_path: str):
        with open(file_path, 'w') as f:
            f.write(self.model_dump_json(indent=4))
            
    def get_items(self, tag, value) -> List[InventoryItem]: 
        items = [i for i in self.items if i.entity_data.os_details['tags'].get(tag, '')==value]
        return items


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
                self.inventory.items_lookup[item.item_id] = item
                # POC:  for single region, single VPC 
                ip = item.entity_data.nics[0].private_ip_addresses[0]
                self.inventory.private_ips_lookup[ip] = item

        return self.inventory


class AWSOrchestrator:

    def __init__(self, regions: list[str], interval: int = 120):
        self.regions = regions

        self.timer = None
        self.interval = interval
        self.is_running = False

        self.fetcher = None
        self._inventory = AwsInventory()

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

            self.timer = threading.Timer(delay, self.fetch)
            self.timer.start()

    def fetch(self, force: bool=False):
        LOG.info(f"Attempting to fetch inventory at {datetime.now()}")

        if self.fetcher:
            LOG.warn("Fetching in progress")
        else:
            new_ver = self._inventory.version + 1
            self.fetcher = AWSInventoryFetcher(new_ver, self.regions)
            try:
                self._inventory = self.fetcher.go()
                self._inventory.write_to_json_file('inventory.json')
                LOG.info(f"Fetched inventory {new_ver=}")
            except Exception as e:
                LOG.error(f"Failed to fetch inventory {new_ver=}, {e=}")
            finally:
                self.fetcher = None
        if not force:
            self._schedule_next_run()

    @property
    def inventory(self): 
        return self._inventory