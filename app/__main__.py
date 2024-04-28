import os
from flask import Flask
from common.logger import log_init, get_logger
from .aws_orchestrator import AWSOrchestrator
from .routes import route

LOGGER_NAME = "aws-agent"

logger = log_init(
    base_name=LOGGER_NAME,
    use_stderr=True,
    show_thread_name=True,
    log_path=os.getcwd()
)

LOG = get_logger(module_name=__name__)

app = Flask(__name__)

if __name__ == '__main__':
    orchestrator = AWSOrchestrator(regions=['us-east-1'])
    orchestrator.start()

    route(app, orchestrator)
    app.run(debug=False, host='0.0.0.0', port=8080)