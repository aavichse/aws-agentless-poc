from dataclasses import dataclass
from typing import Dict, List, Set, Tuple
import requests
from retry import retry
from common.logger import get_logger
from aggregator.model.integrations.v1.common.inventory import (
    InventoryItem,
    VMData,
    ManagedServiceData,
)

LOG = get_logger(module_name=__name__)


INVENTORY_URL = "http://localhost/v1/provider/inventory"
# INVENTORY_URL = "http://ec2-44-204-113-186.compute-1.amazonaws.com/v1/provider/inventory"

@dataclass 
class IPAddressInfo: 
    ip_addresses: Set[str]
    item: InventoryItem
    
#@retry(tries=5, delay=2, backoff=2, exceptions=(requests.exceptions.RequestException,))
def fetch_inventory_with_retries():
    """
    Request discovered resources from the inventory service
    """
    try:
        response = requests.get(INVENTORY_URL)
        response.raise_for_status()  # Raise an HTTPError for bad responses
        data = response.json()
        return data
    except requests.exceptions.HTTPError as http_err:
        LOG.error(f"HTTP error occurred: {http_err}")
        raise
    except Exception as err:
        LOG.error(f"An error occurred: {err}")
        raise


def fetch_inventory() -> List[InventoryItem]:
    try:
        result = fetch_inventory_with_retries() or []
        inventory = [InventoryItem(**item) for item in result]
        LOG.info(f"total items={len(inventory)}")
        return inventory
    except Exception as e:
        LOG.error(f"Failed to fetch inventory after retries: {e}")


def get_connection_item_info(dicovered_item: InventoryItem) -> InventoryItem:
    return InventoryItem.model_validate(
        {
            'item-id': dicovered_item.item_id,
            'item-type': dicovered_item.entity_type,
            'external-ids': dicovered_item.external_ids
        })
  
    
# def populate_eni(inventory_items: List[InventoryItem]) -> Dict[Tuple[str, str], Tuple[Set[str], InventoryItem]]:
#     ipmap = {}   # key=(eni ID, vpc ID), value=(list of ip addresses, inventory item)
#     for item in inventory_items:
#         entity_data = item.entity_data
#         if isinstance(entity_data, (VMData, ManagedServiceData)) and entity_data.nics:
#             for nic in entity_data.nics:
#                 key = (nic.id, nic.network)
#                 connection_item_info = get_connection_item_info(item)
#                 ips = set(nic.private_ip_addresses) + set(nic.public_ip_addresses)
                
#     return ipmap


def popuplate_ipmap(inventory_items: List[InventoryItem]) -> Dict[str, InventoryItem]:
    ipmap = {}
    for item in inventory_items:
        entity_data = item.entity_data
        if isinstance(entity_data, (VMData, ManagedServiceData)) and entity_data.nics:
            for nic in entity_data.nics:
                if nic.private_ip_addresses:
                    for ip in nic.private_ip_addresses:
                        ipmap[(ip, nic.network)] = get_connection_item_info(item)
                if nic.public_ip_addresses:
                    for ip in nic.public_ip_addresses:
                        ipmap[(ip, nic.network)] = get_connection_item_info(item)
    return ipmap


def fetch_ipmap() -> Dict[str, InventoryItem]:
    inventory = fetch_inventory() or []
    ipmap = popuplate_ipmap(inventory_items=inventory)
    return ipmap
