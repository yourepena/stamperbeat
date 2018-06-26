package beater

import (
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/yourepena/qwcfp-client-go"
	"github.com/yourepena/stamperbeat/config"
	"time"
)

type Stamperbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig

	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Stamperbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Stamperbeat) Run(b *beat.Beat) error {
	logp.Info("stamperbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		username := bt.config.Username
		password := bt.config.Password
		dnsServer := bt.config.UrlWS
		groupName := bt.config.Group
		rootConfig := bt.config.RootConfigXML

		fileVersionArray, errorR := getFiles(username, password, dnsServer, groupName, rootConfig)

		if errorR != nil {
			logp.Info("Aqui finalizou")
			logp.Error(errorR)
			return errorR
		}

		for u := 0; u < len(fileVersionArray); u++ {

			FileName := fileVersionArray[u].FileName
			Path := fileVersionArray[u].Path
			FileVersionId := fileVersionArray[u].FileVersionId
			FileId := fileVersionArray[u].FileId
			Groupid := fileVersionArray[u].Groupid

			logp.Info("FileName: %s\nPath: %s\nFileVersionId: %d\nFileId: %d\nGroupid: %d\n\n\n", FileName, Path, FileVersionId, FileId, Groupid)

		}

		bt.client.Close()
		close(bt.done)

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Stamperbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

type FileVersionRetorno struct {
	FileName      string
	Path          string
	FileVersionId int
	FileId        int
	Groupid       int
}

func getFiles(username string, password string, dnsServer string, groupName string, rootConfig string) ([]FileVersionRetorno, error) {

	loginKey, err := soap.Login(username, password, dnsServer, rootConfig)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var retorno = []FileVersionRetorno{}

	groupid, errG := soap.GetGroup("GFNS_PROCESSADO", loginKey, dnsServer, rootConfig)

	if errG != nil {
		fmt.Println(errG)
		return nil, err
	}

	fileVersionArray, errF := soap.GetFilesFromQWCFP(loginKey, groupName, dnsServer, rootConfig)

	for i := 0; i < len(fileVersionArray); i++ {
		fvr := FileVersionRetorno{
			FileName:      fileVersionArray[i].FileName,
			Path:          fileVersionArray[i].Path,
			FileVersionId: fileVersionArray[i].FileVersionId,
			FileId:        fileVersionArray[i].FileId,
			Groupid:       groupid,
		}

		retorno = append(retorno, fvr)

	}

	if errF != nil {
		fmt.Println(errF)
		return nil, errF
	}

	return retorno, nil

}

//strconv.Itoa

func moveFile(FileVersionId string, groupid string, loginKey string, dnsServer string, rootConfig string) error {
	_, errC := soap.MoveFile(FileVersionId, groupid, loginKey, dnsServer, rootConfig)
	return errC
}
