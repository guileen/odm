package dynamo

import (
	"fmt"
	"testing"

	"git.devops.com/go/odm"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func resetDB(t *testing.T) {
	db, err := odm.Open("dynamo", dbpath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db.DropTable("book")
	table := db.Table(&Book{})
	// touch the table.
	err = table.GetItem("A", "B", nil, nil)
	assert.NoError(t, err)
}

func Test_ConnectString(t *testing.T) {
	// s := "AccessKey=[You Access Key];SecretKey=[You Secret Key]+;Region=cn-northwest-1"
	s := "AccessKey=123;SecretKey=456;Token=789;Region=localhost;Endpoint=http://127.0.0.1:8000"
	cfg, err := ParseConnectString(s)
	assert.NoError(t, err)
	assert.Equal(t, &aws.Config{
		Credentials: credentials.NewStaticCredentials("123", "456", "789"),
		Endpoint:    aws.String("http://127.0.0.1:8000"),
		Region:      aws.String("localhost"),
	}, cfg)
}

type Account struct {
	Id      int   `odm:"PK" json:"id"`
	Balance int64 `json:"balance"`
}

type Bag struct {
	Uid       int    `odm:"PK" json:"uid"`
	ProductId string `odm:"SK" json:"product_id"`
	Count     int    `json:"count"`
}

type Product struct {
	Id    string `odm:"PK" json:"id"`
	Price int    `json:"price"`
}

type Order struct {
	Uid      int            `odm:"PK" json:"uid"`
	Tid      int            `odm:"SK" json:"tid"`
	Cart     map[string]int `json:"product_id"`
	TotalFee int
	Status   int8 `json:"status"`
}

const (
	StatusWaitPay = 0
	StatusPayed   = 1
	StatusDone    = 2
	StatusCancel  = -1
)

func ExampleDB_TransactWriteItems() {
	db, _ := odm.Open("dynamo", dbpath)
	accounts, _ := db.ResetTable(&Account{})
	bags, _ := db.ResetTable(&Bag{})
	products, _ := db.ResetTable(&Product{})
	orders, _ := db.ResetTable(&Order{})
	uid := 10
	tid := 1234
	// 用户账户余额1000
	accounts.PutItem(&Account{
		Id:      uid,
		Balance: 100000,
	}, nil, nil)
	// 添加商品
	products.PutItem(&Product{
		Id:    "iPhone",
		Price: 6000,
	}, nil, nil)
	products.PutItem(&Product{
		Id:    "Huawei",
		Price: 3000,
	}, nil, nil)
	// 用户选择商品
	cart := map[string]int{
		"iPhone": 1,
		"Huawei": 1,
	}
	fee := 0
	// TODO: 一致性读取商品价格
	// db.TransactGetItems([]odm.TransGet{

	// })
	// caculate fee
	fee = 9000

	// 保存Order
	err := orders.PutItem(&Order{
		Uid:      uid,
		Tid:      tid,
		Cart:     cart,
		Status:   StatusWaitPay,
		TotalFee: fee,
	}, &odm.WriteOption{
		// 仅当数据不存在时插入
		Condition: "attribute_not_exists(tid)",
	}, nil)
	if err != nil {
		fmt.Printf("Error to save order %v", err)
		return
	}
	// 用户付款时传入订单id，从余额中扣款，改Order状态为已支付
	writeItems := []*odm.TransactWrite{}
	// 扣钱
	writeItems = append(writeItems, &odm.TransactWrite{
		Update: &odm.Update{
			TableName: "account",
			Key: odm.Map{
				"id": uid,
			},
			Expression: "SET balance=balance-:fee",
			WriteOption: &odm.WriteOption{
				Condition: "balance >= :fee",
				ValueParams: odm.Map{
					":fee": fee,
				},
			},
		},
	})
	// 改状态
	writeItems = append(writeItems, &odm.TransactWrite{
		Update: &odm.Update{
			TableName: "order",
			Key: odm.Map{
				"uid": uid,
				"tid": tid,
			},
			Expression: "SET #status=:status",
			WriteOption: &odm.WriteOption{
				Condition: "#status=:preStatus",
				NameParams: map[string]string{
					"#status": "status",
				},
				ValueParams: odm.Map{
					":preStatus": StatusWaitPay,
					":status":    StatusPayed,
				},
			},
		},
	})
	for pid, count := range cart {
		writeItems = append(writeItems, &odm.TransactWrite{
			Update: &odm.Update{
				TableName: "bag",
				Key: odm.Map{
					"uid":        uid,
					"product_id": pid,
				},
				// Expression: "SET #count=#count+:count",
				Expression: "ADD #count :count",
				WriteOption: &odm.WriteOption{
					NameParams: map[string]string{
						"#count": "count",
					},
					ValueParams: odm.Map{
						":count": count,
					},
				},
			},
		})
	}
	err = db.TransactWriteItems(writeItems)
	if err != nil {
		fmt.Printf("Fail to execute transaction %v", err)
		return
	}
	bagItems := []Bag{}
	bags.Query(&odm.QueryOption{
		KeyFilter: "uid=:uid",
		ValueParams: odm.Map{
			":uid": uid,
		},
		Limit: 10,
	}, nil, &bagItems)
	fmt.Println(bagItems)
	// Output:
	// [{10 Huawei 1} {10 iPhone 1}]
}
