package server

import "github.com/stockyard-dev/stockyard-homestead/internal/license"

type Limits struct { MaxBookmarks int; MaxNotes int }
var freeLimits = Limits{MaxBookmarks: 20, MaxNotes: 10}
var proLimits = Limits{MaxBookmarks: 0, MaxNotes: 0}

func LimitsFor(info *license.Info) Limits {
	if info != nil && info.IsPro() { return proLimits }
	return freeLimits
}
func LimitReached(limit, current int) bool { return limit > 0 && current >= limit }
