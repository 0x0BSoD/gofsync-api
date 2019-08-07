package DB

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

//func CheckOvr(scId int, match string, ctx *user.GlobalCTX) int {
//
//	var id int
//	//fmt.Printf("select id from override_values where sc_id=%d and `match`=%s\n", scId, match)
//	stmt, err := ctx.Config.Database.DB.Prepare("select id from override_values where sc_id=? and `match`=?")
//	if err != nil {
//		logger.Warning.Printf("%q, checkSC", err)
//	}
//	defer utils.DeferCloseStmt(stmt)
//
//	err = stmt.QueryRow(scId, match).Scan(&id)
//	if err != nil {
//		return -1
//	}
//	return id
//}

// ======================================================
// GET
// ======================================================

// ======================================================
// INSERT
// ======================================================

// ======================================================
// DELETE
// ======================================================
func DeleteSmartClass(host string, foremanId int, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM smart_classes WHERE host=? and foreman_id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, foremanId)
	if err != nil {
		logger.Warning.Printf("%q, DeleteSmartClass	", err)
	}
}

func DeleteOverride(scId int, foremanId int, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM override_values WHERE sc_id=? AND foreman_id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(scId, foremanId)
	if err != nil {
		logger.Warning.Printf("%q, DeleteSmartClass	", err)
	}
}
