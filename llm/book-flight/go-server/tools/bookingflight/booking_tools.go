/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bookingflight

import (
	"encoding/json"

	"context"
	"fmt"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
)

/*
SearchFlightTicket
*/
type SearchFlightTicket struct {
	tools.BaseTool
}

type serachFlightTicketData struct {
	Origin             string `json:"origin" validate:"required"`
	Destination        string `json:"destination" validate:"required"`
	Date               string `json:"date" validate:"required"`
	DepartureTimeStart string `json:"departure_time_start"`
	DepartureTimeEnd   string `json:"departure_time_end"`
}

func NewSearchFlightTicket(name string, description string) SearchFlightTicket {
	return SearchFlightTicket{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(serachFlightTicketData{}), "", "", ""),
	}
}

// origin string, destination string, date string, departureTimeStart string, departureTimeEnd string
func (stt SearchFlightTicket) Call(ctx context.Context, input string) (string, error) {
	data := serachFlightTicketData{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), err
	}

	return stt.searchFlightTicket(data)
}

func (stt SearchFlightTicket) searchFlightTicket(data serachFlightTicketData) (string, error) {
	// 此处只做出发地校验，其他信息未进行校验
	if data.Origin != "北京" {
		return "未查询到相关内容", nil
	}

	rst := flightInformation()
	rst_json, err := json.Marshal(rst)
	return string(rst_json), err
}

/*
PurchaseFlightTicket
*/
type PurchaseFlightTicket struct {
	tools.BaseTool
}

type purchaseFlightTicketData struct {
	FlightNumber string `json:"flight_number" validate:"required"`
}

func NewPurchaseFlightTicket(name string, description string) PurchaseFlightTicket {
	return PurchaseFlightTicket{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(purchaseFlightTicketData{}), "", "", ""),
	}
}

func (ptt PurchaseFlightTicket) Call(ctx context.Context, input string) (string, error) {
	data := purchaseFlightTicketData{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), err
	}

	return ptt.purchaseFlightTicket(data)
}

func (ptt *PurchaseFlightTicket) purchaseFlightTicket(data purchaseFlightTicketData) (string, error) {
	flightInfo := flightInformation()
	for _, info := range flightInfo {
		if data.FlightNumber == info["flight_number"] {
			rst_json, err := json.Marshal(info)
			return string(rst_json), err
		}
	}

	return fmt.Sprintf("The flight was not found: %v", data.FlightNumber), nil
}

/*
FinishPlaceholder
*/
type FinishPlaceholder struct {
	tools.BaseTool
}

func NewFinishPlaceholder(name string, description string) FinishPlaceholder {
	return FinishPlaceholder{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt FinishPlaceholder) Call(ctx context.Context, input string) (string, error) {
	return "FINISH", nil
}

func flightInformation() []map[string]string {
	return []map[string]string{
		{
			"flight_number":  "MU5100",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": "2024-06-01 07:00",
			"arrival_time":   "2024-06-01 09:15",
			"price":          "900.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "MU6865",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": "2024-06-01 07:20",
			"arrival_time":   "2024-06-01 09:25",
			"price":          "1160.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "HM7601",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": "2024-06-01 07:30",
			"arrival_time":   "2024-06-01 09:55",
			"price":          "1080.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CA1515",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": "2024-06-01 15:45",
			"arrival_time":   "2024-06-01 17:55",
			"price":          "1080.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "GS9012",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": "2024-06-01 19:00",
			"arrival_time":   "2024-06-01 23:00",
			"price":          "1250.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "GS9013",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": "2024-06-01 18:30",
			"arrival_time":   "2024-06-01 22:00",
			"price":          "1200.00",
			"seat_type":      "头等舱",
		},
	}
}
