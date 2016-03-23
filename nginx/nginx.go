package nginx

import (
	"github.com/xplacepro/common"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	AVAILABLE_PATH = "/etc/nginx/sites-available"
	ENABLED_PATH   = "/etc/nginx/sites-enabled"
)

type Domain struct {
	Domain string
	Config string
}

func (d *Domain) AvailablePath() string {
	return path.Join(AVAILABLE_PATH, d.Domain)
}

func (d *Domain) EnabledPath() string {
	return path.Join(ENABLED_PATH, d.Domain)
}

func (d *Domain) Create() error {
	log.Printf("Creating domain %s", d.Domain)
	if err := ioutil.WriteFile(d.AvailablePath(), []byte(d.Config), 755); err != nil {
		log.Printf("Error reating domain %s, %s", d.Domain, err)
		return err
	}
	symlink := func() error {
		return os.Symlink(d.AvailablePath(), d.EnabledPath())
	}
	if err := symlink(); err != nil {
		if os.IsExist(err) {
			os.Remove(d.EnabledPath())
			if err_again := symlink(); err_again != nil {
				log.Printf("Error reating domain %s, %s", d.Domain, err_again)
				return err_again
			}
		} else {
			log.Printf("Error reating domain %s, %s", d.Domain, err)
			return err
		}
	}
	log.Printf("Created domain %s", d.Domain)
	return nil
}

func (d *Domain) Exists() error {
	if err := os.Remove(d.AvailablePath()); err != nil {
		return err
	}
	if err := os.Remove(d.EnabledPath()); err != nil {
		return err
	}
	return nil
}

func (d *Domain) Delete() error {
	if err := os.Remove(d.AvailablePath()); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.Remove(d.EnabledPath()); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func Reload() error {
	if _, err := common.RunCommand("/usr/sbin/service", []string{"nginx", "reload"}); err != nil {
		log.Printf("Error reloading nginx, %s", err)
		return err
	}
	return nil
}

func Test() error {
	if _, err := common.RunCommand("nginx", []string{"-t"}); err != nil {
		return err
	}
	return nil
}
