from common.logger import log_init, get_logger


LOGGER_NAME = "test-use-common-py"

logger = log_init(
    base_name=LOGGER_NAME,
    # use_stderr=True,
    # show_thread_name=True,
    # use_colors=False
)

LOG = get_logger(module_name=__name__)

LOG.info("this is my first lline")