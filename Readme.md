# Basic

Start Project with:

```bash
go run ./cmd/api
```

Access Database:

```bash
docker exec -it event-postgres psql -U event_user -d event_platform
```

# Setup Paypal

Create a sandbox environment via paypal.
After create a webhook with all permissions, with the url form ngrok for the webhook.

# Integration Tests

This project contains integration tests for the main application flows.

## Run All Integration Tests

Run all integration tests:

```bash
go test -p 1 ./tests/integration/... -v
```

`-p 1` is used to run tests sequentially and avoid database conflicts between tests.

---

# Test Structure

```
tests/
├── integration/
│   ├── payment/
│   │   └── paypal/
│   ├── purchase/
│   └── ticket/
│       ├── lifecycle/
│       ├── scan/
│       └── verification/
└── helpers/
```

---

# Payment Tests

Payment tests verify the complete payment lifecycle.

## Run Payment Tests

```bash
go test -p 1 ./tests/integration/payment/... -v
```

## PayPal Tests

```bash
go test -p 1 ./tests/integration/payment/paypal/... -v
```

Covered functionality:

* PayPal client creation
* Order creation
* Checkout flow
* Webhook verification
* Payment capture
* Capture completion handling
* Payment confirmation
* Payment event storage
* Webhook idempotency

---

# Purchase Tests

Run purchase integration tests:

```bash
go test -p 1 ./tests/integration/purchase/... -v
```

Covered functionality:

* Creating purchases
* Purchase status changes
* Payment assignment
* Purchase confirmation flow

---

# Ticket Tests

## Ticket Lifecycle

```bash
go test -p 1 ./tests/integration/ticket/lifecycle/... -v
```

Tests:

* Ticket generation after successful payment
* Ticket state transitions
* Ticket ownership

---

## Ticket Scanning

```bash
go test -p 1 ./tests/integration/ticket/scan/... -v
```

Tests:

* Ticket scan creation
* Scan validation
* Duplicate scan prevention

---

## Ticket Verification

```bash
go test -p 1 ./tests/integration/ticket/verification/... -v
```

Tests:

* Verification windows
* Valid ticket verification
* Expired verification handling

---

# Full Payment Flow

The expected payment lifecycle:

```
Create Purchase
        |
        v
Create PayPal Checkout
        |
        v
Customer approves payment
        |
        v
CHECKOUT.ORDER.APPROVED webhook
        |
        v
Capture payment
        |
        v
PAYMENT.CAPTURE.COMPLETED webhook
        |
        v
Confirm Payment
        |
        v
Purchase marked PAID
        |
        v
Tickets generated
```

---

# Required Environment Variables

Integration tests require the application environment configuration.

Example:

```env
PAYPAL_CLIENT_ID=
PAYPAL_CLIENT_SECRET=
PAYPAL_BASE_URL=
PAYPAL_RETURN_URL=
PAYPAL_CANCEL_URL=
PAYPAL_WEBHOOK_ID=

DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
```

---

# Database Requirements

Tests require a running PostgreSQL database.

Before running tests:

1. Start PostgreSQL
2. Configure environment variables
3. Ensure migrations are applied

The application uses GORM migrations during startup.

---

# Running a Single Test

Run one test function:

```bash
go test -p 1 ./tests/integration/payment/paypal -run TestName -v
```

Example:

```bash
go test -p 1 ./tests/integration/payment/paypal -run TestCapturePayment -v
```

---

# Test Goals

The integration suite verifies:

* Authentication flow
* Party and purchase creation
* PayPal payment processing
* Webhook security verification
* Payment event idempotency
* Purchase state updates
* Ticket generation
* Ticket validation
* Database persistence
