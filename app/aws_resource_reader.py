from typing import Any
import boto3
from .aws_resource_converter import AWSResourceConverter


class AWSResourcePaginationReader:

    def __init__(self, describe_api: str, session: boto3.Session, filters: list = []):
        self.describe_api = describe_api
        self.filters = filters
        self.session = session
        self.iterator = self._reader()

    def __iter__(self):
        return self

    def __next__(self):
        return next(self.iterator)

    def _reader(self):
        self.ec2_client = self.session.client("ec2")
        paginator = self.ec2_client.get_paginator(self.describe_api)
        page_iterator = paginator.paginate(Filters=self.filters)
        for page in page_iterator:
            for item in self.describe(page):
                yield item

    def describe(self, page):
        raise NotImplementedError()


class AWSEC2Reader(AWSResourcePaginationReader):

    def __init__(self, session: boto3.Session, filters: list = []):
        super().__init__("describe_instances", session, filters)

    def describe(self, page):
        for reservation in page["Reservations"]:
            for ec2 in reservation["Instances"]:
                yield AWSResourceConverter.ec2_to_InventoryItem(
                    ec2=ec2, 
                    vm_data = self._get_vm_data(ec2)
                )

    def _get_vm_data(self, ec2) -> Any:
        # Note: for the POC we will take only the FIRST eni 
        eni = ec2['NetworkInterfaces'][0]
        network_interface = AWSResourceConverter.nic_to_NetworkInterfaceData(eni)

        return AWSResourceConverter.to_VMData(
            nic=network_interface, 
            private_ip=network_interface.private_ip_addresses[0], 
            security_groups=[group['GroupId'] for group in eni['Groups']], 
            tags={item['Key']: item['Value'] for item in ec2.get('Tags', [])}
        )