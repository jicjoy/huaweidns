package huaweidns

import (
	"context"
	"fmt"
	"sync"

	"github.com/libdns/libdns"
)

// Provider implements the libdns interfaces for Alicloud.
type Provider struct {
	client *DnsClient
	// The API Key ID Required by Aliyun's for accessing the Aliyun's API
	AccKeyID string `json:"access_key_id"`
	// The API Key Secret Required by Aliyun's for accessing the Aliyun's API
	AccKeySecret string `json:"access_key_secret"`
	// Optional for identifing the region of the Aliyun's Service,The default is zh-hangzhou
	RegionID string `json:"region_id,omitempty"`
	ZoneName string `json:"zone_name,omitempty"`
	mu       sync.Mutex
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	var rls []libdns.Record

	fmt.Printf("AppendRecords:%s\n", zone)
	p.getClient(ctx, zone)
	for _, rec := range recs {
		ar := ToHuaweiDnsRecord(rec, zone)
		fmt.Printf("ar: %+v", ar)
		p.GetZoneByName(ctx, ar.ZoneName)
		if ar.ID == "" {
			rId, _ := p.client.GetRecordLists(ctx, ar.Name, ar.Type)
			if len(rId.Response) > 0 {
				ar.ID = rId.Response[0].ID
			}
		}

		res, err := p.UpdateOrcreateRecord(ctx, &ar)
		if err != nil {
			return rls, err
		}

		rls = append(rls, res.LibdnsRecord())
		fmt.Printf("res: %+v", res.LibdnsRecord())
	}
	return rls, nil
}

// DeleteRecords deletes the records from the zone. If a record does not have an ID,
// it will be looked up. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	var rls []libdns.Record
	p.getClient(ctx, zone)
	for _, rec := range recs {
		ar := ToHuaweiDnsRecord(rec, zone)
		fmt.Printf("ar: %+v", ar)
		p.GetZoneByName(ctx, ar.ZoneName)
		if len(ar.ID) == 0 {
			r0, err := p.client.GetRecordLists(ctx, ar.Name, ar.Type)
			ar.ID = r0.Response[0].ID
			if err != nil {
				return nil, err
			}
		}
		_, err := p.client.DeleteRecord(ctx, ar.ID)
		if err != nil {
			return nil, err
		}
		rls = append(rls, ar.LibdnsRecord())
	}
	return rls, nil
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	var rls []libdns.Record
	fmt.Printf("get:%s", zone)
	p.getClient(ctx, zone)

	recs, err := p.client.GetRecordLists(ctx, "", "")
	if err != nil {
		return nil, err
	}
	for _, rec := range recs.Response {
		rls = append(rls, rec.LibdnsRecord())
	}
	return rls, nil
}

// SetRecords sets the records in the zone, either by updating existing records
// or creating new ones. It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	var rls []libdns.Record
	fmt.Printf("Set:%s", zone)
	p.getClient(ctx, zone)
	for _, rec := range recs {
		ar := ToHuaweiDnsRecord(rec, zone)
		res, err := p.UpdateOrcreateRecord(ctx, &ar)
		if err != nil {
			return nil, err
		}
		rls = append(rls, res.LibdnsRecord())
	}
	return rls, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
