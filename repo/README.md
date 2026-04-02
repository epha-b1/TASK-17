# ParkOps

Start the platform:

`docker compose up --build`

Run all tests:

`run_tests.sh`

Access the login page:

`http://localhost:8080/login`

## Default Admin Credentials

**Login URL**: http://localhost:8080/login

**Admin Account**:
- **Username**: `admin`
- **Password**: `AdminPass1234`

**API Testing with curl**:
```bash
# Login and save session cookie
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "AdminPass1234"}' \
  -c cookies.txt

# Test authenticated endpoint
curl -X GET http://localhost:8080/api/me \
  -b cookies.txt

# Test admin endpoint
curl -X GET http://localhost:8080/api/admin/users \
  -b cookies.txt
```

**Additional Test Users** (create via API after admin login):
- Fleet Manager: `fleet1` / `FleetPass1234`
- Dispatch Operator: `dispatch1` / `DispatchPass1234`
- Auditor: `auditor1` / `AuditorPass1234`
