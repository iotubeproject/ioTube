# Runbook: missing / deleted source transaction

## When you see this alert

A witness received a "stale height" from the relayer but the source chain's
node/API has pruned that block, so the witness cannot fetch the transfer to
sign it. The transfer is stuck `new` on the relayer; the witness re-tries and
re-alerts every tick until it is resolved.

The alert arrives as a Lark card carrying the **cashier** and **block height**,
with a single **Mute** button.

You have two choices:

1. **Resolve it** — sign & submit manually (below), then Mute.
2. **Abandon it** — just click Mute (only if the transfer should never settle).

## Step 1 — Find the missing transaction(s)

You have the cashier address and block height from the alert.

- Relayer DB:

  ```sql
  SELECT token, tidx, sender, recipient, amount, payload, sourceTxHash
  FROM transfers
  WHERE cashier = '<cashier>' AND blockHeight = <height> AND status = 'new';
  ```

  (Or use the relayer `List` / `Lookup` gRPC.)
- Cross-check on the block explorer: open the cashier contract at that height
  and confirm the Receipt event's token / index / sender / recipient / amount.

> Note: `token` in the relayer row is the **co-token on the destination chain**
> — that is the `-cotoken` value below.

## Step 2 — Sign (offline, no API needed)

```sh
sign-witness -config <witness-config.yaml> -secret <secret.yaml> -cashier <cashier-id> \
  -cotoken <coToken> -index <tidx> -sender <sender> -recipient <recipient> -amount <amount> \
  [-payload <hex>] [-to-solana]
```

Copy the printed `Signature` and `Transfer ID`.

## Step 3 — Submit to the relayer

```sh
submit-witness -config <relayer-config.yaml> -secret <secret.yaml> \
  -validator <validator> -cashier <cashier> -token <coToken> -index <tidx> \
  -sender <sender> -recipient <recipient> -amount <amount> [-payload <hex>] \
  -signatures <sig-from-step-2> -dry-run    # inspect first, then re-run without -dry-run
```

Once enough witnesses' signatures land, the relayer settles the transfer and it
stops appearing in `StaleHeights`.

## Step 4 — Mute

Click **Mute** on the Lark card so this witness stops retrying/alerting on that
height. Mute is in-memory: if the witness restarts before the transfer settles,
the alert may reappear once — just Mute again.
