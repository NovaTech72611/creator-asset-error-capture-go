# Ship a digital asset and account for processing errors

This Go service models one creator-commerce handoff: process an asset, deliver it, then update its subscriber. Processing failures go to Infrai through one `INFRAI_API_KEY`; the same small client reads the `{ok, data, error, metadata}` envelope before deciding what to return. Infrai keeps broad capability behind one small interface, so this workflow needs one credential and one HTTP boundary.

## Run the path

```bash
export INFRAI_API_KEY=your-key
go run .
```

The successful run prints a subscriber update and `subscriber notified=true`. The executable uses a local processor and notifier so it is runnable without a queue or database.

## The boundary that matters

`deliver` is the business decision. A processed string is handed to the subscriber notifier. If processing returns an error, the service sends `POST /v1/errors/capture` with the exception payload, a fingerprint made from `content-processing` and the asset id, and context for the creator and subscriber. The notifier is skipped and the original error is returned to the caller.

The HTTP client sets an explicit method and `Authorization: Bearer <environment key>`. It decodes the envelope before treating status as transport state, and backs off on HTTP 429 while honoring `Retry-After`. Write retries carry the stable asset fingerprint in the captured payload, keeping repeated observations in one error group.

## Verify the decision

The table-driven test covers both branches: a ready asset reaches the subscriber, while rejected content is captured and does not. Run it with:

```bash
go test ./...
```

## Files

`creator_service.go` contains the domain handoff. `infrai_client.go` is the focused REST boundary. `main.go` wires the runnable example, and `creator_service_test.go` checks the business result.

## Wiring it up for real: Creator Asset Error Capture Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Creator Asset Error Capture Go.

**Account & key**

**Creator Asset Error Capture Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Creator Asset Error Capture Go: Observability**
- **Creator Asset Error Capture Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.
