# 07: Deployable image

**What to build:** A reviewer or operator can start the whole app with one `docker run` and use the frontend and API from a single port, without installing Go or Node.

**Blocked by:** 01 (Add: first working calculation)

**Status:** done

- [x] One multi-stage root Dockerfile builds the frontend and backend into one image, using the same Go and Node versions as the tool pins
- [x] The final image is distroless and runs as non-root
- [x] `make docker-run` builds the image and runs it on port 8080
- [x] The app and API work from one origin in the container, with no CORS configuration
- [x] The server serves built frontend files from `STATIC_DIR`, and runs API-only when it's unset
- [x] `PORT` sets the listening port and defaults to 8080
- [x] `GET /healthz` returns 200 and `{"status": "ok"}`
- [x] The server has timeouts configured and shuts down gracefully, letting in-flight requests finish
- [x] API tests cover the health check and static serving, including API-only mode

Not yet verified: the Docker image hasn't been built or run, because the Docker daemon wasn't available during implementation. The server binary was verified directly, outside a container.
