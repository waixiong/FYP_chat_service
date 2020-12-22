package grpcClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"

	mailnotificationv1 "getitqec.com/server/chat/pkg/api/clients/mailnotification/v1"
	"getitqec.com/server/chat/pkg/commons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetMailNotificationClient() (mailnotificationv1.MailNotificationServiceClient, *grpc.ClientConn, error) {
	b, _ := ioutil.ReadFile(CertFilePath)
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		log.Fatalf("fail to dial: %v", errors.New("credentials: failed to append certificates"))
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            cp,
	}
	creds := credentials.NewTLS(config)

	a := commons.ENVVariable("MAIL_NOTIFICATION_SERVER_ADDR")
	if len(a) == 0 {
		a = "0.0.0.0:8120"
	}
	conn, err := grpc.Dial(a, grpc.WithTransportCredentials(creds))
	return mailnotificationv1.NewMailNotificationServiceClient(conn), conn, err
}
