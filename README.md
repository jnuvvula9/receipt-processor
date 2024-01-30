# Receipt Processor Service

## Description
This Go-based service provides an API for processing receipts. It calculates points based on specific rules and criteria, running within a Docker container.

## Prerequisites
- Docker
- Go 1.21.6

## Installation & Setup
1. Clone the repository: `git clone https://github.com/jnuvvula9/receipt-processor.git`
2. Navigate to the project directory: `cd receipt-processor`

## Running the Application
1. Build the Docker image: `docker build -t receipt-processor-app .`
2. Run the application: `docker run -d -p 8080:63342 receipt-processor-app`
3. The service is now accessible at `http://localhost:8080`.

## Using the Application
- To process a receipt, POST JSON data to `http://localhost:8080/receipts/process`.
- Once a receipt is processed, you can fetch the points gained for that receipt by sending a GET request.
- You will need the unique ID that was returned when you processed the receipt http://localhost:8080/receipts/[receipt_id]/points

