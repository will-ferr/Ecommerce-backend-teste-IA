package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type JobType string

const (
	JobTypeEmailSend      JobType = "email_send"
	JobTypeReportGenerate JobType = "report_generate"
	JobTypeDataCleanup    JobType = "data_cleanup"
)

type Job struct {
	ID          string    `json:"id"`
	Type        JobType   `json:"type"`
	Payload     []byte    `json:"payload"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
	CreatedAt   time.Time `json:"created_at"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type JobQueue interface {
	Enqueue(ctx context.Context, job Job) error
	Dequeue(ctx context.Context) (*Job, error)
	Complete(ctx context.Context, jobID string) error
	Fail(ctx context.Context, jobID string, reason string) error
	GetStats(ctx context.Context) (map[string]int, error)
}

type RedisJobQueue struct {
	client *redis.Client
}

func NewRedisJobQueue(addr, password string) (*RedisJobQueue, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1, // Use DB 1 for job queue
		PoolSize: 10,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisJobQueue{client: rdb}, nil
}

func (q *RedisJobQueue) Enqueue(ctx context.Context, job Job) error {
	job.ID = generateJobID()
	job.CreatedAt = time.Now()
	job.Attempts = 0
	job.MaxAttempts = 3

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Add to queue list
	return q.client.LPush(ctx, "job_queue", jobJSON).Err()
}

func (q *RedisJobQueue) Dequeue(ctx context.Context) (*Job, error) {
	result, err := q.client.BRPop(ctx, "job_queue", 10*time.Second).Result()
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil // No jobs available
	}

	var job Job
	if err := json.Unmarshal([]byte(result), &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (q *RedisJobQueue) Complete(ctx context.Context, jobID string) error {
	// Remove from processing queue
	return q.client.LRem(ctx, "processing_jobs", jobID).Err()
}

func (q *RedisJobQueue) Fail(ctx context.Context, jobID string, reason string) error {
	// Move to dead letter queue for analysis
	deadJob := map[string]interface{}{
		"job_id":    jobID,
		"reason":    reason,
		"failed_at": time.Now(),
	}

	deadJobJSON, _ := json.Marshal(deadJob)
	return q.client.LPush(ctx, "dead_letter_queue", deadJobJSON).Err()
}

func (q *RedisJobQueue) GetStats(ctx context.Context) (map[string]int, error) {
	queueLen, err := q.client.LLen(ctx, "job_queue").Result()
	if err != nil {
		return nil, err
	}

	processingLen, err := q.client.LLen(ctx, "processing_jobs").Result()
	if err != nil {
		return nil, err
	}

	deadLen, err := q.client.LLen(ctx, "dead_letter_queue").Result()
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"queue_length":      int(queueLen),
		"processing_length": int(processingLen),
		"dead_letter_count": int(deadLen),
	}, nil
}

func generateJobID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
