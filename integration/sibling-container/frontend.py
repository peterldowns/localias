from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, HTTPServer


FRONTEND_PORT = 5173  # Must match the upstream and healthcheck in compose.yaml.


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        # Echo the Host header so the proxy test can verify that it is preserved.
        body = self.headers["Host"].encode()
        self.send_response(HTTPStatus.OK)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


HTTPServer(("0.0.0.0", FRONTEND_PORT), Handler).serve_forever()
