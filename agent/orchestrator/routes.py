import json
from flask import Flask, Response, jsonify, request
from .aws_orchestrator import AWSOrchestrator
from common.logger import get_logger

count = 0
LOG = get_logger(module_name=__name__)

def route(api: Flask, orchestrator: AWSOrchestrator):

    @api.route('/', defaults={'path': ''})
    @api.route('/<path:path>', methods=['GET', 'POST', 'PUT', 'DELETE'])
    def catch_all(path: str):
        global count
        count += 1
        LOG.info(f'Method: {request.method}, Path: {request.path}')

        #headers = {k: v for k, v in request.headers.items()}
        #LOG.info(f'Headers: {headers}')

        try:
            data = request.get_json()
            body = json.dumps(data, indent=4)
        except Exception as e:
            body = request.get_data(as_text=True) if request.data else 'No body'

        #LOG.info(f'BODY: {body}')

        return jsonify({'message': f'#{count}'}), 200


    @api.route('/inventory/v1/provider/inventory', methods=['GET'])
    def get_inventory():
        inventory= orchestrator.inventory

        LOG.info(f'Request inventory: Total items={len(inventory.items)}')
        inventory_items = [item.model_dump_json(by_alias=True) for item in inventory.items]
        inventory_json = '[' + ','.join(inventory_items) + ']'
        inventory_bytes = inventory_json.encode('utf-8')
        return Response(inventory_bytes, mimetype='application/json', status=200)