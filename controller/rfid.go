package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/response"
	"github.com/gin-gonic/gin"
)

func RfidDecoder(g *gin.Context) {
	barcode := ParseRfidToBarcode(g.Param("rfid"), conf.RFID_ENCODE_COUNT)
	g.JSON(http.StatusOK, &response.Response{
		Code:    response.OK,
		Message: "",
		Data:    barcode,
	})
}

func RfidEncoder(g *gin.Context) {
	rfid := ParseBarcodeToRfid(g.Param("code"), conf.RFID_ENCODE_COUNT)
	g.JSON(http.StatusOK, &response.Response{
		Code:    response.OK,
		Message: "",
		Data:    rfid,
	})
}

func ParseBarcodeToRfid(barcode string, encodeCount int) string {
	var head, rfidString string
	for barcode != "" {
		var tmpStr string
		head, barcode = barcode[:1], barcode[1:]
		tmpStr = strconv.Itoa(int([]byte(head)[0]))
		if len(tmpStr) < encodeCount {
			tmpStr = strings.Repeat("0", encodeCount-len(tmpStr)) + tmpStr
		}
		rfidString = rfidString + tmpStr
	}
	if len(rfidString)%4 != 0 {
		rfidString = rfidString + strings.Repeat("0", 4-len(rfidString)%4)
	}
	return rfidString
}

func ParseRfidToBarcode(rfidCode string, encodeCount int) string {
	var barcode string
	if len(rfidCode)%encodeCount != 0 {
		rfidCode = rfidCode[:len(rfidCode)-len(rfidCode)%encodeCount]
	}
	for rfidCode != "" {
		var head, tmpStr string
		head, rfidCode = rfidCode[:encodeCount], rfidCode[encodeCount:]

		if i, err := strconv.Atoi(head); err != nil {
			fmt.Println(err)
		} else {
			tmpStr = string([]byte{uint8(i)})
		}
		barcode = barcode + tmpStr
	}
	return barcode
}
