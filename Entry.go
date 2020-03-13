package main

import (
	"github.com/fatih/structs"
	"github.com/golang/protobuf/proto"
	types "github.com/golang/protobuf/ptypes"
	aws "github.com/polyglotDataNerd/zib-Go-Utils/aws"
	c "github.com/polyglotDataNerd/zib-Go-Utils/database"
	p "github.com/polyglotDataNerd/zib-Go-Utils/utils"
	pb "github.com/sweetgreen/sg-gRPC-orderhistory/definiton"
	"sync"
)

func main() {
	p.Warning.Println("EntryPoint")

	collection := pb.HistoryResponse{}
	var userorders []*pb.HistoryMap
	var collectHistory []map[string]interface{}
	var wg sync.WaitGroup
	uuid := "d35e848c-f44e-11e9-abac-8d25e63738f4"

	props := p.Mutator{
		SetterKeyEnv:    "host",
		SetterValueEnv:  "cassandra.us-east-1.amazonaws.com",
		SetterKeyUser:   "user",
		SetterValueUser: aws.SSMParams("/cassandra/mcs/ServiceUserName", 0),
		SetterKeyPW:     "pw",
		SetterValuePW:   aws.SSMParams("/cassandra/mcs/ServicePassword", 0),
	}

	clientConfig := c.CQLProps{
		Mutator: props,
	}

	client := c.CQL{
		CQLProps: clientConfig,
		Wg:       wg,
		SSLPath:  "/Users/gerardbartolome/.mac-ca-roots",
	}

	session := client.CassandraSession()
	query := "SELECT * FROM sg_cass.order_history where gid =" + "'" + uuid + "'"
	resultSet, rerr := client.CassReadOrderHistory(query, session)
	if rerr != nil {
		p.Info.Println(rerr)
	}
	for _, v := range resultSet {
		timeConvert, _ := types.TimestampProto(v.OrderDate)
		userorders = append(userorders, &pb.HistoryMap{
			Gid:        v.Gid,
			OrderDate:  timeConvert,
			OrderID:    v.OrderId,
			Entree:     v.Entree,
			CustomerId: v.CustomerID,
		})
	}
	session.Close()
	collection.Collection = userorders
	for _, o := range collection.Collection {
		collectHistory = append(collectHistory, structs.Map(o))
	}
	data, err := proto.Marshal(&collection)
	if err != nil {
		p.Error.Println("you suck")
	}
	/*byte array wireframe serialized*/
	/*deserialize*/
	newCollection := pb.HistoryResponse{}
	proto.Unmarshal(data, &newCollection)
	p.Info.Println(newCollection.GetCollection())

}