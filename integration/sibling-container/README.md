# Sibling container integration test

From the repository root, run:

```sh
docker compose -f integration/sibling-container/compose.yaml up --build --abort-on-container-exit --exit-code-from test
```

The test container starts embedded Caddy with a generated HTTP alias and proxies
through Docker DNS to `frontend:5173`. The frontend returns the received Host
header so the test verifies it is preserved. No host ports, hosts-file changes,
certificates, or HAProxy are needed. Afterwards, remove the test containers and
network with:

```sh
docker compose -f integration/sibling-container/compose.yaml down
```
