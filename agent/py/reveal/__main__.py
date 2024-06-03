from argparse import ArgumentParser
import sys
from retry import retry
from .vpclog_worker import VPCLogWorker
from .ipmap_fetcher import fetch_ipmap
from common.logger import log_init, get_logger

# DELME
# VPCLOG_S3 = "nabcert-tooling-nonprod-v1975-flow-log"
# VPCLOG_S3_BUCKET = "nabcert-private-nonprod-v1975-flow-log"

REPORTING_ENTITY_ID = "cf178df3-1d8c-46a5-86b7-974a941c4d80"
LOG = None

def init_logger(args):
    global LOG
    log_init(
        base_name="REVEAL-" + args.s3_bucket_name,
        use_stderr=True,
        show_thread_name=True,
        log_path="/var/log/agentless",
    )
    LOG = get_logger(module_name=__name__)


def main(argv=None):
    argv = argv or sys.argv[1:]
    parser = ArgumentParser(description="Read VPC Flow Log Records")

    parser.add_argument(
        "--s3-bucket-name",
        type=str,
        help="S3 bucket name of VPC flowlog",
    )
    parser.add_argument(
        "--interval",
        type=int,
        default=60,
        help="Interval fetching logflows (default=60s)",
    )
    parser.add_argument(
        '--dump',
        action='store_true',
        help='Dump read records',
    )
    args = parser.parse_args(argv)

    init_logger(args)

    LOG.info("start VPC flowlog={args.s3_bucket_name}")

    vpclog = VPCLogWorker(
        vpclog_s3_bucket_name=args.s3_bucket_name,
        interval=args.interval,
        reporting_entity_id=REPORTING_ENTITY_ID,  # FIXME: hardcoded for PoC
        dump=args.dump
    )

    vpclog.start()


if __name__ == "__main__":
    main()
