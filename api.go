package nylechain

import (
	"github.com/dedis/student_19_nylechain/service"
	"go.dedis.ch/cothority/v3"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
)

// Client is a structure to communicate with the template
// service
type Client struct {
	*onet.Client
}

// StoreTree stores the input tree in that ServerIdentity
func (c *Client) StoreTree(si *network.ServerIdentity, tree *onet.Tree) error {
	void := &service.VoidReply{}
	err := c.SendProtobuf(si, &service.StoreTreeArg{Tree: tree}, void)
	if err != nil {
		return err
	}
	return nil

}

// Setup sends a SetupArgs to every server. It returns an error if there was one for any of the servers.
func (c *Client) Setup(servers []*onet.Server, translations map[onet.TreeID][]byte,
	localityTrees map[string][]*onet.Tree) error {
	var serverIDS []*network.ServerIdentity
	for _, server := range servers {
		serverIDS = append(serverIDS, server.ServerIdentity)
	}
	void := &service.VoidReply{}
	sArgs := &service.SetupArgs{
		ServerIDS: serverIDS, Translations: translations,
	}
	for _, si := range serverIDS {
		err := c.SendProtobuf(si, sArgs, void)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenesisTx sends a GenesisArgs to every server. It returns an error if there was one for any of the servers.
func (c *Client) GenesisTx(servers []*onet.Server, id []byte, coinID []byte, rPK kyber.Point) error {
	void := &service.VoidReply{}
	receiverPK, err := rPK.MarshalBinary()
	if err != nil {
		return err
	}
	gArgs := &service.GenesisArgs{
		ID: id, CoinID: coinID, ReceiverPK: receiverPK,
	}
	for _, server := range servers {
		err := c.SendProtobuf(server.ServerIdentity, gArgs, void)
		if err != nil {
			return err
		}
	}
	return nil
}

// TreesBLSCoSi sends a CoSiTrees to the specified Server, and returns a CoSiReplyTrees or an eventual error.
func (c *Client) TreesBLSCoSi(si *network.ServerIdentity, treeIDs []onet.TreeID, message []byte) (*service.CoSiReplyTrees, error) {
	reply := &service.CoSiReplyTrees{}
	err := c.SendProtobuf(si, &service.CoSiTrees{TreeIDs: treeIDs, Message: message}, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// NewClient instantiates a new template.Client
func NewClient() *Client {
	return &Client{Client: onet.NewClient(cothority.Suite, service.ServiceName)}
}
