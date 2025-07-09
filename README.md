# CLI App for Saving in USDT

Build a simple CLI app in Go that allows a user to save an amount in USDT (a stablecoin pegged to USD). The user is protected from local currency inflation.

## Project Structure

- `cmd/` - Entry point(s) for the CLI application.
- `internal/` - Internal application logic (not intended for public use).
- `pkg/` - Reusable packages that could be used by other projects.

## User Flow

1. User starts the CLI app.
2. User enters the amount they want to save in USD (will be converted to USDT 1:1).
3. App asks for:
   + Wallet address to receive USDT.
   + Local currency amount (if needed for conversion display).
4. App fetches or displays conversion rates (for example, local currency to USD if needed).
5. App simulates or logs:
   - Payment processing steps.
   - Blockchain transaction creation.
   - Generates and displays a transaction hash as proof.
6. Return success or failure message to the user.

## API Usage

### Start the API Server

```
go run cmd/api.go
```

### Endpoint: POST /api/quote

**Request Body:**
```json
{
  "amount": 100,
  "wallet": "your_lisk_wallet_address"
}
```

**Response Example:**
```json
{
  "outputAmount": 0.77,
  "message": "Quote retrieved successfully"
}
```

- `amount`: Amount in KES (minimum 50)
- `wallet`: LISK network wallet address

Returns the approximate USDT amount you will receive for the given KES amount on the LISK network.


