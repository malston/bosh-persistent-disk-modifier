package bosh

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/malston/bosh-persistent-disk-modifier/vc"
	"github.com/vmware/govmomi/vim25"
)

const (
	UpdatePersistentDisk     = "UPDATE persistent_disks SET disk_cid=$1 WHERE disk_cid=$2"
	GetPersistentDiskMapping = `SELECT disk_cid, cid FROM persistent_disks INNER JOIN instances ON persistent_disks.instance_id=instances.id INNER JOIN deployments ON instances.deployment_id=deployments.id INNER JOIN vms ON persistent_disks.instance_id=vms.instance_id WHERE deployments.name=$1;`
)

const DefaultDiskID = "3"

type Repository struct {
	DB *sqlx.DB
}

type diskMappingsRow struct {
	VmCID   string `db:"cid"`
	DiskCID string `db:"disk_cid"`
}

func (r Repository) UpdatePersistentDiskCIDs(deployment string, c *vim25.Client) error {
	diskMappings, err := r.getPersistentDiskMappings(deployment)
	if err != nil {
		return err
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(UpdatePersistentDisk)
	if err != nil {
		return err
	}

	for _, m := range diskMappings {
		err, diskName := vc.GetPersistentDiskName(context.Background(), c, m.VmCID)
		if err != nil {
			_ = stmt.Close()
			return err
		}
		if diskName == "" {
			diskName = fmt.Sprintf("%s_%s", m.VmCID, DefaultDiskID)
		}

		fmt.Printf("UPDATE persistent_disks SET disk_cid=%s WHERE disk_cid=%s\n", diskName, m.DiskCID)
		_, err = stmt.Exec(diskName, m.DiskCID)
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

func (r Repository) getPersistentDiskMappings(deployment string) ([]diskMappingsRow, error) {
	var diskMappings []diskMappingsRow
	if err := r.DB.Select(&diskMappings, GetPersistentDiskMapping, deployment); err != nil {
		return diskMappings, err
	}

	return diskMappings, nil
}
