package idserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eyedeekay/sam3"
	"github.com/zieckey/goini"
)

func (i IDServer) ClientConfigBytes() []byte {
	str := `[` + i.TunName + `]
type=client
inbound.length=2
outbound.length=2
inbound.quantity=6
outbound.quantity=6
inbound.backupQuantity=4
outbound.backupQuantity=4
i2cp.fastReceive=true
i2cp.messageReliability=BestEffort
i2cp.gzip=true
i2cp.dontPublishLeaseSet=true
destination=` + i.I2PKeys.Addr().Base32() + `
keys=` + i.TunName + `
`
	return []byte(str)
}

func (i IDServer) ClientKeys() (sam3.I2PKeys, error) {
	if dir, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(dir, i.Conf.TunName+"client.i2pkeys")); os.IsNotExist(err) {
			if file, err := os.Create(filepath.Join(dir, i.Conf.TunName+"client.i2pkeys")); err != nil {
				defer file.Close()
				return sam3.I2PKeys{}, err
			} else {
				defer file.Close()
				if keys, err := i.SAM.NewKeys(); err == nil {
					if err = sam3.StoreKeysIncompat(keys, file); err != nil {
						return keys, err
					}
					return keys, nil
				} else {
					return keys, err
				}

			}
		}
		if file, err := os.Open(filepath.Join(i.Conf.FilePath, i.Conf.TunName+".i2pkeys")); err != nil {
			defer file.Close()
			return sam3.I2PKeys{}, err
		} else {
			defer file.Close()
			return sam3.LoadKeysIncompat(file)
		}
	} else {
		return sam3.I2PKeys{}, err
	}
}

func (i IDServer) AccessList() string {
	var s string
	for _, v := range i.Conf.AccessList {
		s += v + ","
	}
	return strings.TrimRight(s, ",")
}

func (i IDServer) ServerConfigBytes() []byte {
	keys, err := i.ClientKeys()
	if err != nil {
		panic(err)
	}
	str := `[` + i.TunName + `]
type=server
inbound.length=2
outbound.length=2
inbound.quantity=6
outbound.quantity=6
inbound.backupQuantity=4
outbound.backupQuantity=4
i2cp.fastReceive=true
i2cp.messageReliability=BestEffort
i2cp.gzip=true
i2cp.enableWhitelist=true
i2cp.accessList=` + i.AccessList() + "," + keys.Addr().Base64() + `
destination=` + i.I2PKeys.Addr().Base32() + `.b32.i2p
keys=` + i.TunName + `
`
	return []byte(str)
}

func (i IDServer) GenerateClientConfig() error {
	if !i.generatekeys {
		return nil
	}
	raw := i.ClientConfigBytes()
	ini := goini.New()
	err := ini.Parse(raw, "\n", "=")
	if err != nil {
		return fmt.Errorf("parse INI memory data failed : %v", err.Error())
	}
	if dir, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(dir, i.Conf.TunName+"client.ini")); os.IsNotExist(err) {
			if file, err := os.Create(filepath.Join(dir, i.Conf.TunName+"client.ini")); err != nil {
				defer file.Close()
				return err
			} else {
				if file, err := os.Open(filepath.Join(dir, i.Conf.TunName+"client.ini")); err != nil {
					defer file.Close()
					return err
				} else {
					defer file.Close()
					return ini.Write(file)
				}
			}
		} else {
			if file, err := os.Open(filepath.Join(dir, i.Conf.TunName+"client.ini")); err != nil {
				defer file.Close()
				return err
			} else {
				defer file.Close()
				return ini.Write(file)
			}
		}
	} else {
		return err
	}
	return nil
}
