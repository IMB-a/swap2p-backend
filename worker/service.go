package worker

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/pkg/errors"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/tracker"
)

type TradeEventType int

const (
	TradeEventTypeCreate TradeEventType = 1 + iota
	TradeEventTypeAccept
	TradeEventTypeReject
)

type TradeEvent struct {
	Type TradeEventType
	api.Trade
}

type Service struct {
	ac      *abi.ABI
	idName  map[string]*abi.Method
	idEvent map[string]*abi.Event

	TradeChan chan TradeEvent
	ErrChan   chan error

	cAddress ethgo.Address
	jsonCli  *jsonrpc.Client

	t *tracker.Tracker
}

func NewService(cfg *Config) (*Service, error) {
	s := &Service{
		TradeChan: make(chan TradeEvent),
		ErrChan:   make(chan error),
	}

	abiContract, err := abi.NewABI(cfg.AbiJSON)
	if err != nil {
		return nil, errors.Wrap(err, "new abi")
	}

	s.ac = abiContract
	s.idName = map[string]*abi.Method{}
	s.idEvent = map[string]*abi.Event{}
	for _, v := range abiContract.Methods {
		s.idName[hex.EncodeToString(v.ID())] = v
	}
	for _, v := range abiContract.Events {
		s.idEvent[v.ID().String()] = v
	}
	s.cAddress = ethgo.HexToAddress(cfg.ContractAddress)

	client, err := jsonrpc.NewClient(cfg.JSONRPCClient)
	if err != nil {
		return nil, errors.Wrap(err, "json rpc client")
	}

	s.jsonCli = client

	track, err := tracker.NewTracker(client.Eth(), tracker.WithFilter(&tracker.FilterConfig{
		Address: []ethgo.Address{s.cAddress},
		Start:   cfg.BlockFrom,
	}))
	if err != nil {
		return nil, errors.Wrap(err, "tracker")
	}

	s.t = track

	return s, nil
}

func (s *Service) Run(ctx context.Context) {
	go func() {
		for evt := range s.t.EventCh {
			for _, e := range evt.Added {
				for _, t := range e.Topics {
					if v, ok := s.idEvent[t.String()]; ok {
						dataI, err := v.Inputs.Decode(e.Data)
						if err != nil {
							s.ErrChan <- err
							continue
						}
						data := dataI.(map[string]interface{})
						switch v.Name {
						case "EscrowCreated":
							id := data["escrowIndex"].(*big.Int)
							s.TradeChan <- TradeEvent{
								Type: TradeEventTypeCreate,
								Trade: api.Trade{
									Id:       int(id.Int64()),
									XAddress: data["xOwner"].(ethgo.Address).String(),
									XAmount:  data["xAmount"].(*big.Int).String(),
									XAsset:   data["xTokenContractAddr"].(ethgo.Address).String(),
									YAddress: data["yOwner"].(ethgo.Address).String(),
									YAmount:  data["yAmount"].(*big.Int).String(),
									YAsset:   data["yTokenContractAddr"].(ethgo.Address).String(),
								},
							}
						case "EscrowAccepted":
							s.TradeChan <- TradeEvent{
								Type: TradeEventTypeAccept,
							}
						case "EscrowRejected":
							s.TradeChan <- TradeEvent{
								Type: TradeEventTypeReject,
							}
						}
					}
				}
			}
		}
	}()
	s.t.Sync(ctx)
}
