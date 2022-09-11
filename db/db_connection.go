package db

import (
	"database/sql"
	"discord_crypto/dto"
	config "discord_crypto/util"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func AddGenesis(messageDTO dto.MessageDTO, config config.Config) {
	db, err := sql.Open("mysql", config.DB_PATH)

	if err != nil {
		fmt.Println(fmt.Sprintf("DB Error: %s", err.Error()))
		return
	}

	defer db.Close()

	sqlString := fmt.Sprintf("INSERT INTO player_order (player_name, player_id, message, message_id, refer_id, is_genesis, status) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		messageDTO.ByWho, messageDTO.ByWhoID, messageDTO.Message, messageDTO.MsgId, messageDTO.InReply, "1", messageDTO.Status)

	if messageDTO.Status == "edited" {
		sqlString = fmt.Sprintf("INSERT INTO player_order (message_id, player_name, player_id, message, is_genesis, refer_id, status) SELECT '%s', player_name, player_id, '%s', '%s', '%s', '%s' FROM player_order WHERE message_id = %s",
			messageDTO.MsgId, messageDTO.Message, "1", messageDTO.InReply, messageDTO.Status, messageDTO.InReply)
	}

	insert, err := db.Query(sqlString)

	if err != nil {
		fmt.Println(sqlString)
		fmt.Println(fmt.Sprintf("DB Error: %s", err.Error()))
		return
	}

	defer insert.Close()

	fmt.Println("DB: Order Recorded!")
}

func AddFollowUp(messageDTO dto.MessageDTO, config config.Config) {
	db, err := sql.Open("mysql", config.DB_PATH)

	if err != nil {
		fmt.Println(fmt.Sprintf("DB Error: %s", err.Error()))
		return
	}

	defer db.Close()

	editSQL := fmt.Sprintf("INSERT INTO player_order (message_id, player_name, player_id, message, is_genesis, refer_id, pnl, status) SELECT '%s', player_name, player_id, '%s', '%s', '%s', '%f', '%s' FROM player_order WHERE message_id = %s",
		messageDTO.MsgId, messageDTO.Message, "0", messageDTO.InReply, messageDTO.PnL, messageDTO.Status, messageDTO.InReply)

	if messageDTO.Status != "edited" {

		insert, err := db.Query(editSQL)

		if err != nil {
			fmt.Println(editSQL)
			fmt.Println(fmt.Sprintf("DB Error: %s", err.Error()))
			return
		}
		defer insert.Close()
	}

	// Update
	updateSQL := fmt.Sprintf("UPDATE player_order SET pnl=%f, status='%s' WHERE message_id=%s",
		messageDTO.PnL, messageDTO.Status, messageDTO.InReply)

	if messageDTO.Status == "stoploss" && messageDTO.PnL == 0.0 {
		updateSQL = fmt.Sprintf("UPDATE player_order SET status='%s' WHERE message_id=%s",
			messageDTO.Status, messageDTO.InReply)
	}
	update, err := db.Query(updateSQL)

	if err != nil {
		fmt.Println(updateSQL)
		fmt.Println(fmt.Sprintf("SQL Error: %s", err.Error()))
		return
	}

	defer update.Close()

	fmt.Println("Order Recorded!")
}
