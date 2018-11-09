package main

import (
	"database/sql"

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
	}

	return &pb.Empty{}, c.Insert(s.Env.DBH, boil.Infer())
}

func (s *server) RemoveListening(ctx context.Context, r *pb.SSERemoveListeningRequest) (*pb.Empty, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	_, err := m.SseConfigs(
		qm.Where(m.SseConfigColumns.SourceReceiver+"=? and "+m.SseConfigColumns.StellarAccount+"=?", r.SourceReciver, r.StellarAccount),
	).DeleteAll(s.Env.DBH)

	return &pb.Empty{}, err
}

func (s *server) GetData(ctx context.Context, r *pb.SSEGetDataRequest) (*pb.SSEGetDataResponse, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	sqlStr := querying.GetSQLKeyString(`update @sse_data set @status=$1 where @status=$2 and @source_receiver=$3 limit $4 returning *`,
		map[string]string{
			"@sse_data":        m.TableNames.SseData,
			"@status":          m.SseDatumColumns.Status,
			"@source_receiver": m.SseDatumColumns.SourceReceiver,
		})

	var d []m.SseDatum
	err := queries.Raw(sqlStr, m.SseDataStatusSelected, m.SseDataStatusNew, r.SourceReciver, r.Count).Bind(nil, s.Env.DBH, &d)
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
	_, err = m.SseData(qm.Where(m.SseDatumColumns.Status+"=?", m.SseDataStatusSelected)).DeleteAll(s.Env.DBH)
	return &ret, err
}
