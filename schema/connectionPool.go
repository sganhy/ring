package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"strings"
	"sync"
)

type connectionPool struct {
	connectionString string
	maxConnection    int
	minConnection    int
	cursor           int
	swapIndex        int
	lastIndex        int
	putRequestCount  uint16
	poolId           int32
	provider         databaseprovider.DatabaseProvider
	pool             []*connection
	syncRoot         sync.Mutex
}

const (
	initialMinValue           uint16 = 1
	initialMaxValue           uint16 = 2
	connStringApplicationName string = " application_name"
)

var (
	currentPoolId int32 = 0
)

func (connPool *connectionPool) Init(schemaId int32, connectionString string, provider databaseprovider.DatabaseProvider, minConnection uint16, maxConnection uint16) {
	currentPoolId++
	connPool.connectionString = connectionString //connPool.getConnectionString(schemaId, connectionString)
	connPool.provider = provider
	connPool.poolId = currentPoolId

	// add login
	if maxConnection > initialMaxValue {
		connPool.maxConnection = int(maxConnection)
	} else {
		connPool.maxConnection = int(initialMaxValue)
	}
	if minConnection > initialMinValue && minConnection <= maxConnection {
		connPool.minConnection = int(minConnection)
	} else {
		connPool.minConnection = int(initialMinValue)
	}
	connPool.pool = make([]*connection, 0, connPool.maxConnection)

	for i := 0; i < connPool.minConnection; i++ {
		connection := new(connection)
		err := connection.Init(i+1, connPool.connectionString, provider.String())
		if err != nil {
			panic(err)
		}
		connPool.pool = append(connPool.pool, connection)
	}
	// cursor pointing on last element
	connPool.cursor = connPool.minConnection - 1
	connPool.lastIndex = connPool.maxConnection - 1
}

//******************************
// getters and setters
//******************************
func (connPool *connectionPool) IsDisabled() bool {
	return connPool.poolId == -1
}
func (connPool *connectionPool) GetMaxConnection() uint16 {
	return uint16(connPool.maxConnection)
}
func (connPool *connectionPool) GetMinConnection() uint16 {
	return uint16(connPool.minConnection)
}
func (connPool *connectionPool) GetConnectionString() string {
	return connPool.connectionString
}
func (connPool *connectionPool) GetDatabaseProvider() databaseprovider.DatabaseProvider {
	return connPool.provider
}
func (connPool *connectionPool) setConnectionString(connectionString string) {
	connPool.connectionString = connectionString
}
func (connPool *connectionPool) setDatabaseProvider(provider databaseprovider.DatabaseProvider) {
	connPool.provider = provider
}
func (connPool *connectionPool) setDisabled() {
	connPool.poolId = -1
}

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
func (connPool *connectionPool) get() *connection {
	connPool.syncRoot.Lock()
	if connPool.cursor >= 0 {
		var result = connPool.pool[connPool.cursor]
		connPool.pool[connPool.cursor] = nil
		connPool.cursor--
		connPool.syncRoot.Unlock()
		return result
	}
	connPool.syncRoot.Unlock()

	newConn := new(connection)
	fmt.Println("connectionPool::newConn := new(connection)")
	_ = newConn.Init(-1, connPool.connectionString, connPool.provider.String())
	//TODO add login
	return newConn
}

func (connPool *connectionPool) put(conn *connection) {
	connPool.syncRoot.Lock()
	if connPool.cursor < connPool.lastIndex {
		connPool.cursor++
		connPool.putRequestCount++
		connPool.swapIndex = int(connPool.putRequestCount) % (connPool.cursor + 1)
		connPool.pool[connPool.cursor] = connPool.pool[connPool.swapIndex]
		connPool.pool[connPool.swapIndex] = conn
		connPool.syncRoot.Unlock()
		return
	}
	connPool.syncRoot.Unlock()
	fmt.Println("connectionPool::conn.close()")

	conn.close()
}

func (connPool *connectionPool) getConnectionString(schemaid int32, connectionString string) string {
	if !strings.Contains(strings.ToLower(connectionString), connStringApplicationName) {
		return connectionString + connStringApplicationName + fmt.Sprintf("=Ring(%d)", schemaid)
	}
	return connectionString
}
