# AI Prompt Log

The key prompts sent to the AI coding assistant while building this project, in order.

- **Prompts** appear in quote blocks, exactly as they were sent. The one exception is a personal email address, replaced with `[my email]`.
- **Context and Result** lines are added afterwards, to explain what each prompt was answering and what it produced.
- Commands starting with `/mattpocock-skills:` run a reusable skill, a packaged set of workflow instructions for the assistant.

## Conversations

| # | Focus | Prompts | Tool |
|---|---|---:|---|
| [1](#conversation-1-design-spec-and-first-ticket) | Design, spec, tickets, and the first working slice | 9 | Claude Code (Claude Opus 5.5) |
| [2](#conversation-2-remaining-tickets-and-final-check) | Tickets 02 to 08, built in parallel and verified | 3 | Claude Code (Claude Opus 5.5) |
| [3](#conversation-3-architecture-review-docker-compose-and-readme) | Code review, Docker Compose, coverage reports, and the full README | 1 | Claude Code (Claude Opus 5.5) |

---

## Conversation 1: Design, spec, and first ticket

**Date:** November 3, 2026

**Outcome:** The design was settled through a question-and-answer session, then written up as a README, a spec, and 8 tickets. Ticket 01 was built, reviewed, and committed.

| # | Prompt | Stage |
|---:|---|---|
| [1](#1-define-the-project) | Define the project | Design |
| [2](#2-answer-the-first-design-round) | Answer the first design round | Design |
| [3](#3-answer-the-second-design-round) | Answer the second design round | Design |
| [4](#4-write-the-readme) | Write the README | Documentation |
| [5](#5-write-the-spec) | Write the spec | Planning |
| [6](#6-break-the-spec-into-tickets) | Break the spec into tickets | Planning |
| [7](#7-approve-the-ticket-breakdown) | Approve the ticket breakdown | Planning |
| [8](#8-implement-ticket-01) | Implement ticket 01 | Implementation |
| [9](#9-prepare-the-next-conversation) | Prepare the next conversation | Handoff |

### 1. Define the project

> /mattpocock-skills:grill-me I need to build a full-stack calculator monorepo with a Go backend (/backend) and a React + TypeScript frontend (/frontend) using Vite. It needs to support basic arithmetic (add, subtract, multiply, divide) plus exponentiation, square root, and percentage.
>
> \- Frontend (React):  
> Intuitive UI for entering input and displaying results  
> Input validation and error handling  
> Responsive design (basic mobile support)
>
> \- Backend (REST API):  
> Expose endpoints for calculator operations  
> Validate input and handle edge cases (division by zero, invalid data)  
> Return results in JSON format
>
> \- Non-Functional  
> Clean, readable, and idiomatic code (frontend and backend)  
> Unit tests covering key functionality for both layers  
> Documentation: setup instructions, API usage, and design rationale  
> Dockerfile for full-stack deployment

**Result:** The assistant checked the repo and toolchain, then asked a first round of 8 design questions, each with a recommended answer.

### 2. Answer the first design round

**Context:** The questions covered the quality bar, interaction model, endpoint shape, percentage meaning, number representation, error contract, deployment shape, and dependencies.

> yes to all, lets use a go LTS version

**Result:** All recommendations were accepted. The assistant explained that Go has no LTS releases and proposed Go 1.27, the newest supported version, in a second round of 9 questions.

### 3. Answer the second design round

**Context:** The questions covered the Go and Node versions, backend layout, frontend serving, maths edge cases, form behaviour, client-side validation, testing, and docs and git.

> yes to all

**Result:** The design was complete, and the assistant summarised the full plan for confirmation.

### 4. Write the README

> create the README.md

**Result:** The root `README.md` was written with setup instructions, API usage, and design rationale, ahead of the code it describes.

### 5. Write the spec

> /mattpocock-skills:to-spec use the docs folder

**Result:** [`fullstack-calculator-spec.md`](fullstack-calculator-spec.md) was written with 61 user stories, implementation decisions, and testing decisions.

### 6. Break the spec into tickets

> /mattpocock-skills:to-tickets create tickets for every operation and save them in the docs folder

**Result:** The assistant proposed 9 tickets with blocking relationships and asked 4 questions before writing any files.

### 7. Approve the ticket breakdown

**Context:** The questions asked where errors with no specific input should appear, whether to keep the deployment and final-check tickets, whether the granularity was right, and whether the blocking relationships were correct.

> q1 yes, q2 8 and 9 are ok, q3 2 and 3 can be merged, q4 yes

**Result:** Subtract and multiply were merged, giving 8 tickets in [`tickets/`](tickets/). The spec was updated with the decision on where those errors appear.

### 8. Implement ticket 01

> /mattpocock-skills:implement ticket 01

**Result:** The add operation was built test-first through the whole stack. Two review agents found issues, and the ones that broke explicit requirements were fixed and committed.

### 9. Prepare the next conversation

> /mattpocock-skills:handoff the next section will be the implementation of the tickets that 01 unblocked, I'd like to implement the other tickets in parallel

**Result:** A handoff document was written outside the repo, including an analysis of where parallel tickets would conflict.

---

## Conversation 2: Remaining tickets and final check

**Date:** November 3, 2026

**Outcome:** Tickets 02 to 07 were built by parallel agents and reviewed, then ticket 08 verified the finished app in Docker and locally and corrected the README. All 8 tickets are done.

| # | Prompt | Stage |
|---:|---|---|
| [1](#1-implement-the-unblocked-tickets) | Implement the unblocked tickets | Implementation |
| [2](#2-start-docker-and-choose-the-404) | Start Docker and choose the 404 | Review |
| [3](#3-run-the-final-check) | Run the final check | Verification |

### 1. Implement the unblocked tickets

> /mattpocock-skills:implement using subagents the unblocked tickets by 01, first read the handoff document /var/folders/9y/2snjbqxx3gxfs2pr5c_ch89r0000gn/T/fullstack-calculator-handoff.md

**Result:** Agents in separate git worktrees built tickets 02, 03, 06, and 07, then 04 and 05, because power reuses the `division_by_zero` error that ticket 03 adds. Each ticket landed on `main` as its own commit. A first `Percentage` used `math/big`, which broke the spec's `float64` decision, so it was redone with `float64` only. Two review agents and a Chrome check followed, and the check found the operation picker wrapping on desktop, which was fixed. The Docker image was still unverified because the daemon wasn't running.

### 2. Start Docker and choose the 404

**Context:** The assistant asked whether Docker could be started, and whether `POST /missing` should return 404 when serving static files, where it returned 405. It recommended keeping the 405.

> 1. docker is running, listening on :8080  
> 2. 404 is ok for the /missing endpoint

**Result:** Unknown paths now return 404 for every method, and real files still reject other methods with 405. The image was built and tested for the first time: 86 HTTP checks passed, it ran as a non-root user, and it stopped cleanly.

### 3. Run the final check

> go ahead with task 08

**Result:** A fresh clone passed `make test`, `make lint`, and `make build`, and 100 `curl` checks passed against the container across every route. The README was corrected: exact response formats, the operation table wording, and two undocumented edge cases. The local setup couldn't be run exactly as written, because a container was using port 8080.

---

## Conversation 3: Architecture review, Docker Compose, and README

**Date:** November 3, 2026

**Outcome:** A code review led to three small cleanups that didn't change behaviour. Coverage reports, separate backend and frontend images, and a Compose file were added, and the README was rewritten.

| # | Prompt | Stage |
|---:|---|---|
| [1](#1-review-the-code-and-add-docker-compose) | Review the code and add Docker Compose | Review and deployment |

### 1. Review the code and add Docker Compose

> /mattpocock-skills:improve-codebase-architecture Do a quick review of both /backend and /frontend to make sure the code is idiomatic, modules have clean boundaries, and there's no unnecessary complexity or duplication. Prioritize correctness, clarity, and maintainability over extra features.
>
> Once everything is clean and all Go and Vitest tests pass with coverage reports:  
> 1\. Add multi-stage Dockerfiles for the backend and frontend, plus a root docker-compose.yml so the entire stack can be spun up with a single command.  
> 2\. Make sure that the README.md is comprehensive covering: prerequisites (Go, Node, Docker), local setup and run instructions, how to run tests and coverage reports, curl examples for the REST API (both success and error cases), key design decisions.

**Result:** The review flattened `staticFiles`, merged the repeated validation test tables into one that runs against all six two-operand routes, and removed two unneeded type exports. `make coverage` was added, along with `@vitest/coverage-v8`. The stack was tested through Compose on both ports, including replacing the backend container and stopping and restarting it. The README's 18 `curl` examples with documented output were run against the server and all matched. It wasn't checked in a real browser. The review also suggested removing the root `Dockerfile` and Go static serving if Compose was enough.
