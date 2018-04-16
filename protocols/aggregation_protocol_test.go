package protocolsunlynxsmc

import (
	"github.com/dedis/onet"
	"testing"
	"time"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"
	"github.com/henrycg/prio/share"
	"github.com/lca1/unlynx-smc/lib"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"math/big"
	lib "github.com/lca1/unlynx-smc/lib"
)

// the field cardinality must be superior to nbclient*2^b where b is the maximum number of bit a client need to encode its valu

//2 random numbers to test, you can test it with smaller number to see the sum yourself
var secret1 = big.NewInt(int64(55189642165))
var secret2 = big.NewInt(int64(5518495792165))

//the share of them
var secret1Share = share.Share(share.IntModulus, lib.NbServers, secret1)
var secret2Share = share.Share(share.IntModulus, lib.NbServers, secret2)

func TestAggregationProtocol(t *testing.T) {
	local := onet.NewLocalTest(libunlynx.SuiTe)

	// You must register this protocol before creating the servers
	onet.GlobalProtocolRegister("AggregationTest", NewAggregationTest)
	_, _, tree := local.GenTree(lib.NbServers, true)
	defer local.CloseAll()

	p, err := local.CreateProtocol("AggregationTest", tree)
	if err != nil {
		t.Fatal("Couldn't start protocol:", err)
	}

	protocol := p.(*AggregationProtocol)
	start := time.Now()
	protocol.Start()
	timeout := network.WaitRetry * time.Duration(network.MaxRetryConnect*5*2) * time.Millisecond

	//verify results
	expectedResults := big.NewInt(int64(0))
	expectedResults.Add(expectedResults, secret1)
	expectedResults.Add(expectedResults, secret2)
	expectedResults.Mod(expectedResults, field)
	println("EXPECTED", expectedResults.Int64())

	select {
	case Result := <-protocol.Feedback:
		log.Lvl1("time elapsed is ", time.Since(start))
		//JS: get the result and reduce it modulo the field we are working on
		result := libunlynxsmc.Decode(Result, "sum")
		result.Mod(result, field)
		println("RESULT", result.Int64())
		assert.Equal(t, expectedResults, result)
	case <-time.After(timeout):
		t.Fatal("Didn't finish in time")
	}
}

//inject Test data
func NewAggregationTest(tni *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	pi, err := NewAggregationProtocol(tni)
	protocol := pi.(*AggregationProtocol)

	//here assign a share of each secret to the server.
	// Meaning if 2 server, secret1 = [share1,share2] each of them goes to different server
	// (1 and 2 respectively even if order does not matter)
	//You use AFE encoding to encode the shares.
	protocol.Modulus = share.IntModulus
	protocol.Shares = make([][]*big.Int, 0)
	protocol.Shares = append(protocol.Shares, libunlynxsmc.Encode(secret1Share[tni.Index()], "sum"))
	protocol.Shares = append(protocol.Shares, libunlynxsmc.Encode(secret2Share[tni.Index()], "sum"))
	return protocol, err
}