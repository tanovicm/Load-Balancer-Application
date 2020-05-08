echo "registering user1" 
curl -XPOST -d '{"Username": "user1", "Password": "123"}' 'http://localhost:8090/register'

echo "login user1" 
PTOKEN=$(curl -XPOST -d '{"Username": "user1", "Password": "123"}' 'http://localhost:8090/login')

echo "create bank account for user1"
ACCOUNTID=$(curl -XPOST -H "Token: $PTOKEN" -d '{"Name": "pata"}' 'http://localhost:8090/bank_account' | jq ".accountID")

echo "fetch bank account for user1"
curl -XGET -H "Token: $PTOKEN" -d "{\"AccountID\": $ACCOUNTID}" 'http://localhost:8090/bank_account'

echo "create expense for user1"
curl -XPOST -H "Token: $PTOKEN" -d "{\"AccountID\":$ACCOUNTID,\"Name\": \"pokon\", \"Amount\":100}" 'http://localhost:8090/expense'

echo "fetch bank account for user1"
curl -XGET -H "Token: $PTOKEN" -d "{\"AccountID\": $ACCOUNTID}" 'http://localhost:8090/bank_account'

echo "create expense for user1"
EXPENSEID=$(curl -XPOST -H "Token: $PTOKEN" -d "{\"AccountID\":$ACCOUNTID,\"Name\": \"dugi pokon\", \"Amount\":200}" 'http://localhost:8090/expense' | jq ".expenseID") 

echo "fetch expense for user1"
curl -XGET -H "Token: $PTOKEN" -d "{\"ExpenseID\": $EXPENSEID}" 'http://localhost:8090/expense'

echo "fetch bank account for user1"
curl -XGET -H "Token: $PTOKEN" -d "{\"AccountID\": $ACCOUNTID}" 'http://localhost:8090/bank_account'

echo "registering user2" 
curl -XPOST -d '{"Username": "user2", "Password": "123"}' 'http://localhost:8090/register'

echo "login user2" 
STOKEN=$(curl -XPOST -d '{"Username": "user2", "Password": "123"}' 'http://localhost:8090/login')

echo "fetch bank account for user1 by srce"
curl -XGET -H "Token: $STOKEN" -d "{\"AccountID\": $ACCOUNTID}" 'http://localhost:8090/bank_account'

echo "fetch expense for user1 by user2"
curl -XGET -H "Token: $STOKEN" -d "{\"ExpenseID\": $EXPENSEID}" 'http://localhost:8090/expense'

echo "delete expense for user1"
curl -XDELETE -H "Token: $PTOKEN" -d "{\"ExpenseID\": $EXPENSEID}" 'http://localhost:8090/expense'

echo "fetch bank account for user1"
curl -XGET -H "Token: $PTOKEN" -d "{\"AccountID\": $ACCOUNTID}" 'http://localhost:8090/bank_account'

echo "delete bank account for user1"
curl -XDELETE -H "Token: $PTOKEN" -d "{\"AccountID\": $ACCOUNTID}" 'http://localhost:8090/bank_account'

