package bosh

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	//UpdatePersistentDisk = `UPDATE persistent_disks SET disk_cid=? WHERE disk_cid=?`
	GetPersistentDiskMapping  = `SELECT disk_cid, cid FROM persistent_disks INNER JOIN vms on persistent_disks.instance_id=vms.instance_id;`
)

type BOSH struct {
	DB *sqlx.DB
}

type diskMappingsRow struct {
	VmCID   string `db:"cid"`
	DiskCID string `db:"disk_cid"`
}

//func (b BOSH) UpdatePersistentDiskCIDs() {
//	return
//}

func (b BOSH) GetPersistentDiskMappings() (bool, error){
	var diskMappings []diskMappingsRow
	if err := b.DB.Select(&diskMappings, GetPersistentDiskMapping); err != nil {
		return false, err
	}

	for _, m := range diskMappings {
		fmt.Printf("vm cid: %s, disk cid: %s", m.VmCID, m.DiskCID)
	}

	return true, nil
}
