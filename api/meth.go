package api

func (q *QTradeClosed) Bool() *bool {
	if q == nil {
		return nil
	}
	b := bool(*q)
	return &b
}
