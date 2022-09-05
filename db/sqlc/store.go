package db

import (
	"Edtech_Golang/util"
	"context"
	"database/sql"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type Store interface {
	Querier
	OnlineUserMaps
	ScoreInterface
	SendChildNotification(ctx context.Context, userId int32, notificationType, body string) error
	SendParentsNotification(ctx context.Context, userId int32, notificationType, body string) error
	SearchList(ctx context.Context, arg SearchListParams) ([]SearchListRow, error)
	UploadToS3(fileName string, data []byte) error
	GetBookImageUrl(bookId int32) string
	GetBookPdfUrl(bookId int32) string
}

//SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	db          *sql.DB
	Config      *util.Config
	fApp        *firebase.App
	awsSession  *awsSession.Session
	emailClient *ses.SES
	*Queries
	*OnlineUserMap
}

//NewStore creates a new store
func NewStore(db *sql.DB, config *util.Config) Store {
	awsConf := aws.Config{
		Credentials: awsCredentials.NewStaticCredentials(config.AWS.ID, config.AWS.Secret, ""),
		Region:      &config.AWS.Region,
	}
	sess := awsSession.Must(awsSession.NewSession(&awsConf))
	ouMap := newOnlineUserMap()
	return &SQLStore{
		db:            db,
		Queries:       New(db),
		OnlineUserMap: &ouMap,
		Config:        config,
		fApp:          initializeAppWithServiceAccount(config),
		awsSession:    sess,
		emailClient:   ses.New(sess, &awsConf),
	}
}
func initializeAppWithServiceAccount(config *util.Config) *firebase.App {
	absPath, err := filepath.Abs(config.Firebase.CredentialPath)
	if err != nil || config.Firebase.CredentialPath == "" {
		logrus.Fatal("firebase credential file not found in " + filepath.Join(config.Firebase.CredentialPath))
	}
	opt := option.WithCredentialsFile(absPath)
	fConfig := &firebase.Config{ProjectID: "esikshya-6c829"}
	app, err := firebase.NewApp(context.Background(), fConfig, opt)
	if err != nil {
		logrus.Fatal(err)
	}
	return app
}
