package huaweidns

import (
	"fmt"
	"strings"
	"time"

	"github.com/libdns/libdns"
)

type Tags struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type links struct {
	Self string `json:"self,omitempty"`
}
type RecordTag struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Ttl         uint32   `json:"ttl,omitempty"`
	Tags        []*Tags  `json:"tags,omitempty"`
	Links       links    `json:"links,omitempty"`
	Status      string   `json:"status,omitempty"`
	ZoneID      string   `json:"zone_id,omitempty"`
	ZoneName    string   `json:"zone_name,omitempty"`
	Type        string   `json:"type,omitempty"`
	ErrorMsg    string   `json:"message,omitempty"`
	Records     []string `json:"records,omitempty"`
}

type Zones struct {
	Response []RecordTag `json:"zones,omitempty"`
}

type Recordsets struct {
	Response []RecordTag `json:"recordsets,omitempty"`
}

type Recrod struct {
	Response RecordTag
}

func (r *RecordTag) LibdnsRecord() libdns.Record {
	return libdns.Record{
		ID:    r.ID,
		Type:  r.Type,
		Name:  r.Name,
		Value: strings.Join(r.Records, " "),
		TTL:   time.Duration(r.Ttl) * time.Second,
	}
}

func ToHuaweiDnsRecord(rec libdns.Record, zone string) RecordTag {

	return RecordTag{
		ID:       rec.ID,
		Name:     strings.Trim(libdns.AbsoluteName(libdns.RelativeName(rec.Name, zone), zone), ".") + ".",
		Ttl:      uint32(rec.TTL.Seconds()),
		Type:     strings.ToUpper(rec.Type),
		ZoneName: strings.Trim(zone, ".") + ".",
		Records:  getRecords(rec.Type, rec.Value),
	}
}

func getRecords(rType string, value string) []string {
	switch rType {
	case "TXT":
		value = fmt.Sprintf("\"%s\"", value)

	}
	return []string{value}
	//return []string{fmt.Sprintf("\"%s\"", value)}
}

func ValidateZone(zone string) bool {
	if len(zone) == 0 {
		return false
	}
	zoneArr := strings.Split(strings.Trim(zone, "."), ".")
	fmt.Println(zoneArr)
	return len(zoneArr) > 1

}
