import os
from retry import retry
from .vpclog_worker import VPCLogWorker
from .ipmap_fetcher import fetch_ipmap
from common.logger import log_init, get_logger

LOGGER_NAME = "reveal"
VPCLOG_S3 = "poc-examples"
REPORTING_ENTITY_ID = 'cf178df3-1d8c-46a5-86b7-974a941c4d80'

log_init(
    base_name=LOGGER_NAME,
    use_stderr=True,
    show_thread_name=True,
    log_path="/var/log/guardicore",
)

log_init(
    base_name=LOGGER_NAME,
    log_path="/var/log/guardicore",
    show_thread_name=False,
    use_stderr=False,
    logger_name="vpcflow",
)

LOG = get_logger(module_name=__name__)


if __name__ == "__main__":

    LOG.info("start")
    
    vpclog = VPCLogWorker(
        region="us-east-1", location=VPCLOG_S3, interval=60, 
        reporting_entity_id=REPORTING_ENTITY_ID
    )
    vpclog.start()
