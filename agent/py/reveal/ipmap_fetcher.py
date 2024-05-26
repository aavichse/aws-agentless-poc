from typing import Dict, List
import requests
from retry import retry
from common.logger import get_logger
from aggregator.model.integrations.v1.common.inventory import (
    InventoryItem,
    VMData,
    ManagedServiceData,
)

LOG = get_logger(module_name=__name__)


# Nginx forward the request to inventory service
INVENTORY_URL = "http://localhost:8080/v1/provider/inventory"


@retry(tries=5, delay=2, backoff=2, exceptions=(requests.exceptions.RequestException,))
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
    result = []
    try:
        result = fetch_inventory_with_retries()
        inventory = [InventoryItem(**item) for item in result]
        LOG.info(f"total items={len(inventory)}")
        return inventory
    except Exception as e:
        LOG.error(f"Failed to fetch inventory after retries: {e}")


def get_inventory_item_info(id: str, type: str, external_ids: List[str]) -> InventoryItem:
    return InventoryItem.model_validate(
        {
            'item-id': id,
            'item-type': type,
            'external-ids': external_ids
        })

def popuplate_ipmap(inventory_items: List[InventoryItem]) -> Dict[str, InventoryItem]:
    ipmap = {}
    for item in inventory_items:
        entity_data = item.entity_data
        if isinstance(entity_data, (VMData, ManagedServiceData)) and entity_data.nics:
            for nic in entity_data.nics:
                if nic.private_ip_addresses:
                    for ip in nic.private_ip_addresses:
                        ipmap[(ip, nic.network)] = get_inventory_item_info(item.item_id, item.entity_type, item.external_ids)
                if nic.public_ip_addresses:
                    for ip in nic.public_ip_addresses:
                        ipmap[(ip, nic.network)] = get_inventory_item_info(item.item_id, item.entity_type, item.external_ids)
    return ipmap


def fetch_ipmap() -> Dict[str, InventoryItem]:
    inventory = fetch_inventory()
    ipmap = popuplate_ipmap(inventory_items=inventory)
    return ipmap