# Adding a microservice (or a new RPC)

This guide documents how the Calidum Rotae backend is structured so new
microservices can be added without reverse-engineering the wiring. It also
covers the slightly-longer path of adding a new workflow to the github
provider.

## Architecture recap

```
HTTP client
   │  POST /discord /email /command /github /github/user /cluster (X-API-KEY)
   ▼
calidum-rotae-service (gin, port 3000)
   │  cmd/calidum-rotae-service/server/http.go  →  pkg/calidum (CalidumClient)
   ▼  each method JSON-decodes the body into the matching proto request
   gRPC (plaintext by default, TLS via CERTIFICATE_FILE_PATH)
   ▼
<name>-provider (port 4000..8000)
   │  cmd/<name>-provider/server/operations.go implements the proto service
   ▼
external API (Resend, Discord webhook, GitHub Actions dispatch, ...)
```

Every provider follows the same skeleton:

- `api/<name>_provider.proto` — the contract
- `pkg/proto-gen/<name>-provider/` — generated Go (never edit by hand)
- `cmd/<name>-provider/` — `main.go` (cobra + viper), `config/config.go`,
  `server/grpc.go` (register service), `server/operations.go` (business logic),
  `Dockerfile`

## 1. Adding a brand-new microservice

1. **Define the contract** in `api/<name>_provider.proto`:

   ```proto
   syntax = "proto3";
   option go_package = "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/<name>-provider;<pkg>";

   package <pkg>;

   message MyRequest { string Field = 1; }
   message MyResponse { string Status = 1; }

   service <Pkg>Provider {
     rpc DoThing(MyRequest) returns (MyResponse);
   }
   ```

