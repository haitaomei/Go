package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"gopkg.in/mgo.v2"
)

func main() {
	os.Setenv("MONGO_CONNECTSTR", "mongodb://mongo-test:27017/")
	os.Setenv("MONGO_SSLCAFILE_LOCATION", "~/mongo.pem")

	db, err := setupMongoKK()
	if err == nil {
		//db.C("Person").DropCollection()
		db.Logout()
	}
}

func setupMongoKK() (*mgo.Database, error) {

	var mongoConnectStr string
	ok := false
	if mongoConnectStr, ok = os.LookupEnv("MONGO_CONNECTSTR"); !ok {
		return nil, fmt.Errorf("MongoDB connection string not specified. Please set MONGO_CONNECTSTR")
	}

	dialInfo, err := mgo.ParseURL(mongoConnectStr)

	/* **************************************************
	 * read env and test if the sslCAFile exists
	 * If existing, connect with SSL
	 * **************************************************/
	var CALocaton string
	if CALocaton, ok = os.LookupEnv("MONGO_SSLCAFILE_LOCATION"); ok {
		if _, err := os.Stat(CALocaton); err == nil {
			roots := x509.NewCertPool()
			if ca, err := ioutil.ReadFile(CALocaton); err == nil {
				roots.AppendCertsFromPEM(ca)
			}

			tlsConfig := &tls.Config{}
			tlsConfig.RootCAs = roots

			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
				return conn, err
			}
		}
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	db := session.DB("oddk")

	return db, nil
}
