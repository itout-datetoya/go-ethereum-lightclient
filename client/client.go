package client

import (
	"context"
	"time"
	"errors"
	"log"
	"itout/go-ethereum-lightclient/sync"
	"itout/go-ethereum-lightclient/util"
	"itout/go-ethereum-lightclient/configs"
	"github.com/protolambda/ztyp/tree"

)

type Client struct {
	BeaconBaseURL string
	TrustedRoot string
	Spec configs.Spec
}

func (c *Client) StartClient(ctx context.Context) error {
	trustedRoot := util.HexstrTo32Bytes(c.TrustedRoot)
	bootstrap := sync.GetBootstrap(trustedRoot, c.BeaconBaseURL)
	store, err := sync.InitStore(tree.Root(trustedRoot), bootstrap)
	if err != nil {
		return errors.New("Client failed to start")
	}

	finalityTicker := time.NewTicker(time.Second)
	defer finalityTicker.Stop()
	updateTicker := time.NewTicker(60 * time.Second)
	defer updateTicker.Stop()

	for {
		select {
		case <-finalityTicker.C:
			update := sync.GetFinalityUpdate(c.BeaconBaseURL)
			if store.Header.Slot < update.AttestedHeader.Slot {
				err := store.FinalityUpdateStore(update, &c.Spec)
				if err != nil {
					log.Printf("%+v", err)
				} else {
					log.Printf("Update: Slot %d", store.Header.Slot)
				}
			}

		case <-updateTicker.C:
			update := sync.GetUpdate(store.Header.Slot, c.BeaconBaseURL)
			err := store.UpdateStore(update, &c.Spec)
			if err != nil {
				log.Printf("%+v", err)
			} else {
				log.Printf("Update: Sync Committee")
			}

		case <-ctx.Done():
			log.Printf("Stopping client")
			return nil
		}
	}
}