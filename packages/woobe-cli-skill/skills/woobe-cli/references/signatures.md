# Server signatures and explicitly pinned trust

Use CLI 0.31+ and a backend with instance receipt signing configured. Read-only
export uses the authorized lifecycle owner's exact record; no receipt body is
accepted from the client.

```sh
woobe agent UUID receipt publication PUBLICATION_UUID --path publication.yaml
woobe network UUID receipt deployment DEPLOYMENT_UUID --env production --path deployment.yaml
```

`revision` takes a revision ID; `evaluation`, `publication`, `deployment` take
receipt UUIDs. Existing files are not overwritten. YAML/JSON follows extension.
Native UUIDs do not require `.woobe-config`.

Pin only a fingerprint confirmed independently by the instance operator:

```sh
woobe trust pin publication.yaml --root-fingerprint ROOT_FINGERPRINT --path trust.yaml
woobe receipt verify publication.yaml --trust trust.yaml
woobe trust refresh fresh-publication.yaml --trust trust.yaml --path trust-next.yaml
```

The uppercase values are placeholders, including the operator-confirmed
`sha256:...` fingerprint. Verification is offline and checks exact payload,
root/key signatures, scope, validity, revocation and manifest expiry. Trust
refresh rejects lower epochs, same-epoch rewrites and changed instance/root.
It writes a separate public file for review; it does not replace trust silently.

Report the trust epoch/expiry. A signature does not claim current selection,
executable availability or permission, and does not turn declared Git metadata
into verified provenance. Offline trust cannot discover newer revocation until
refreshed. Keep private seeds, credentials and `.state` out of Git. Never delete
an uncertain operation checkpoint to get around a verification failure.

A receipt carrying a newer valid root-signed manifest requires explicit trust
refresh before verification; it never silently replaces the pinned manifest.
Signed deployment exports omit private lease proofs and identify their public
projection plus the original owner receipt digest.
