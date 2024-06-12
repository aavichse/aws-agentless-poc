import os
from pydantic import BaseModel, Field, ValidationError, model_serializer
from .vpclog_reader import FlowRecord
from common.logger import get_logger
from typing import Callable, Dict, Optional, Tuple, Union
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
    EnforcementInfo,
    Verdict
)

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

    class Config:
        extra = 'ignore'  # Ignore reported fields not mentioned in model

    # # Added to base model: check no double process info and convert event_type & direction to ConnectionEventType
    # @model_serializer(mode='wrap')
    # def serialize_model(self, handler):
    #     # if self.source_process_info and self.dest_process_info:
    #     #     raise ValidationError("source process info and dest process info can't co-exist")
    #     self.event_type = {
    #             (EventType.SUCCESSFUL.value, Direction.OUTBOUND.value): ConnectionEventType.NewSuccessOutgoingConnection,
    #             (EventType.SUCCESSFUL.value, Direction.INBOUND.value): ConnectionEventType.NewSuccessIncomingConnection,
    #             (EventType.SUCCESSFUL.value, Direction.BI_DIRECTIONAL.value): ConnectionEventType.NewSuccessMatchedConnection,
    #             (EventType.SUCCESSFUL.value, Direction.OUTBOUND_ONLY.value): ConnectionEventType.NewSuccessOutgoingConnection,
    #             (EventType.FAILED.value, Direction.OUTBOUND.value): ConnectionEventType.NewFailedOutgoingConnection,
    #             (EventType.FAILED.value, Direction.INBOUND.value): ConnectionEventType.NewFailedIncomingConnection,
    #             (EventType.FAILED.value, Direction.BI_DIRECTIONAL.value): ConnectionEventType.NewFailedMatchedConnection,
    #             (EventType.FAILED.value, Direction.OUTBOUND_ONLY.value): ConnectionEventType.NewFailedOutgoingConnection,
    #         }.get((self.event_type.value, self.direction.value))
    #     return handler(self)


LOG = get_logger(module_name=__name__)
MSG_LOG = get_logger(module_name=__name__, logger_name='vpcflow')

GCAPP_FLOWLOGS_TOPIC = 'gcapp-flowlogs-ehub'
GCAPP_BROKER = os.getenv('GCAPP_BROKER', '') + ':9093'

producer = KafkaProducer(
    bootstrap_servers=[GCAPP_BROKER],
    value_serializer=lambda x: x.model_dump_json(exclude_none=True, by_alias=True).encode('utf-8'))

# producer = KafkaProducer(
#	bootstrap_servers=['ec2-44-204-113-186.compute-1.amazonaws.com:9093'],
#    value_serializer=lambda x: x.model_dump_json(exclude_none=True, by_alias=True).encode('utf-8'))


# UNKNOWN_INFO_ITEM = InventoryItem.model_validate(
#             {
#                 'item-id': 'unkown', 
#                 'item-type': ItemType.ASSET,
#                 'external-ids': ['unkown'],   
#             })

class Reveal: 

    def __init__(self, ipmap: Dict[str, InventoryItem], reporting_entity_id: str):
        self.ipmap = ipmap
        self.reporting_entity_id = reporting_entity_id

    def filter_by_tcp_flags(self, rec: FlowRecord) -> bool: 
        return rec.tcp_flags != 2 and rec.tcp_flags != 3

    def resolve_ip(self, ip: str, network: str) -> InventoryItem:
        return self.ipmap.get((ip, network), None)

    def calc_direction(self, 
                       record: FlowRecord, 
                       src_item_info: InventoryItem, 
                       dst_item_info: InventoryItem) -> Direction: 
        if src_item_info and dst_item_info: 
            return Direction.BI_DIRECTIONAL
        if src_item_info:
            return Direction.OUTBOUND
        if dst_item_info: 
            return Direction.INBOUND
        
        return None
    

    def calc_event_type(self, rec: FlowRecord, direction: Direction):
        if rec.action == "ACCEPT": 
            if direction == Direction.OUTBOUND:
                return ConnectionEventType.NewSuccessOutgoingConnection
            if direction == Direction.INBOUND:
                return ConnectionEventType.NewSuccessIncomingConnection
            if direction == Direction.BI_DIRECTIONAL:
                return ConnectionEventType.NewSuccessMatchedConnection
            
        # FAILED
            if direction == Direction.OUTBOUND:
                return ConnectionEventType.NewFailedOutgoingConnection
            if direction == Direction.INBOUND:
                return ConnectionEventType.NewFailedIncomingConnection
            if direction == Direction.BI_DIRECTIONAL:
                return ConnectionEventType.NewFailedMatchedConnection
        
        
    #     self.event_type = {
    #             (EventType.SUCCESSFUL.value, Direction.OUTBOUND.value): ConnectionEventType.NewSuccessOutgoingConnection,
    #             (EventType.SUCCESSFUL.value, Direction.INBOUND.value): ConnectionEventType.NewSuccessIncomingConnection,
    #             (EventType.SUCCESSFUL.value, Direction.BI_DIRECTIONAL.value): ConnectionEventType.NewSuccessMatchedConnection,
    #             (EventType.SUCCESSFUL.value, Direction.OUTBOUND_ONLY.value): ConnectionEventType.NewSuccessOutgoingConnection,
    #             (EventType.FAILED.value, Direction.OUTBOUND.value): ConnectionEventType.NewFailedOutgoingConnection,
    #             (EventType.FAILED.value, Direction.INBOUND.value): ConnectionEventType.NewFailedIncomingConnection,
    #             (EventType.FAILED.value, Direction.BI_DIRECTIONAL.value): ConnectionEventType.NewFailedMatchedConnection,
    #             (EventType.FAILED.value, Direction.OUTBOUND_ONLY.value): ConnectionEventType.NewFailedOutgoingConnection,
    #         }.get((self.event_type.value, self.direction.value))
    #     return handler(self)


    def send(self, rec: FlowRecord) -> ConnectionInfo:

        if rec.dstport > rec.srcport: 
            return

        src_item = self.resolve_ip(rec.srcaddr, rec.vpc_id)
        dst_item = self.resolve_ip(rec.dstaddr, rec.vpc_id)
        
        if src_item is None and dst_item is None: 
            return 
        
        direction = self.calc_direction(record=rec, src_item_info=src_item, dst_item_info=dst_item)
        
        msg = ConnectionInfo.model_validate({
            'direction': direction,
            'event-type': EventType.SUCCESSFUL.value if rec.action == 'ACCEPT' else EventType.FAILED.value,
            'source-ip': rec.srcaddr,
            'dest-ip': rec.dstaddr,
            'dest-port': rec.dstport,
            'ip-protocol': IpProtocol.TCP,  # FIXME
            'ip-version': IpVersion.IPV4,
            'enforcement-state': EnforcementState.REVEAL_ONLY,
            'source-inventory-item': src_item,
            'dest-inventory-item': dst_item,
            'start-time': 1718010000, #rec.start,
            'end-time': 1718014423, #rec.end,
            'count': rec.count,
            'reporting-entity': ReportingEntity(uuid=self.reporting_entity_id, type='cloud_aws'),
        })

        LOG.info(f'PUBLISH: {msg.model_dump_json(by_alias=True, exclude_none=True)}')
        
        producer.send(GCAPP_FLOWLOGS_TOPIC, value=msg)
        producer.flush() 

