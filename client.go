package huaweidns

import (
	"context"
	"fmt"
)

func (p *Provider) getClient(ctx context.Context, zoneName string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ZoneName = zoneName
	if p.client == nil {
		config := &WithConfig{
			AccKeyID:     p.AccKeyID,
			AccKeySecret: p.AccKeySecret,
			Version:      "v2",
			RegionID:     p.RegionID,
		}
		b := new(BuilderApi).WithConfig(config)
		p.client, _ = b.Build()

	}

	if ValidateZone(zoneName) {
		p.GetZoneByName(ctx, zoneName)
	}
	fmt.Printf("zoneName: %s", zoneName)

}

func (p *Provider) GetZoneByName(ctx context.Context, name string) {

	if len(p.client.ZoneID) == 0 && len(name) > 0 {
		p.client.GetZoneList(ctx, name, false)
	}

}
func (p *Provider) ValidateZone() error {
	if len(p.client.ZoneID) == 0 {
		return fmt.Errorf("the Zone Name error %s", p.ZoneName)
	}

	return nil
}
func (p *Provider) UpdateOrcreateRecord(ctx context.Context, rec *RecordTag) (RecordTag, error) {

	if rec.ID == "" {
		return p.client.CreateRecord(ctx, rec)
	}
	return p.client.UpdateRecord(ctx, rec)
}
