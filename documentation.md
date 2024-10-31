# Distribution Management System Documentation

## Table of Contents
1. [Overview](#overview)
2. [System Architecture](#system-architecture)
3. [Data Models](#data-models)
4. [API Usage Examples](#api-usage-examples)
5. [Deployment](#deployment)

## Overview
This is a Go-based distribution management system that handles geographical distribution rights and permissions. The system manages distributors, their hierarchies, and their distribution rights across different geographical areas (countries, provinces, and cities).

## System Architecture

### Core Components

1. **API Layer**
   - Uses Gin framework for HTTP routing
   - Handles distributor management and distribution rights verification
   - Endpoints:
     - POST `/add` - Add new distributor
     - GET `/all` - List all distributors
     - GET `/check` - Check distribution permissions

2. **Service Layer**
   - Business logic implementation in `distributorService.go`
   - Distributor management
   - Distribution rights validation
   - Area code verification

3. **Repository Layer**
   - In-memory data storage using `database/distributor`
   - Distributor data management
   - Location/area data management via CSV

4. **Configuration**
   - JSON-based configuration using Viper
   - Logging setup with Logrus
   - Application paths and settings

## Data Models

### Distributor
```go
type Distributor struct {
    Id                string   `json:"id,omitempty"`
    IncludeCode       []string `json:"includeCode,omitempty"`
    ExcludeCode       []string `json:"excludeCode,omitempty"`
}
```
## API Usage Examples

### 1. Add/Update Distributor
Endpoint to add new distributors or update existing distributor relationships.

**Endpoint:** `POST /add`

#### Use Case 1: Add New Distributor with Area Codes
**Request:**
```bash
curl --location 'localhost:8090/add' \
--header 'Content-Type: application/json' \
--data '{
    "id": "101",
    "includeCode": ["MH-IN","HP-IN"],
    "excludeCode": ["RJ-IN"]
}'
```

**Response:**
```json
{
    "message": "Distributor added successfully"
}
```

#### Use Case 2: Assign Sub-Distributor
**Request:**
```bash
curl --location 'localhost:8090/add' \
--header 'Content-Type: application/json' \
--data '{
    "id": "101",
    "subDistributor": "102"
}'
```

**Response:**
```json
{
    "message": "Distributor added successfully"
}
```

### 2. List All Distributors
Retrieves a list of all registered distributors and their distribution rights.

**Endpoint:** `GET /all`

**Request:**
```bash
curl --location 'localhost:8090/all'
```

**Response:**
```json
{
    "listOfDistributor": [
        {
            "id": "101",
            "includeCode": ["MH-IN", "HP-IN"],
            "excludeCode": ["RJ-IN"]
        },
        {
            "id": "102",
            "includeCode": ["UNYRK-RJ-IN"],
            "excludeCode": ["UDAIR-RJ-IN"],
            "headDistributor": "101"
        }
    ]
}
```

### 3. Check Distribution Permission
Verifies if a distributor has permission for a specific geographical area.

**Endpoint:** `GET /check`

**Request:**
```bash
curl --location --request GET 'localhost:8090/check' \
--header 'Content-Type: application/json' \
--data '{
    "id": "102",
    "area": "UNYRK-RJ-IN"
}'
```

**Response:**
```json
{
    "distributionPermission": false
}
```

**Notes:**
- The `/add` endpoint supports both creating distributors and managing hierarchies
- Area codes can be in two formats:
  - Province level: `{PROVINCE_CODE}-{COUNTRY_CODE}` (e.g., "MH-IN")
  - City level: `{CITY_CODE}-{PROVINCE_CODE}-{COUNTRY_CODE}` (e.g., "UNYRK-RJ-IN")
- The `/all` endpoint returns a stringified JSON response
- The `/check` endpoint verifies distribution rights considering:
  - Directly included/excluded areas
  - Inherited permissions from head distributors
  - Hierarchical relationships between distributors



## Deployment

### Docker Deployment

The application is distributed as a Docker image. Deploy using the following steps:

1. **Extract and Load Docker Image**
   ```bash
   unzip distribution-system.zip
   docker load < distribution-system.tar
   ```

2. **Run Container**
   ```bash
   sudo docker run -d \
     --name cinemacontainer \
     --net=host \
     -v /var/log/cinema:/app/logs \
     cinema
   ```
   **Note:** The `-v /var/log/cinema:/app/logs` option mounts the local directory `/var/log/cinema` to container's `/app/logs` for accessing log files on the host machine.
