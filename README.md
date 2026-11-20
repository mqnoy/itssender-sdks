# itssender SDKs

Official typed client libraries for the itssender omnichannel messaging gateway. Both packages wrap the
same REST API and ship the Standard Webhooks v1.0.0 signature verification helper, model every monetary
value and quota as a pure integer (minor units), and expose timestamps as int64 Unix milliseconds.

## Packages

| Language | Directory | Install | Package | Runtime requirement |
| --- | --- | --- | --- | --- |
| TypeScript/JavaScript | `typescript/` | `pnpm add @itssender/sdk` | `@itssender/sdk` (ESM + CJS, .d.ts) | Node >=18 or any fetch-capable edge runtime |
| Go | `go/` | `go get github.com/mqnoy/itssender-sdks/go` | package `itssender` | Go >=1.22 |

## TypeScript quick start

```ts
import { ItssenderClient, verifyWebhookSignature } from '@itssender/sdk';

const client = new ItssenderClient({
  baseURL: 'https://itssender.imzalab.com',
  apiKey: process.env.ITSSENDER_API_KEY!,
});

const sent = await client.messages.sendText('dev_01j9k2m3n4', {
  to: '6281234567890',
  text: 'Hello from itssender!',
});
console.log(sent.message_id);

const envelope = verifyWebhookSignature({
  payload: rawBody,
  headers: request.headers,
  secret: process.env.ITSSENDER_WEBHOOK_SECRET!,
});
```

Signature verification is also importable from its own entry point: `import { verifyWebhookSignature } from
'@itssender/sdk/crypto';`.

## Go quick start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mqnoy/itssender-sdks/go"
)

func main() {
	client := itssender.NewClient(
		itssender.WithBaseURL("https://itssender.imzalab.com"),
		itssender.WithAPIKey(os.Getenv("ITSSENDER_API_KEY")),
	)

	resp, err := client.Messages.SendText(context.Background(), "dev_01j9k2m3n4", itssender.SendTextMessageRequest{
		To:   "6281234567890",
		Text: "Hello from itssender!",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.MessageID)
}
```

Every Go method takes a `context.Context` first, so cancellation and deadlines propagate to the underlying
HTTP request. `WithBearerToken`, `WithTimeout` and `WithRetry` are available as additional options.

## Repository layout

```text
sdks/
├── go/                     Go SDK (package itssender)
├── typescript/             TypeScript SDK (@itssender/sdk)
└── .github/workflows/      Continuous integration
```

## Development

```sh
cd typescript && pnpm install && pnpm run typecheck && pnpm run build && pnpm run test
cd go && go vet ./... && go test ./... -race
```

## Security

Report vulnerabilities privately through the repository security advisories; do not open public issues.

## License

Apache-2.0. See [LICENSE](LICENSE).