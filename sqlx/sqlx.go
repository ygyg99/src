package main
import(
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)
var db *sqlx.DB

type User struct {
	Age  int `db:"age"`
	Name string `db:"name"`
}  

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := db.Exec(query, args...)
	return err
}

// BatchInsertUsers3 使用NamedExec实现批量插入
func BatchInsertUsers3(users []*User) error {
	_, err := db.NamedExec("INSERT INTO user (name, age) VALUES (:name, :age)", users)
	return err
}

func initDB()(err error){
	dsn := "root:123456@tcp(0.0.0.0:3306)/sql_test?charset=utf8mb4&parseTime=True"
	
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	// 参数设置
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}


func main() {
	if err :=initDB();err != nil{
		fmt.Printf("init DB failed, err: %v \n", err)
		return
	}
	fmt.Println("init DB sucess")
	// queryRowDemo()
	// queryMultiRowDemo()
	u1 := User{Name: "七米", Age: 18}
	u2 := User{Name: "q1mi", Age: 28}
	u3 := User{Name: "小王子", Age: 38}


	// 方法2
	users2 := []interface{}{u1, u2, u3}
	err := BatchInsertUsers2(users2)
	if err != nil {
		fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
	}
	// // 方法1
	// users := []*User{&u1, &u2, &u3}
	// err := BatchInsertUsers3(users)
	// if err != nil {
	// 	fmt.Printf("BatchInsertUsers failed, err:%v\n", err)
	// }
	
}
