package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sysdiglabs/stackdriver-webhook-bridge/converter"
	"github.com/sysdiglabs/stackdriver-webhook-bridge/model"
	"google.golang.org/genproto/googleapis/cloud/audit"
	"log"
	"os"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var entry model.SavedLoggingEntry
		err := json.Unmarshal([]byte(s.Text()), &entry)
		if err != nil {
			log.Fatalf("Could not decode log entry as json: %v", err)
		}
		var auditPayload audit.AuditLog
		err = jsonpb.UnmarshalString(entry.AuditPayload, &auditPayload)
		if err != nil {
			log.Fatalf("Could not unmarshal log entry: %v", err)
		}
		actualAuditEvent, err := converter.ConvertLogEntrytoAuditEvent(entry.Entry, &auditPayload)
		if err != nil {
			log.Fatalf("Could not conver log entry: %v", err)
		}
		actualAuditEventJSON, err := json.MarshalIndent(*actualAuditEvent, "", "  ")
		fmt.Println(string(actualAuditEventJSON))
	}
}
