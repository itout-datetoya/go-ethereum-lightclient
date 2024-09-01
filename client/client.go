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
	BaseURL string
	TrustedRoot string
	Spec configs.Spec
}

func (c *Client) StartClient(ctx context.Context) error {
	trustedRoot := util.HexstrTo32Bytes(c.TrustedRoot)
	bootstrap := sync.GetBootstrap(trustedRoot)
	store, err := sync.InitStore(tree.Root(trustedRoot), bootstrap)
	if err != nil {
		return errors.New("Client failed to start")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			update := sync.GetUpdate(store.Header.Slot)
			err := store.UpdateStore(update, &c.Spec)
			if err != nil {
				log.Printf("Error Updating data")
			} else {
				log.Printf("Update: Slot %d", store.Header.Slot)
			}
			continue
		case <-ctx.Done():
			log.Printf("Stopping client")
			return nil
		}
	}
}