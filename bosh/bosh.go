package bosh

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	UpdatePersistentDisk     = "UPDATE persistent_disks SET disk_cid=$1 WHERE disk_cid=$2"
	GetPersistentDiskMapping = `SELECT disk_cid, cid FROM persistent_disks INNER JOIN instances ON persistent_disks.instance_id=instances.id INNER JOIN deployments ON instances.deployment_id=deployments.id INNER JOIN vms ON persistent_disks.instance_id=vms.instance_id WHERE deployments.name=$1;`
)

const DefaultDiskID = "3"

type BOSH struct {
	DB *sqlx.DB
}

type diskMappingsRow struct {
	VmCID   string `db:"cid"`
	DiskCID string `db:"disk_cid"`
}

func parseDiskID(id string) string {
	if id == "" {
		return DefaultDiskID
	}

	if strings.HasPrefix(id, "_") {
		return strings.SplitAfter(id, "_")[1]
	}

	return id
}

func (b BOSH) UpdatePersistentDiskCIDs(deployment string, id string) error {
	diskID := parseDiskID(id)

	diskMappings, err := b.getPersistentDiskMappings(deployment)
	if err != nil {
		return err
	}

	db := b.DB
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(UpdatePersistentDisk)
	if err != nil {
		return err
	}

	for _, m := range diskMappings {
		fmt.Printf("UPDATE persistent_disks SET disk_cid=%s_%s WHERE disk_cid=%s\n", m.VmCID, diskID, m.DiskCID)
		_, err = stmt.Exec(m.VmCID+diskID, m.DiskCID)
		if err != nil {
			_ = stmt.Close()
			return err
		}
	}
	err = stmt.Close()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (b BOSH) getPersistentDiskMappings(deployment string) ([]diskMappingsRow, error) {
	var diskMappings []diskMappingsRow
	if err := b.DB.Select(&diskMappings, GetPersistentDiskMapping, deployment); err != nil {
		return diskMappings, err
	}

	return diskMappings, nil
}
