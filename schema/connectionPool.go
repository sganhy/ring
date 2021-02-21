package schema

import (
	"ring/schema/databaseprovider"
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
	poolId           int
	provider         databaseprovider.DatabaseProvider
	pool             []*connection
	syncRoot         sync.Mutex
}

var currentPoolId = 0

const initialMinValue = 1
const initialMaxValue = 2

func newConnectionPool(connectionString string, provider databaseprovider.DatabaseProvider, minConnection uint16, maxConnection uint16) (*connectionPool, error) {
	var newPool = new(connectionPool)
	currentPoolId++
	newPool.connectionString = connectionString
	newPool.provider = provider
	newPool.poolId = currentPoolId

	// add login
	if maxConnection > initialMaxValue {
		newPool.maxConnection = int(maxConnection)
	} else {
		newPool.maxConnection = initialMaxValue
	}
	if minConnection > initialMinValue && minConnection <= maxConnection {
		newPool.minConnection = int(minConnection)
	} else {
		newPool.minConnection = initialMinValue
	}
	newPool.pool = make([]*connection, 0, newPool.maxConnection)
	// is it unitest
	for i := 0; i < newPool.minConnection; i++ {
		connection, err := newConnection(i+1, connectionString, provider.ToString())
		if err != nil {
			return nil, err
		}
		newPool.pool = append(newPool.pool, connection)
	}
	// cursor pointing on last element
	newPool.cursor = newPool.minConnection - 1
	newPool.lastIndex = newPool.maxConnection - 1
	return newPool, nil
}

func (pool *connectionPool) get() *connection {
	pool.syncRoot.Lock()
	if pool.cursor >= 0 {
		var result = pool.pool[pool.cursor]
		pool.cursor--
		pool.syncRoot.Unlock()
		return result
	}
	pool.syncRoot.Unlock()
	var newConn, _ = newConnection(-1, pool.connectionString, pool.provider.ToString())
	//TODO add login
	return newConn
}

func (pool *connectionPool) put(conn *connection) {
	pool.syncRoot.Lock()
	if pool.cursor < pool.lastIndex {
		pool.cursor++
		pool.putRequestCount++
		pool.swapIndex = int(pool.putRequestCount) % (pool.cursor + 1)
		pool.pool[pool.cursor] = pool.pool[pool.swapIndex]
		pool.pool[pool.swapIndex] = conn
		pool.syncRoot.Unlock()
		return
	}
	pool.syncRoot.Unlock()
	conn.destroy()
}
