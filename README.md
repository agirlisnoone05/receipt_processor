
# Receipt Processor

This project implements a web service for processing receipts and calculating points based on specific rules. The service provides two endpoints:

1. `/receipts/process` (POST): Accepts receipt data and returns a unique receipt ID.
2. `/receipts/{id}/points` (GET): Returns the points awarded for a specific receipt.

---

### Prerequisites
Ensure the following are installed on the system:
- Go (version 1.23 used here)
- Docker

---

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/receipt_processor.git
   cd receipt_processor
   ```

2. Build the Docker image:
   ```bash
   docker build -t receipt_processor .
   ```

3. Run the Docker container:
   ```bash
   docker run -p 8080:8080 receipt_processor
   ```

The service will now be running at `http://localhost:8080`.

---

### Usage

#### 1. Process Receipts (POST)
**Endpoint**: `/receipts/process`

**Example Request**:
```bash
curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },
    {
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },
    {
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },
    {
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}'
```

**Example Response**:
```json
{
    "id": "3482536161526637850"
}
```

---

#### 2. Get Points (GET)
**Endpoint**: `/receipts/{id}/points`

**Example Request**:
```bash
curl http://localhost:8080/receipts/3482536161526637850/points
```

**Example Response**:
```json
{
    "points": 23
}
```

---

### Running the Application Locally Without Docker
If you want to run the application locally without Docker:

```bash
go run main.go
```
