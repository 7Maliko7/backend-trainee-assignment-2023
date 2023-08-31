package segments

import (
	"context"
	"database/sql"
	"flag"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/config"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/middleware"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/service"
	transport "github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/structs"
	"github.com/7Maliko7/backend-trainee-assignment-2023/pkg/db/driver/postgres"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"math/rand"
	"os"
	"testing"
)

var (
	configPath string
	repository *postgres.Repository
	srvc       service.SegmentsService
	user_id    int64
)

func TestNewService(t *testing.T) {
	flag.StringVar(&configPath, "c", "", "Custom config path")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "backend-trainee-assignment-2023",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	appConfig, err := config.New(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var db *sql.DB
	{
		db, err = sql.Open("postgres", appConfig.RWDB.ConnectionString)
		if err != nil {
			t.Fatal(err)
		}
	}

	repository, err = postgres.New(db)
	if err != nil {
		t.Fatal(err)
	}

	srvc = NewService(repository, logger)
	srvc = middleware.LoggingMiddleware(logger)(srvc)
}

func TestAddSegment(t *testing.T) {
	resp, err := srvc.AddSegment(context.TODO(), transport.AddSegmentRequest{
		Slug: "TEST_SEGMENT",
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("no response from AddSegment")
	}

	if resp.Id == 0 {
		t.Fatal("segment can't be added")
	}

	t.Logf("created segment id = %v", resp.Id)
}

func TestUpdateUserSegmentAdd(t *testing.T) {
	user_id = rand.Int63n(5000)

	resp, err := srvc.UpdateUserSegment(context.TODO(), transport.UpdateUserSegmentRequest{
		UserId: user_id,
		Segments: []transport.SegmentAction{{
			Slug:   "TEST_SEGMENT",
			Action: "add",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("no response from UpdateUserSegment")
	}

	if resp.Status != "success" {
		t.Fatal("non success response")
	}
}

func TestGetSegments(t *testing.T) {
	resp, err := srvc.GetSegments(context.TODO(), transport.GetSegmentsRequest{
		UserId: user_id,
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if resp == nil {
		t.Log("no response from UpdateUserSegment")
		t.Fail()
	}

	if resp.Segments[0] != "TEST_SEGMENT" {
		t.Logf("invalid segment = %v", resp.Segments[0])
		t.Fail()
	}
}

func TestUpdateUserSegmentDelete(t *testing.T) {
	resp, err := srvc.UpdateUserSegment(context.TODO(), transport.UpdateUserSegmentRequest{
		UserId: user_id,
		Segments: []transport.SegmentAction{{
			Slug:   "TEST_SEGMENT",
			Action: "delete",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("no response from UpdateUserSegment")
	}

	if resp.Status != "success" {
		t.Fatal("non success response")
	}
}

func TestDeleteSegment(t *testing.T) {
	resp, err := srvc.DeleteSegment(context.TODO(), transport.DeleteSegmentRequest{
		Slug: "TEST_SEGMENT",
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("no response from DeleteSegment")
	}

	if resp.Status != "success" {
		t.Fatal("non success response")
	}
}

func TestCloseup(t *testing.T) {
	err := repository.Close()
	if err != nil {
		t.Fatal(err)
	}
}