2. **Generate the Go code** with `make grpc` (or `nix develop -c make grpc`
   if protoc isn't on your PATH). This regenerates every `api/*.proto`.

3. **Server scaffolding** — copy `cmd/<name>-provider/` from any existing
   provider and adapt:
   - `server/grpc.go`: `ConfigureGrpc()` builds a `grpc.NewServer()` and calls
     `<pkg>.Register<Pkg>ProviderServer(grpcServer, s)`.
   - `server/operations.go`: `type Server struct { <pkg>.<Pkg>ProviderServer }`
     then implement each RPC. Embedding `UnimplementedXServer` makes
     unimplemented RPCs fail with `UNIMPLEMENTED` instead of panicking.
   - `main.go`: cobra flags (`--port`, `--loglevel`), `logger.Initialize`,
     `serverutils.NewGrpcServer(...).Run(ctx, ...)`.

4. **Gateway** — wire the new client through these four places:
   - `cmd/calidum-rotae-service/config/config.go`: add
     `--<name>_provider_hostname` / `--<name>_provider_port` flags.
   - `cmd/calidum-rotae-service/client/client.go`: add an
     `init<Name>ProviderClient` (copy an existing one) and return it from
     `InitFromViper`.
   - `pkg/calidum/calidum.go`: add the client to `CalidumService`,
     `Dependencies`, a method on the `CalidumClient` interface (e.g.
     `SendMyRpcRequest`) and its implementation (JSON-unmarshal the body into
     the proto request, call the gRPC stub, return `[]byte`).
   - `cmd/calidum-rotae-service/server/http.go`: register a route in
     `initHTTPServerHandler`, e.g.
     `g.POST("/mything", func(g *gin.Context) { ... })`. Follow the existing
     per-RPC pattern: one `xxxPostRequest` handler (auth check → read body →
     gRPC span → dispatch → JSON response) plus one `sendXxxRpcRequestWithSpan`
     wrapper that owns the gRPC span and logs success/failure. See
     `githubPostRequest` / `sendGithubRpcRequestWithSpan` and the discord,
     email, shell and cluster equivalents.

5. **Docker + CI**:
   - Add a service block to `docker-compose.yml` and `docker-compose.dev.yml`
     (the dev file also spins up the OTEL stack). Expose the gRPC port and set
     `environment:` from `$`[vars loaded from `.env`].
   - Add the image name to the `image` choices in
     `.github/workflows/main.yml`.
   - Add a `Dockerfile` under `cmd/<name>-provider/`.

6. **Tests**: `cmd/calidum-rotae-service/server/http_test.go` uses
   `fakeCalidumClient` — it must implement every `CalidumClient` method the
   router needs. Add your RPC to the fake and a `POST` route test. Run
   `go build ./... && go vet ./... && go test ./...` (CI runs these too).

## 2. Adding a new RPC to an existing provider

Same as above but only steps 1, 2, parts of 4 (a new `CalidumClient` method +
a new route), and 6. The provider server file gets one new method. On the gin
side, follow the per-RPC pattern from the architecture section: each RPC gets
its own `sendXxxRpcRequestWithSpan` wrapper and `xxxPostRequest` handler, so
spans and log lines are tagged with the specific RPC.

## 3. Github provider: dispatches vs dedicated RPCs

The github provider distinguishes *operations* from *workflows*:

- **`RequestDeployment`** — one RPC that picks a workflow to run via the
  `Workflow` field (`"outline"`, `"grav"`, ...). Add a workflow here when the
  new workflow fits the existing deployment path.
- **Dedicated RPCs** (e.g. `AddCedilleUser`) — a new RPC when the workflow has
  a different input shape from deployments.

### Adding a new deployment workflow

In `cmd/github-provider/server/operations.go`:

1. Add a `..._WORKFLOW_URL` const pointing at the workflow dispatch endpoint:
   `https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/<file>.yml/dispatches`.
2. Add a `"<name>": <URL>` entry to the `workflowURLs` map in
   `RequestDeployment`.
3. All deployment workflows share the same inputs as `grav` (`nom_club` +
   `domaine`); they are sent to `dispatchWorkflow` in a single map. Input keys
   must match each workflow's `workflow_dispatch.inputs` **exactly** or GitHub
   returns `422`. If a future workflow needs different inputs, give it a
   dedicated RPC (like `AddCedilleUser`) instead.

### `AddCedilleUser` (example of a dedicated RPC)

Dispatches the `add-new-member.yml` workflow (const `CEDILLE_USER_WORKFLOW_URL`
in `cmd/github-provider/server/operations.go`) via the shared
`dispatchWorkflow`.

Workflow inputs (must match `workflow_dispatch` in the YAML):

| proto field          | workflow input       | JSON body key       |
|----------------------|----------------------|---------------------|
| `GithubUsername`     | `github_username`    | `githubUsername`    |
| `GithubEmail`        | `github_email`       | `githubEmail`       |
| `TeamSre`            | `team_sre`           | `teamSre`           |
| `ClusterRole`        | `cluster_role`       | `clusterRole`       |
| `NetdataRole`        | `netdata_role`       | `netdataRole`       |
| `UserUID`            | (caller identifier)  | `userUID`           |

Booleans are sent as workflow inputs as `"true"`/`"false"` strings
(`strconv.FormatBool`). Curl:

```bash
curl -X POST localhost:3000/github/user \
  -H "X-API-KEY: $CALIDUM_ROTAE_SERVICE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "githubUsername": "alexvegas22",
    "githubEmail": "alex@cedille.club",
    "teamSre": false,
    "clusterRole": "Reader",
    "netdataRole": "observer",
    "userUID": "2"
  }'
```

HTTP routes: `POST /github` (deployment, body carries `"workflow"`), and
`POST /github/user` (add member). Both return the workflow run URL in
`{"response": "..."}`.

### GitHub dispatch semantics baked into `dispatchWorkflow`

- On success, GitHub returns **204 No Content** with no body. The dispatch URL
  (from the response's `Location` header, if GitHub sets one) is returned in
  `WorkflowRunUrl`; there is no run metadata to fetch — a 204 is an
  acknowledgment, not a resource.
- Non-2xx responses (e.g. `422` for missing/unexpected inputs, `401` for a bad
  token) are returned as errors, so callers see why the dispatch failed.
- Requires the `GITHUB_TOKEN` env var (wired from `.env` in docker-compose).