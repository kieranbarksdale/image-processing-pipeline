# Image Processing Pipeline

## Requirements
### Functional:

Merchants can upload an image via a REST API
The upload endpoint returns immediately without waiting for processing
The system produces three versions of every uploaded image: thumbnail (150x150), medium (800x800), and large (1920x1920)
Merchants can query the status of their upload — pending, processing, completed, failed
Merchants can retrieve the URLs of their processed images once complete
Failed processing jobs must be retried automatically up to 3 times before being marked as permanently failed

### Non-functional:

The upload endpoint must respond in under 200ms regardless of image size
The system must handle 500 concurrent uploads without degradation
Processed images must be stored durably — not on the local filesystem
The system must be observable — you need to know how many jobs are pending, processing, failed at any given time

## Constraints
- Use Go
- Use PostgreSQL for job metadata
- Use Redis with Asynq as your message queue (this is what a real Go job queue looks like)
- Use MinIO as your object storage — it's an open source S3-compatible storage server you can run locally via Docker. The API is identical to AWS S3 so switching to real S3 in production requires changing one config value.
- All infrastructure via Docker Compose

## Deliverable
A running service where you can POST /upload with an image, get back a job ID, poll GET /jobs/{id} and watch the status move from pending → processing → completed, then call GET /jobs/{id}/images and get back three URLs pointing to the processed versions.


1. What does this service do. 
This service is an image processing pupeline where a user can upload an image, and it will instantaniously get a response of the job's id and start processing the images into 3 different sizes. The user can query the status of the job using the job id and see the status of their image being processed. Once it is done, the user can get three urls of the processed images in different sizes.
2. What are the components 
- A database to hold all of the metadata about the images and their location in the storage container
- A message queue to hold all of the jobs and their order 
- An API gateway to handle the requests 
- A storage container similar to an S3 bucket to store the processed images
- 
3. What breaks at scale?
- Adding lots of images at once can overwhelm the system
- If a message queue fails, jobs can be lost
- If the database goes down, all job metadata is lost
- If the storage container goes down, all processed images are lost
- IF a user submits to many images at once, the system can become overwhelmed 



## Database Schema

### jobs
| Column     | Type    | Description                    |
|------------|---------|--------------------------------|
| id         | UUID    | Primary key                    |
| status     | ENUM    | pending, processing, completed, failed |
| retries    | INT     | Number of processing attempts  |
| error      | TEXT    | Error message if failed, null otherwise |
| created_at | TIMESTAMP | When the job was created     |

### images
| Column     | Type    | Description                    |
|------------|---------|--------------------------------|
| id         | UUID    | Primary key                    |
| job_id     | UUID    | Foreign key → jobs.id          |
| url        | TEXT    | Location of processed image in object storage |
| size       | ENUM    | thumbnail, medium, large       |
| created_at | TIMESTAMP | When the record was created  |