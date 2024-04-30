import os
import time
from flask import Flask
from common.logger import log_init, get_logger
from .orchestrator import AWSOrchestrator
from .enforcement import Enforcement
from .routes import route
from .vpclog_worker import VPCLogWorker
from .vpclog_reveal import Reveal

LOGGER_NAME = "aws-agent"
VPCLOG_S3 = "gca-flowlogs"

log_init(
    base_name=LOGGER_NAME, use_stderr=True, show_thread_name=True, log_path=os.getcwd()
)

log_init(
    base_name="vpcflow",
    log_path=os.getcwd(),
    show_thread_name=False,
    use_stderr=False,
    logger_name="vpcflow",
)

LOG = get_logger(module_name=__name__)

app = Flask(__name__)

if __name__ == "__main__":

    # orchestrator
    orchestrator = AWSOrchestrator(regions=["us-east-1"])
    orchestrator.start()

    while orchestrator.inventory.version == 0:
        LOG.info("Waiting for initial inventory....")
        time.sleep(5)

    def get_inventory():
        return orchestrator.inventory

    enforcement = Enforcement(10, get_inventory)
    enforcement.start()

    def get_plan():
        return enforcement.plan

    # POC: is for single region within a single account.
    reveal = Reveal(get_inventory_handler=get_inventory, get_plan_handler=get_plan)
    vpclog = VPCLogWorker(
        region="us-east-1", location=VPCLOG_S3, interval=60, reveal=reveal
    )
    vpclog.start()

    route(app, orchestrator, enforcement)
    app.run(debug=False, host="0.0.0.0", port=8080)
