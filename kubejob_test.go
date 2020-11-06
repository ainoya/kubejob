package kubejob_test

import (
	"context"
	"testing"

	"github.com/goccy/kubejob"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

var (
	cfg *rest.Config
)

func init() {
	c, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	cfg = c
}

func Test_Run(t *testing.T) {
	job, err := kubejob.NewJobBuilder(cfg, "default").BuildWithJob(&batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "test",
							Image:   "golang:1.15",
							Command: []string{"echo", "hello"},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to build job: %+v", err)
	}
	if err := job.Run(context.Background()); err != nil {
		t.Fatalf("failed to run: %+v", err)
	}
}

func Test_RunnerWithExecutionHandler(t *testing.T) {
	job, err := kubejob.NewJobBuilder(cfg, "default").BuildWithJob(&batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "test",
							Image:   "golang:1.15",
							Command: []string{"echo", "foo"},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to build job: %+v", err)
	}
	if err := job.RunWithExecutionHandler(context.Background(), func(executors []*kubejob.JobExecutor) error {
		for _, exec := range executors {
			exec.Exec()
		}
		return nil
	}); err != nil {
		t.Fatalf("failed to run: %+v", err)
	}
}
