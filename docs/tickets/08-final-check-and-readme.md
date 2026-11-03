# 08: Final check and README

**What to build:** The finished project is verified as a whole. A reviewer can trust that the app behaves as documented and that following the README works.

**Blocked by:** 02 (Subtract and multiply), 03 (Divide), 04 (Power), 05 (Square root), 06 (Percentage), 07 (Deployable image)

**Status:** ready-for-agent

- [ ] `make test`, `make lint`, and `make build` all pass
- [ ] The Docker image builds and starts, and the frontend loads from it
- [ ] Every API route is exercised with `curl` against the running container, including one success and each of its error cases, and the results are reported
- [ ] The local setup instructions in the README work as written
- [ ] The README is corrected wherever it differs from what was built, including paths, commands, error codes, and example responses
