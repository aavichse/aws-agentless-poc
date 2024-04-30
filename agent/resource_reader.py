from typing import Any
import boto3
from .resource_converter import AWSResourceConverter
from common.logger import get_logger

LOG = get_logger(module_name=__name__)


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
        
        
class AWSResourceReader: 
    def __init__(self, session: boto3.Session):
        self.session = session
        self.ec2_client = session.client('ec2')
        
    def try_get_exists_security_group(self, vpc_id: str, name: str) -> Any:
        try:
            response = self.ec2_client.describe_security_groups(
                Filters=[
                    {
                        'Name': 'vpc-id',
                        'Values': [vpc_id]
                    },
                    {
                        'Name': 'group-name',
                        'Values': [name]
                    }
                ]
            )
            if response['SecurityGroups']:
                return response['SecurityGroups'][0]
        except Exception as e:
            return None
    
    
    def try_create_security_group(self,
                              vpc_id: str,
                              name: str,
                              description: str = None) -> str:
        LOG.info(f'Creating security group {vpc_id=}, {name=}, {description=}')
        sg = self.try_get_exists_security_group(vpc_id, name)
        if sg: 
            LOG.info(f'Security group already exists: {sg["GroupId"]}')
            return sg
         
        try:
            sg = self.ec2_client.create_security_group(
                GroupName=name,
                Description=description or f'CloudApp generated security group - f{name}',
                VpcId=vpc_id,
            )

            id = sg['GroupId']
            
            self.ec2_client.create_tags(
                Resources=[id],
                Tags=[
                    {'Key': 'Name', 'Value': name}
                ]
            )

            LOG.info(
                f'Security group created: {id=}, {vpc_id=}, {name=}, {description=}')

            return sg

        except Exception as e:
            LOG.exception(
                f'Failed create security group {vpc_id=}, {name=}, {e=}')
            return None
        
    def revoke_ingress_rules(self, sg):
        ip_permissions = sg['IpPermissions']
        sg_id=sg['GroupId']
        if ip_permissions:
            try:
                self.ec2_client.revoke_security_group_ingress(
                    GroupId=sg_id, IpPermissions=ip_permissions)
                LOG.info(
                    f"Revoked ingress rules for security group {sg_id}")
            except Exception as e:
                LOG.info(
                    f"Error revoking ingress rules for security group {sg_id}: {e}")
        else:
            LOG.info(
                f"No ingress rules to revoke for security group {sg_id}")
            
            
    def get_security_groups_by_prefix(self, prefix):
        # Retrieve all security groups that have names starting with the specified prefix
        response = self.ec2_client.describe_security_groups(Filters=[
            {'Name': 'group-name', 'Values': [prefix + '*']}
        ])
        return response['SecurityGroups']
    
    def delete_gc_security_groups(self):
        sgs = self.get_security_groups_by_prefix('gc_')
        for sg in sgs: 
            self.revoke_ingress_rules(sg)
            
        for sg in sgs:
            self.delete_security_group(sg['GroupId'])
        
        return [i['GroupId'] for i in sgs]
    
    def delete_security_group(self, id: str) -> bool:

        ec2_client = boto3.client('ec2', region_name='us-east-1')
        ec2_resource = boto3.resource('ec2', region_name='us-east-1')

        response = ec2_client.describe_instances(
            Filters=[
                {
                    'Name': 'instance.group-id',
                    'Values': [id]
                }
            ]
        )

        instances_to_modify = []
        for reservation in response['Reservations']:
            for instance in reservation['Instances']:
                instances_to_modify.append(instance['InstanceId'])

        for instance_id in instances_to_modify:
            instance = ec2_resource.Instance(instance_id)
            # Get current security groups attached to the instance, excluding the target security group
            new_sg_ids = [sg['GroupId']
                          for sg in instance.security_groups if sg['GroupId'] != id]
            # If instance is associated with more security groups, modify the instance to use the new list
            if new_sg_ids:
                instance.modify_attribute(Groups=new_sg_ids)

        # Step 3: Delete the security group
        try:
            # self.delete_all_inbound_rules(security_group_id)
            ec2_client.delete_security_group(GroupId=id)
            print(
                f"Security Group {id} has been detached from all instances and deleted.")
        except Exception as e:
            if 'does not exists' in str(e):
                return True
            if 'depend' in str(e):
                return False
            print(f"Error deleting security group {id}: {e}")
            raise e

        return True

    def create_allow_ingress(self,
                             src_sg: str,
                             dst_sg: str,
                             port: int, 
                             description: str = ''):

        src_sg_id = src_sg['GroupId']
        dst_sg_id = dst_sg['GroupId']
        
        try:
            self.ec2_client.authorize_security_group_ingress(
                GroupId=dst_sg_id,
                IpPermissions=[
                    {
                        'IpProtocol': 'tcp',
                        'FromPort': port,
                        'ToPort': port,
                        'UserIdGroupPairs': [
                            {
                                'GroupId': src_sg_id,
                                'Description': description
                            }
                        ]
                    }
                ]
            )
            # LOG.info(
            #     f'Created ingress rule from {src_sg_id=} to {dst_sg_id=} port={port}')
        except Exception as e:
            LOG.error(
                f'Error creating ingress rule from {src_sg_id=} to {dst_sg_id=}: {e}')
            
    def update_instance_with_security_groups(self, id: str, security_groups: list[str]):
        try:
            self.ec2_client.modify_instance_attribute(
                InstanceId=id,
                Groups=security_groups
            )
            LOG.info(
                f"Updated instance {id=} with {security_groups=}")
        except Exception as e:
            LOG.error(
                f"Error updating instance {id=} with security group {security_groups=}: {e}")