
### Specs: 
It's an API for creating a sequence. Sequence consists of multiple steps, where the user defines how the first and 
followup emails will look like and what are the waiting days in-between them.

Sequence has: name (string), openTrackingEnabled (bool), clickTrackingEnabled (bool). 
Each sequence step has email subject and content.

API exposes these endpoints:

- Create a sequence with steps
- Update a sequence step (new subject or content)
- Delete a sequence step
- Update sequence open or click tracking


---
**Run API**: `go run .` (SQLite db will be auto created & schema will be migrated)  
 Routes are defined inside `api/router.go`

**Run Tests**: `go test -v ./... ` (SQLite test db will be auto created & schema will be migrated)




