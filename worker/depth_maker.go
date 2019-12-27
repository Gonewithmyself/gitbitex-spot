package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gitbitex/gitbitex-spot/matching"
	"github.com/gitbitex/gitbitex-spot/pushing"
	"github.com/shopspring/decimal"
	logger "github.com/siddontang/go-log/log"
)

type OrderBookStream struct {
	productId   string
	name        string
	logReader   matching.LogReader
	logCh       chan *logOffset
	orderBook   *orderBook
	snapshotCh  chan interface{}
	lastSaveSeq int64
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

type logOffset struct {
	log    interface{}
	offset int64
}

func NewDepthMaker(productId string, logReader matching.LogReader) *OrderBookStream {
	s := &OrderBookStream{
		productId:  productId,
		name:       "depthmaker_" + productId,
		orderBook:  newOrderBook(productId),
		logCh:      make(chan *logOffset, 100),
		logReader:  logReader,
		snapshotCh: make(chan interface{}, 100),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// try restore snapshot
	snapshot, err := sharedSnapshotStore().getLastFull(productId)
	if err != nil {
		logger.Fatalf("get snapshot error: %v", err)
	}
	if snapshot != nil {
		s.lastSaveSeq = snapshot.Seq
		s.orderBook.Restore(snapshot)
		logger.Infof("%v order book snapshot loaded: %+v", s.productId, snapshot)
	}

	s.logReader.RegisterObserver(s)
	return s
}

func (s *OrderBookStream) Start() {
	logOffset := s.orderBook.logOffset
	if logOffset > 0 {
		logOffset++
	}
	go s.logReader.Run(s.orderBook.logSeq, logOffset)
	go s.runApplier()
	go s.runPublish()
}

func (s *OrderBookStream) Stop() {
	logger.Info("stopping..", s.name)
	s.logReader.Stop()
	s.cancel()
	s.wg.Wait()
	if s.orderBook.seq != s.lastSaveSeq {
		snap := s.orderBook.SnapshotFull()
		if err := sharedSnapshotStore().storeFull(s.productId, snap); err != nil {
			logger.Error("save snapshot: ", err)
		}
	}
	logger.Info("stopped", s.name)
}

func (s *OrderBookStream) OnOpenLog(log *matching.OpenLog, offset int64) {
	s.logCh <- &logOffset{log, offset}
}

func (s *OrderBookStream) OnMatchLog(log *matching.MatchLog, offset int64) {
	s.logCh <- &logOffset{log, offset}
}

func (s *OrderBookStream) OnDoneLog(log *matching.DoneLog, offset int64) {
	s.logCh <- &logOffset{log, offset}
}

func (s *OrderBookStream) runApplier() {
	var lastLevel2Snapshot *OrderBookLevel2Snapshot
	var lastFullSnapshot *OrderBookFullSnapshot

	s.wg.Add(1)
	defer s.wg.Done()
	tick := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-s.ctx.Done():
			return

		case logOffset := <-s.logCh:
			var l2Change *pushing.Level2Change

			switch log := logOffset.log.(type) {
			case *matching.DoneLog:
				order, found := s.orderBook.orders[log.OrderId]
				if !found {
					continue
				}

				newSize := order.Size.Sub(log.RemainingSize)
				if newSize.LessThan(decimal.Zero) {
					logger.Fatal(log)
				}
				l2Change = s.orderBook.saveOrder(logOffset.offset, log.Sequence, log.OrderId, newSize, log.Price,
					log.Side)

			case *matching.OpenLog:
				l2Change = s.orderBook.saveOrder(logOffset.offset, log.Sequence, log.OrderId, log.RemainingSize,
					log.Price, log.Side)

			case *matching.MatchLog:
				order, found := s.orderBook.orders[log.MakerOrderId]
				if !found {
					panic(fmt.Sprintf("should not happen : %+v", log))
				}
				newSize := order.Size.Sub(log.Size)
				l2Change = s.orderBook.saveOrder(logOffset.offset, log.Sequence, log.MakerOrderId, newSize,
					log.Price, log.Side)
			}

			if lastLevel2Snapshot == nil || s.orderBook.seq-lastLevel2Snapshot.Seq > 10 {
				lastLevel2Snapshot = s.orderBook.SnapshotLevel2(1000)
				s.snapshotCh <- lastLevel2Snapshot
			}

			if lastFullSnapshot == nil || s.orderBook.seq-lastFullSnapshot.Seq > 10000 {
				lastFullSnapshot = s.orderBook.SnapshotFull()
				s.snapshotCh <- lastFullSnapshot
			}

			if l2Change != nil {
				s.snapshotCh <- l2Change
			}

		case <-tick.C:
			if lastLevel2Snapshot == nil || s.orderBook.seq > lastLevel2Snapshot.Seq {
				lastLevel2Snapshot = s.orderBook.SnapshotLevel2(1000)
				s.snapshotCh <- lastLevel2Snapshot
			}
		}
	}
}

func (s *OrderBookStream) runPublish() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case snapshot := <-s.snapshotCh:
			switch msg := snapshot.(type) {
			case *OrderBookLevel2Snapshot:
				err := sharedSnapshotStore().publishDepthSnapshot(s.productId, msg)
				if err != nil {
					logger.Error(err)
				}
			case *OrderBookFullSnapshot:
				err := sharedSnapshotStore().storeFull(s.productId, msg)
				if err != nil {
					logger.Error(err)
				} else {
					s.lastSaveSeq = msg.Seq
				}
			case *pushing.Level2Change:
				err := sharedSnapshotStore().publishChange(s.productId, msg)
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
}
