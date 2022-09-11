package dto

import (
	"fmt"
	"regexp"
	"strconv"
)

type MessageDTO struct {
	MsgId   string
	ByWho   string
	ByWhoID string
	Time    string
	InReply string
	Message string
	Status  string
	PnL     float64
}

func MessageToDTO(input string) MessageDTO {

	byWho, byWhoID, time, inReply, status := "", "", "", "-1", "created"
	pnl := 0.0

	// Message Infomation
	byRe := regexp.MustCompile(`(By\s?=\s?)([^\(]+)(\()([0-9]+)`)
	timeRe := regexp.MustCompile(`(Time\s?=\s?)([^\n]+)`)

	byWhoList := byRe.FindStringSubmatch(input)
	if len(byWhoList) == 5 {
		byWho = byWhoList[2]
	}

	byWhoIDList := byRe.FindStringSubmatch(input)
	if len(byWhoIDList) == 5 {
		byWhoID = byWhoIDList[4]
	}

	timeList := timeRe.FindStringSubmatch(input)
	if len(timeList) == 3 {
		time = timeList[2]
	}

	// Stop Loss
	stopLossRe := regexp.MustCompile(`(?i)(stoploss([^\n]+)(\s+)?)(loss(\s+)?:(\s+)?)([0-9]+(.)?([0-9]+)?)`)
	stopLossTarget1 := regexp.MustCompile(`(?i)(Closed at trailing stoploss after reaching take profit)`)

	stopLossList := stopLossRe.FindStringSubmatch(input)
	if len(stopLossList) >= 8 {
		fpnl, err := strconv.ParseFloat(stopLossList[7], 4)
		pnl = fpnl * -1

		status = "stoploss"
		if err != nil {
			pnl = 0.0000
			fmt.Println(fmt.Sprintf("Parsing Error(Stoploss): %s", err.Error()))
		}
	}

	stopLossTarget1List := stopLossTarget1.FindStringSubmatch(input)
	if len(stopLossTarget1List) == 2 {
		status = "stoploss"
	}

	// Target Enter
	partialEntryTargetRe := regexp.MustCompile(`(?i)(Entry target 1)`)
	fullEntryTargetRe := regexp.MustCompile(`(?i)(All entry targets achieved)`)

	partialEntryTargetList := partialEntryTargetRe.FindStringSubmatch(input)
	if len(partialEntryTargetList) == 2 {
		status = "enter 1 target"
	}

	fullEntryTargetList := fullEntryTargetRe.FindStringSubmatch(input)
	if len(fullEntryTargetList) == 2 {
		status = "enter all target"
	}

	// Take Profit
	takeProfit1 := regexp.MustCompile(`(?i)(Take-Profit target 1([^\n]+)(\s+)?)(Profit(\s+)?:(\s+)?)([0-9]+(.)?([0-9]+)?)`)
	takeProfitAll := regexp.MustCompile(`(?i)(ll take-profit targets achieved([^\n]+)(\s+)?)(Profit(\s+)?:(\s+)?)([0-9]+(.)?([0-9]+)?)`)

	takeProfit1List := takeProfit1.FindStringSubmatch(input)
	if len(takeProfit1List) >= 8 {
		fpnl, err := strconv.ParseFloat(takeProfit1List[7], 4)
		pnl = fpnl

		status = "take profit target 1"
		if err != nil {
			pnl = 0.0000
			fmt.Println(fmt.Sprintf("Parsing Error(TakeProfit1): %s", err.Error()))
		}
	}

	takeProfitAllList := takeProfitAll.FindStringSubmatch(input)
	if len(takeProfitAllList) >= 8 {
		fpnl, err := strconv.ParseFloat(takeProfitAllList[7], 4)
		pnl = fpnl

		status = "completed"
		if err != nil {
			fmt.Println(fmt.Sprintf("Parsing Error(Completed): %s", err.Error()))
			pnl = 0.0000
		}
	}

	// Edit
	editRe := regexp.MustCompile(`(?i)(Published By: @kmkok)`)
	editReList := editRe.FindStringSubmatch(input)
	if len(editReList) == 2 {
		status = "edited"
	}

	// Cancel
	cancelledRe := regexp.MustCompile(`(?i)(cancelled)`)
	cancelledReList := cancelledRe.FindStringSubmatch(input)
	if len(cancelledReList) == 2 {
		status = "cancelled"
	}

	messageDTO := MessageDTO{
		ByWho:   byWho,
		ByWhoID: byWhoID,
		Time:    time,
		InReply: inReply,
		Message: input,
		Status:  status,
		PnL:     pnl}

	return messageDTO
}
