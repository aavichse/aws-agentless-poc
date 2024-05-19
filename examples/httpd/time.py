def application(environ, start_response):
    import json
    from datetime import datetime
    import socket

    hostname = socket.gethostname()
    server_ip = socket.gethostbyname(hostname)
    current_time = datetime.now().isoformat()

    response = json.dumps({
        "current_time": current_time,
        "server_ip": server_ip
    })

    start_response('200 OK', [('Content-Type', 'application/json')])
    return [response.encode('utf-8')]
