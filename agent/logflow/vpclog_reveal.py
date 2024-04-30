import datetime
from pydantic import BaseModel, Field
from .vpclog_reader import FlowRecord
from common.logger import get_logger
from typing import Callable, Optional, Union
from kafka import KafkaProducer
from aggregator.model.integrations.v1.common.inventory import InventoryItem, ItemType
from aggregator.model.integrations.v1.provider.reveal import (
    Direction, 
    EventType, 
    IpProtocol,
    IpVersion,
    EnforcementState, 
    ProcessInfo,
    ReportingEntity, 
    ConnectionEventType, 
    NatInfo,
    EnforcementInfo
)
from ..orchestrator.aws_orchestrator import AwsInventory


class ConnectionInfo(BaseModel):

    direction: Union[Direction, None] = Field(
        ...,
        alias='direction',
        description='Indicates the direction of the network event; bi-directional is reported for two-sided connection;'
                    ' outbound-only - Destination is PaaS resource behind a private endpoint. ',
    )
    # Changed from base model: instead of just EventType, we later convert it to ConnectionEventType
    event_type: Union[EventType, ConnectionEventType] = Field(
        ...,
        alias='event-type',
        description='Type of network event that occurred; redirected is used for honeypot',
    )
    source_ip: str = Field(
        ..., alias='source-ip', description='Source IP address of the network event.'
    )
    destination_ip: str = Field(
        ..., alias='dest-ip', description='Destination IP address of the network event.'
    )
    destination_port: int = Field(
        ...,
        alias='dest-port',
        description='Destination port number for the network event.',
    )
    ip_protocol: IpProtocol = Field(
        ..., alias='ip-protocol', description='IP protocol used in the network event.'
    )
    ip_version: IpVersion = Field(
        ..., alias='ip-version', description='IP version used in the network event.'
    )
    enforcement_state: EnforcementState = Field(
        ...,
        alias='enforcement-state',
        description="Specifies the enforcement level for a network event: 'monitoring' for policy evaluation but"
                    " allows traffic, and 'reveal only' for data reporting without policy evaluation."
    )
    # Change from base model: changed name from enforcement-info to just enforcement
    enforcement: Optional[EnforcementInfo] = Field(
        None,
        alias='enforcement-info',
        description='Additional enforcement information.',
    )
    nat_info: Optional[NatInfo] = Field(
        None,
        alias='nat-info',
        description='Information about Network Address Translation.',
    )
    source_process_info: Optional[ProcessInfo] = Field(
        None,
        alias='source-process-info',
        description='Information about the source process involved in the network event.',
    )
    dest_process_info: Optional[ProcessInfo] = Field(
        None,
        alias='dest-process-info',
        description='Information about the destination process involved in the network event.',
    )
    # Change from base model: source-inventory-item -> source_inventory_item_info
    source_inventory_item_info: Optional[InventoryItem] = Field(
        None,
        alias='source-inventory-item',
        description='Inventory information for the source in the network event.',
    )
    # Change from base model: dest-inventory-item -> destination_inventory_item_info
    destination_inventory_item_info: Optional[InventoryItem] = Field(
        None,
        alias='dest-inventory-item',
        description='Inventory information for the destination in the network event.',
    )
    source_username: Optional[str] = Field(
        None,
        alias='source-username',
        description='Username associated with the source in the network event.',
    )
    dest_username: Optional[str] = Field(
        None,
        alias='dest-username',
        description='Username associated with the destination in the network event.',
    )
    dest_domain: Optional[str] = Field(
        None,
        alias='dest-domain',
        description='Domain associated with the destination in the network event (FQDN)',
    )
    # Change from base model: start-time -> bucket_start_time
    bucket_start_time: int = Field(
        ...,
        alias='start-time',
        description='Start time of the network event.** for aggregated events (epoch in seconds)',
    )
    # Change from base model: end-time -> bucket_end_time
    bucket_end_time: int = Field(
        ...,
        alias='end-time',
        description='End time of the network event.** for aggregated events (epoch in seconds)',
    )
    count: int = Field(..., description='Number of occurrences of the network event.')
    reporting_entity: ReportingEntity = Field(
        ...,
        alias='reporting-entity',
        description='Entity that reported the network event.',
    )




LOG = get_logger(module_name=__name__)
MSG_LOG = get_logger(module_name=__name__, logger_name='vpcflow')

GCAPP_FLOWLOGS_TOPIC = 'gcapp-flowlogs-ehub'
producer = KafkaProducer(bootstrap_servers=['ec2-44-200-243-86.compute-1.amazonaws.com:9093'],
                          value_serializer=lambda x: x.model_dump_json(exclude_none=True, by_alias=True).encode('utf-8'))

class Reveal: 
    inventory: AwsInventory = AwsInventory()
    

    def __init__(self, get_inventory_handler: Callable):
        self.get_inventory = get_inventory_handler
        self.update_inventory()
        self.unkown_item_info = InventoryItem.model_validate(
            {
                'item-id': 'unkown', 
                'item-type': ItemType.ASSET,
                'external-ids': ['unkown'],   
            })

    def filter_by_tcp_flags(self, rec: FlowRecord) -> bool: 
        return rec.tcp_flags != 2 and rec.tcp_flags != 3

    def update_inventory(self): 
        self.inventory = self.get_inventory()

    def get_item_info_from_ip(self, ip: str) -> InventoryItem:
        item = self.inventory.private_ips_lookup.get(ip, None)
        if not item: 
            return self.unkown_item_info
        
        return InventoryItem.model_validate(
                {
                    'item-id': item.item_id,
                    'item-type': ItemType.ASSET,
                    'external-ids': [item.item_id],
                })


    def send(self, rec: FlowRecord) -> ConnectionInfo:
        
        if self.filter_by_tcp_flags(rec): 
            return None

        MSG_LOG.info(f'{rec.to_short_message()}')

        (enforce_info, enforce_state) = (None, None) 
        
        msg = ConnectionInfo.model_validate({
            'direction': Direction.INBOUND if rec.flow_direction == 'ingress' else Direction.OUTBOUND,
            'event-type': EventType.SUCCESSFUL if rec.action == 'ACCEPT' else EventType.FAILED,
            'source-ip': rec.srcaddr,
            'dest-ip': rec.dstaddr,
            'dest-port': rec.dstport,
            'destination_port': rec.dstport,
            'ip-protocol': IpProtocol.TCP,  # FIXME
            'ip-version': IpVersion.IPV4,
            'enforcement-state': EnforcementState.REVEAL_ONLY,
            'source-inventory-item': self.get_item_info_from_ip(rec.srcaddr),
            'dest-inventory-item': self.get_item_info_from_ip(rec.dstaddr),
            'start-time': rec.start,
            'end-time': rec.end,
            'count': 1,
            'enforcement-info': enforce_info,
            'enforcement-state': enforce_state or EnforcementState.REVEAL_ONLY,
            # FIXME
            'reporting-entity': ReportingEntity(uuid='xxx', type='broker'),
        })

        MSG_LOG.info(msg.model_dump_json(by_alias=True, exclude_none=True))

        producer.send(GCAPP_FLOWLOGS_TOPIC, value=msg)
        producer.flush() 