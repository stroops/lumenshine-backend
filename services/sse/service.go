package main

import (
	"database/sql"
	"strings"

	m "github.com/Soneso/lumenshine-backend/db/horizon/models"
	"github.com/Soneso/lumenshine-backend/db/querying"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	context "golang.org/x/net/context"
)

func (s *server) ListenFor(ctx context.Context, r *pb.SSEListenForRequest) (*pb.Empty, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	c := &m.SseConfig{
		SourceReceiver: r.SourceReciver,
		StellarAccount: r.StellarAccount,
		OperationTypes: r.OpTypes,
		WithResume:     r.WithResume,
		ReturnData:     r.ReturnData,
	}

	err := c.Insert(s.Env.DBH, boil.Infer())
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) RemoveListening(ctx context.Context, r *pb.SSERemoveListeningRequest) (*pb.Empty, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	_, err := m.SseData(
		qm.Where(m.SseDatumColumns.SourceReceiver+"=? and "+m.SseDatumColumns.StellarAccount+"=?", r.SourceReciver, r.StellarAccount),
	).DeleteAll(s.Env.DBH)

	if err != nil {
		return nil, err
	}

	_, err = m.SseConfigs(
		qm.Where(m.SseConfigColumns.SourceReceiver+"=? and "+m.SseConfigColumns.StellarAccount+"=?", r.SourceReciver, r.StellarAccount),
	).DeleteAll(s.Env.DBH)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) ClearSourceRecivers(ctx context.Context, r *pb.SSEClearSourceReciversRequest) (*pb.Empty, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	_, err := m.SseConfigs(
		qm.Where(m.SseConfigColumns.SourceReceiver+"=?", r.SourceReciver),
	).DeleteAll(s.Env.DBH)

	_, err = m.SseData(
		qm.Where(m.SseDatumColumns.SourceReceiver+"=?", r.SourceReciver),
	).DeleteAll(s.Env.DBH)

	return &pb.Empty{}, err
}

func (s *server) GetData(ctx context.Context, r *pb.SSEGetDataRequest) (*pb.SSEGetDataResponse, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	sqlStr := querying.GetSQLKeyString(`with cte as (select id from sse_data where @status=$1 and @source_receiver=$2 limit $3) update sse_data s set @status=$4 from cte where s.id=cte.id returning *`,
		map[string]string{
			"@sse_data":        m.TableNames.SseData,
			"@status":          m.SseDatumColumns.Status,
			"@source_receiver": m.SseDatumColumns.SourceReceiver,
		})

	var d []m.SseDatum
	err := queries.Raw(sqlStr, m.SseDataStatusNew, r.SourceReciver, r.Count, m.SseDataStatusSelected).Bind(nil, s.Env.DBH, &d)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	var ret pb.SSEGetDataResponse
	ret.Data = make([]*pb.SSEGetData, len(d))
	var str = ""
	for i := 0; i < len(d); i++ {
		if d[i].OperationData.IsZero() {
			str = "{}"
		} else {
			str = string(d[i].OperationData.JSON)
		}

		ret.Data[i] = &pb.SSEGetData{
			SseConfigId:    int64(d[i].SseConfigID),
			SourceReceiver: d[i].SourceReceiver,
			StellarAccount: d[i].StellarAccount,
			OperationType:  int64(d[i].OperationType),
			TransactionId:  d[i].TransactionID,
			OperationId:    d[i].OperationID,
			LedgerId:       d[i].LedgerID,
			OperationData:  str,
		}

	}

	//need to delete the data also
	//_, err = m.SseData(qm.Where(m.SseDatumColumns.Status+"=?", m.SseDataStatusSelected)).DeleteAll(s.Env.DBH)
	return &ret, err
}
