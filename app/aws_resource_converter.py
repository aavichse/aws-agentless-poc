from typing import Any, Dict
from aggregator.model.integrations.v1.common.inventory import (
    InventoryItem,
    NetworkInterfaceData,
    ItemType,
    Label,
    Type,
    VMData,
)


class AWSResourceConverter:

    @classmethod
    def ec2_to_InventoryItem(cls, ec2, vm_data: VMData) -> InventoryItem:
        entity_name = vm_data.os_details["tags"].get("Name", ec2["InstanceId"])
        return InventoryItem.model_validate(
            {
                "item-id": ec2["InstanceId"],
                "item-type": ItemType.ASSET,
                "entity-name": entity_name,
                "external-ids": [ec2["InstanceId"]],
                "entity-type": "virtual machine",
                "entity-category": "compute",
                "entity-data": vm_data,
                "labels": [
                    Label(key=item["Key"], value=item["Value"])
                    for item in ec2.get("Tags", [])
                ],
            }
        )

    @classmethod
    def nic_to_NetworkInterfaceData(cls, nic) -> NetworkInterfaceData:
        public_ips = []
        for ip_entry in nic["PrivateIpAddresses"]:
            if "Association" in ip_entry:
                public_ips.append(ip_entry["Association"]["PublicIp"])

        return NetworkInterfaceData.model_validate(
            {
                "id": nic["NetworkInterfaceId"],
                "mac-address": nic["MacAddress"],
                "private-ip-addresses": [nic["PrivateIpAddress"]],
                "public-ip-addresses": public_ips,
                "network": nic["VpcId"],
                "subnet-id": nic["SubnetId"],
            }
        )

    @classmethod
    def to_VMData(
        cls,
        nic: NetworkInterfaceData,
        private_ip: str,
        security_groups: list[str],
        tags: Dict[str, Any],
    ):

        return VMData.model_validate(
            {
                "type": Type.VM,
                "os-details": {
                    "private_ip": private_ip,
                    "security_groups": security_groups, 
                    "tags": tags
                },
                "nics": [nic]
            }
        )
