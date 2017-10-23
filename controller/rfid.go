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
	var formData map[string]string
	err := BindRequestBodyWithTeeReader(g, &formData)
	if err != nil {
		g.JSON(http.StatusBadRequest, &response.Response{
			Code:    response.Error,
			Message: "Arguments error.",
			Data:    "",
		})
		return
	}
	rfids, exist := formData["rfids"]
	if !exist {
		g.JSON(http.StatusBadRequest, &response.Response{
			Code:    response.Error,
			Message: "Arguments error.",
			Data:    "",
		})
		return
	}
	barcodes := map[string]string{}
	for _, rfid := range strings.Split(rfids, ",") {
		barcodes[rfid] = ParseRfidToBarcode(rfid, conf.RFID_ENCODE_COUNT, conf.RFID_REDUNDANT_COUNT)
	}
	g.JSON(http.StatusOK, &response.Response{
		Code:    response.OK,
		Message: "",
		Data:    barcodes,
	})
}

func RfidEncoder(g *gin.Context) {
	var formData map[string]string
	err := BindRequestBodyWithTeeReader(g, &formData)
	if err != nil {
		g.JSON(http.StatusBadRequest, &response.Response{
			Code:    response.Error,
			Message: "Arguments error.",
			Data:    "",
		})
		return
	}
	codes, exist := formData["codes"]
	if !exist {
		g.JSON(http.StatusBadRequest, &response.Response{
			Code:    response.Error,
			Message: "Arguments error.",
			Data:    "",
		})
		return
	}
	rfids := map[string]string{}
	for _, code := range strings.Split(codes, ",") {
		rfids[code] = ParseBarcodeToRfid(code, conf.RFID_ENCODE_COUNT)
	}
	g.JSON(http.StatusOK, &response.Response{
		Code:    response.OK,
		Message: "",
		Data:    rfids,
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

func ParseRfidToBarcode(rfidCode string, encodeCount int, redundantCount int) string {
	var barcode string
	if len(rfidCode) > redundantCount {
		rfidCode = rfidCode[:len(rfidCode)-redundantCount]
	}
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
