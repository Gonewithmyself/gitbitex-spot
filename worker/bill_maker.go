package worker

import (
	"strconv"
	"strings"
	"time"

	"github.com/gitbitex/gitbitex-spot/matching"
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/gitbitex/gitbitex-spot/models/mysql"
	"github.com/gitbitex/gitbitex-spot/service"
	"github.com/shopspring/decimal"
	"github.com/siddontang/go-log/log"
)

type BillMaker struct {
	billCh    chan *models.OffsetBill
	logReader matching.LogReader
	logOffset int64
	logSeq    int64
	group     string
}

type billItem struct {
	bill    models.Bill
	offset  int64
	seq     int64
	product string
}

const (
	db_account = "db_account"
)

func NewBillMaker(logReader matching.LogReader) *BillMaker {
	t := &BillMaker{
		billCh:    make(chan *models.OffsetBill, 1000),
		logReader: logReader,
		group:     "billmaker_" + logReader.GetProductId(),
	}

	lastBill, err := mysql.SharedStore(db_account).GetLastOffset(t.group, 0)
	if err != nil {
		panic(err)
	}

	if lastBill != nil {
		t.logOffset = lastBill.LogOffset
		t.logSeq = lastBill.LogSeq
	}

	t.logReader.RegisterObserver(t)
	return t
}

func (t *BillMaker) Start() {
	if t.logOffset > 0 {
		t.logOffset++
	}
	go t.logReader.Run(t.logSeq, t.logOffset)
	go t.flusher()
}

func getcurr(pair string) (bc, sc string) {
	coins := strings.Split(pair, "-")
	return coins[0], coins[1]
}

func buy(user int64, first, last string, funds, size decimal.Decimal, notes string) []*models.OffsetBill {
	cost := &models.OffsetBill{
		Bill: models.Bill{
			UserId:    user,
			Currency:  last,
			Hold:      funds.Neg(),
			Available: decimal.Zero,
			Type:      models.BillTypeTrade,
			Notes:     notes,
		},
	}

	got := &models.OffsetBill{
		Bill: models.Bill{
			UserId:    user,
			Currency:  first,
			Hold:      decimal.Zero,
			Available: size,
			Type:      models.BillTypeTrade,
			Notes:     notes,
		},
	}

	return []*models.OffsetBill{cost, got}
}

func sell(user int64, first, last string, funds, size decimal.Decimal, notes string) []*models.OffsetBill {
	cost := &models.OffsetBill{
		Bill: models.Bill{
			UserId:    user,
			Currency:  first,
			Hold:      funds.Neg(),
			Available: decimal.Zero,
			Type:      models.BillTypeTrade,
			Notes:     notes,
		},
	}

	got := &models.OffsetBill{
		Bill: models.Bill{
			UserId:    user,
			Currency:  last,
			Hold:      decimal.Zero,
			Available: size,
			Type:      models.BillTypeTrade,
			Notes:     notes,
		},
	}

	return []*models.OffsetBill{cost, got}
}

func (t *BillMaker) OnMatchLog(log *matching.MatchLog, offset int64) {
	funds := log.Size.Mul(log.Price)
	first, last := getcurr(log.ProductId)
	tradeinfo := "-" + strconv.Itoa(int(log.TradeId))
	var bills []*models.OffsetBill
	if log.Side == models.SideBuy {
		bills = buy(log.Taker, first, last, funds, log.Size, tradeinfo)
		bills = append(bills, sell(log.Maker, first, last, funds, log.Size, tradeinfo)...)
	} else {
		buy(log.Maker, first, last, funds, log.Size, tradeinfo)
		bills = append(bills, sell(log.Taker, first, last, funds, log.Size, tradeinfo)...)
	}

	for i := range bills {
		bills[i].LogOffset = offset
		bills[i].LogSeq = log.Sequence
		t.billCh <- bills[i]
	}
}

func (t *BillMaker) OnOpenLog(log *matching.OpenLog, offset int64) {

}

// todo save offset
func (t *BillMaker) OnDoneLog(log *matching.DoneLog, offset int64) {
	if log.Reason != models.DoneReasonCancelled {
		return
	}

	first, last := getcurr(log.ProductId)
	// market order
	curr := first
	funds := log.RemainingSize
	if log.Side == models.SideBuy {
		curr = last
	}

	bill := &models.OffsetBill{
		Bill: models.Bill{
			Currency:  curr,
			UserId:    log.Taker,
			Available: funds,
			Hold:      funds.Neg(),
			Type:      models.BillTypeTrade,
		},
	}
	bill.LogSeq = log.Sequence
	bill.LogOffset = offset
	t.billCh <- bill
}

func (t *BillMaker) flusher() {
	var bills []*models.OffsetBill

	for {
		select {
		case bill := <-t.billCh:
			bills = append(bills, bill)

			if len(t.billCh) > 0 && len(bills) < 1000 {
				continue
			}

			for {
				err := service.AddOffsetBills(t.group, bills)
				if err != nil {
					log.Error(err)
					time.Sleep(time.Second)
					continue
				}

				// commit

				// TODO redis publish
				bills = nil
				break
			}
		}
	}
}
