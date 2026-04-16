# Distributed Image Processing Pipeline

An asynchronous image processing service built in Go. Upload an image, get a job ID 
back immediately, and poll for status while a background worker resizes it into three 
versions.

## Architecture

Two binaries communicate through a Redis-backed task queue:

- **API** (`cmd/api`) — validates uploads, stores the original image in MinIO, writes 
  job metadata to Postgres, enqueues a processing task, and returns a job ID. 
  Responds in under 200ms regardless of image size.
- **Worker** (`cmd/worker`) — pulls tasks from the Asynq queue, fetches the original 
  from MinIO, resizes into three formats, stores results back to MinIO, and updates 
  job status in Postgres. Failed jobs are retried up to 3 times before being marked 
  permanently failed.

## Tech Stack

- **Go** — API and worker binaries
- **PostgreSQL** — job metadata and processed image URLs
- **Redis + Asynq** — distributed task queue
- **MinIO** — S3-compatible object storage
- **Docker Compose** — local infrastructure orchestration

## Getting Started

Prerequisites: Docker and Docker Compose.
```bash
docker compose up --build -d
docker compose logs -f api worker
```

## API

### POST /upload

Upload an image for processing. Returns immediately.
```bash
curl -X POST http://localhost:8080/upload \
  -F "image=@photo.jpg"
```
```json
{ "job_id": "c96be8b8-3127-429f-8628-79ca9198d63e", "status": "pending" }
```

### GET /status/{jobId}

Poll job status. Transitions: `pending` → `processing` → `completed` | `failed`
```bash
curl http://localhost:8080/status/c96be8b8-3127-429f-8628-79ca9198d63e
```
```json
{ "status": "completed" }
```

### GET /images/{jobId}

Returns a zip archive of all three processed versions once the job is complete.
```bash
curl http://localhost:8080/images/c96be8b8-3127-429f-8628-79ca9198d63e \
  --output c96be8b8-3127-429f-8628-79ca9198d63e.zip
```

Archive structure:
```
c96be8b8-3127-429f-8628-79ca9198d63e.zip/
├── 3bb84ab7-668e-47eb-ad23-880c8d52e4c4_2026-04-15_20-47-53.489.jpg
├── 3oidsa9f-ss00-fgio-i852-fa8tvi9apaha_2026-04-15_20-47-53.489.jpg
└── 98sbsksj-0ops-by21-kgu7-0fmfpdisnw35_2026-04-15_20-47-53.489.jpg
```

## Database Schema

### jobs
| Column     | Type      | Description                             |
|------------|-----------|-----------------------------------------|
| id         | UUID      | Primary key                             |
| status     | TEXT      | pending, processing, completed, failed  |
| retries    | INT       | Number of processing attempts           |
| img_key    | TEXT      | Object key of original image in MinIO   |
| error      | TEXT      | Error message if failed, null otherwise |
| created_at | TIMESTAMP | When the job was created                |

### images
| Column     | Type      | Description                            |
|------------|-----------|----------------------------------------|
| id         | UUID      | Primary key                            |
| job_id     | UUID      | Foreign key → jobs.id                  |
| img_key    | TEXT      | Object key of processed image in MinIO |
| size       | TEXT      | thumbnail, medium, large               |
| created_at | TIMESTAMP | When the record was created            |

## Image Sizes

| Size      | Width  |
|-----------|--------|
| Thumbnail | 200px  |
| Medium    | 800px  |
| Large     | 1600px |

## What Breaks at Scale

- **Static worker count** — worker concurrency is fixed in the compose file. Under 
  sustained load the queue grows unbounded. Production would need auto-scaling worker 
  instances based on queue depth.
- **Single MinIO node** — no replication. If the container dies, images are 
  inaccessible until recovery. Production would use S3 or a distributed MinIO cluster.
- **Redis durability** — Redis runs without persistence enabled. A restart loses all 
  pending jobs. Production requires AOF persistence or an external queue with 
  guaranteed delivery.
- **No per-merchant rate limiting** — a single merchant uploading 10,000 images can 
  starve other merchants. Queue priority or per-merchant concurrency limits would 
  be needed.
- **Postgres write amplification** — each job produces 5+ writes across two tables. 
  At high concurrency this becomes a bottleneck. Read replicas and connection pooling 
  would be required at production scale.