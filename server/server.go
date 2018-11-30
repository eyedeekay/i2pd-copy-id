package idserver

import (
	"fmt"
	"github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam3"
	"net"
	"os"
	"path/filepath"
)

type IDServer struct {
	*sam3.StreamListener
	*sam3.StreamSession
	*sam3.SAM
	sam3.I2PKeys
	*i2ptunconf.Conf
	generatekeys bool
	path         string
	host         string
	port         string
}

func (i IDServer) samaddr() string {
	return i.host + ":" + i.port
}

func (i IDServer) ListenAndServe() (sam3.I2PKeys, error) {
	var err error
	i.StreamListener, i.I2PKeys, err = i.Listen()
	if i.HandleKill(err) != nil {
		return i.I2PKeys, err
	}
	for {
		if _, err = i.accept(); err != nil {
			return i.I2PKeys, err
		}
	}
	return i.I2PKeys, nil
}

func (i IDServer) accept() (sam3.I2PKeys, error) {
	SAMConn, err := i.AcceptI2P()
	if i.HandleKill(err) != nil {
		return i.I2PKeys, err
	}
	defer SAMConn.Close()
	readb := make([]byte, 1024)
	if _, err := SAMConn.Read(readb); err != nil {
		return i.I2PKeys, err
	}
	//writeb = make([]byte, 1024)
	if _, err := SAMConn.Write(readb); err != nil {
		//if _, err := i.Write(writeb); err != nil {
		return i.I2PKeys, err
	}
	return i.I2PKeys, nil
}

func (i IDServer) Listen() (*sam3.StreamListener, sam3.I2PKeys, error) {
	var err error
	i.StreamListener, err = i.StreamSession.Listen()
	if i.HandleError(err) != nil {
		return nil, i.I2PKeys, err
	}
	return i.StreamListener, i.I2PKeys, nil
}

func (i IDServer) Accept() (net.Conn, error) {
	return i.AcceptI2P()
}

func (i IDServer) AcceptI2P() (*sam3.SAMConn, error) {
	SAMConn, err := i.AcceptI2P()
	if i.HandleError(err) != nil {
		return nil, err
	}
	return SAMConn, nil
}

func (i IDServer) Close() error {
	return i.HandleKill(fmt.Errorf("%s", "Killing server."))
}

func (i IDServer) WriteConfig() error {
	if file, err := os.Open(filepath.Join(i.path)); err != nil {
		defer file.Close()
		return err
	} else {
		defer file.Close()
		return i.Conf.Config.Write(file)
	}

}

func (i IDServer) HandleKill(err error) error {
	if err != nil {
		if i.HandleError(i.StreamSession.Close()) != nil {
			return err
		}
		if i.HandleError(i.StreamListener.Close()) != nil {
			return err
		}
		if i.HandleError(i.SAM.Close()) != nil {
			return err
		}
	}
	return err
}

func (i IDServer) HandleError(err error, s ...interface{}) error {
	if err != nil {
		return fmt.Errorf("Error encountered: %s", s...)
	}
	return err
}

func (i IDServer) GetKeys() (sam3.I2PKeys, error) {
	if _, err := os.Stat(filepath.Join(i.Conf.FilePath, i.Conf.TunName+".i2pkeys")); os.IsNotExist(err) {
		if file, err := os.Create(filepath.Join(i.Conf.FilePath, i.Conf.TunName+".i2pkeys")); err != nil {
			defer file.Close()
			return i.I2PKeys, err
		} else {
			defer file.Close()
			var err error
			if i.I2PKeys, err = i.SAM.NewKeys(); err == nil {
				if err = sam3.StoreKeysIncompat(i.I2PKeys, file); err != nil {
					return i.I2PKeys, err
				}
				return i.I2PKeys, nil
			} else {
				return i.I2PKeys, err
			}
		}
	}
	if file, err := os.Open(filepath.Join(i.Conf.FilePath, i.Conf.TunName+".i2pkeys")); err != nil {
		defer file.Close()
		return i.I2PKeys, err
	} else {
		defer file.Close()
		return sam3.LoadKeysIncompat(file)
	}
}

func NewIDServer(opts ...func(*IDServer) error) (*IDServer, error) {
	var i IDServer
	i.generatekeys = false
	i.path = "/etc/i2pd/tunnels.conf.d/i2pd-copy-id.conf"
	i.host = "127.0.0.1"
	i.port = "7656"
	for _, o := range opts {
		if err := o(&i); err != nil {
			return nil, err
		}
	}
	var err error
	i.SAM, err = sam3.NewSAM(i.samaddr())
	if i.HandleKill(err) != nil {
		return nil, err
	}
	i.I2PKeys, err = i.GetKeys()
	if i.HandleKill(err) != nil {
		return nil, err
	}
	i.StreamSession, err = i.SAM.NewStreamSession("i2p-copy-id", i.I2PKeys, i.Conf.Print())
	if i.HandleKill(err) != nil {
		return nil, err
	}
	if err = i.GenerateClientConfig(); err != nil {
		return nil, err
	}
	return &i, nil
}
