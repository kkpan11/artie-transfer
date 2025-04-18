package consumer

import (
	"testing"

	"github.com/artie-labs/transfer/lib/config"
	"github.com/artie-labs/transfer/lib/config/constants"
	"github.com/artie-labs/transfer/lib/destination"
	"github.com/artie-labs/transfer/lib/kafkalib"
	"github.com/artie-labs/transfer/lib/mocks"
	"github.com/artie-labs/transfer/models"
	"github.com/stretchr/testify/suite"
)

func SetKafkaConsumer(_topicToConsumer map[string]kafkalib.Consumer) {
	topicToConsumer = &TopicToConsumer{
		topicToConsumer: _topicToConsumer,
	}
}

type FlushTestSuite struct {
	suite.Suite
	fakeConsumer *mocks.FakeConsumer
	cfg          config.Config
	db           *models.DatabaseData
	fakeBaseline *mocks.FakeBaseline
	baseline     destination.Baseline
}

func (f *FlushTestSuite) SetupTest() {
	tc := &kafkalib.TopicConfig{
		Database:     "db",
		Schema:       "schema",
		Topic:        "topic",
		CDCFormat:    constants.DBZPostgresFormat,
		CDCKeyFormat: kafkalib.JSONKeyFmt,
	}

	tc.Load()

	f.cfg = config.Config{
		Mode: config.Replication,
		Kafka: &config.Kafka{
			BootstrapServer: "foo",
			GroupID:         "bar",
			Username:        "user",
			Password:        "abc",
			TopicConfigs:    []*kafkalib.TopicConfig{tc},
		},
		Queue:                constants.Kafka,
		Output:               "snowflake",
		BufferRows:           500,
		FlushIntervalSeconds: 60,
		FlushSizeKb:          500,
	}

	f.fakeBaseline = &mocks.FakeBaseline{}
	f.baseline = f.fakeBaseline
	f.db = models.NewMemoryDB()
	f.fakeConsumer = &mocks.FakeConsumer{}
	SetKafkaConsumer(map[string]kafkalib.Consumer{"foo": f.fakeConsumer})
}

func TestFlushTestSuite(t *testing.T) {
	suite.Run(t, new(FlushTestSuite))
}
