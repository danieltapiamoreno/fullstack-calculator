# 07: Deployable image

**What to build:** A reviewer or operator can start the whole app with one `docker run` and use the frontend and API from a single port, without installing Go or Node.

**Blocked by:** 01 (Add: first working calculation)

**Status:** ready-for-agent

- [ ] One multi-stage root Dockerfile builds the frontend and backend into one image, using the same Go and Node versions as the tool pins
- [ ] The final image is distroless and runs as non-root
- [ ] `make docker-run` builds the image and runs it on port 8080
- [ ] The app and API work from one origin in the container, with no CORS configuration
- [ ] The server serves built frontend files from `STATIC_DIR`, and runs API-only when it's unset
- [ ] `PORT` sets the listening port and defaults to 8080
- [ ] `GET /healthz` returns 200 and `{"status": "ok"}`
- [ ] The server has timeouts configured and shuts down gracefully, letting in-flight requests finish
- [ ] API tests cover the health check and static serving, including API-only mode
